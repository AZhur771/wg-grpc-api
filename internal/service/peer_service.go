package service

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"strings"
	"text/template"
	"time"

	"github.com/AZhur771/wg-grpc-api/internal/app"
	"github.com/AZhur771/wg-grpc-api/internal/entity"
	tmpl "github.com/AZhur771/wg-grpc-api/internal/template"
	"github.com/google/uuid"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type PeerService struct {
	logger   app.Logger
	wgclient app.Wg
	storage  app.PeerStorage

	reservedIPs *reservedIPs
	device      *deviceConfig
}

func NewPeerService(logger app.Logger, wgclient app.Wg, storage app.PeerStorage) *PeerService {
	return &PeerService{
		logger:   logger,
		wgclient: wgclient,
		storage:  storage,
	}
}

func (ps *PeerService) getAvailableIP(ipnet *net.IPNet) (*net.IPNet, error) {
	// this two addresses are not usable
	broadcastAddr := getBroadcastAddr(ipnet).String()
	networkAddr := ipnet.IP.String()

	for ip := ipnet.IP.Mask(ipnet.Mask); ipnet.Contains(ip); incAddr(ip) {
		// check that incremented ip is not reserved
		isReserved := ps.reservedIPs.Contains(ip)
		if isReserved {
			continue
		}

		addr := ip.String()
		if addr != networkAddr && addr != broadcastAddr {
			ps.reservedIPs.Add(ip)
			return &net.IPNet{
				IP:   ip,
				Mask: net.IPv4Mask(255, 255, 255, 255),
			}, nil
		}
	}

	return nil, ErrRunOutOfAddresses
}

func (ps *PeerService) Add(
	ctx context.Context,
	addPresharedKey bool,
	persistentKeepAlive time.Duration,
	description string,
) (uuid.UUID, error) {
	var id uuid.UUID

	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return id, fmt.Errorf("peer service: failed to generate private key: %w", err)
	}

	publicKey := privateKey.PublicKey()

	var presharedKey wgtypes.Key
	var hasPresharedKey bool

	if addPresharedKey {
		presharedKey, err = wgtypes.GenerateKey()
		if err != nil {
			return id, fmt.Errorf("peer service: failed to generate preshared key: %w", err)
		}
		hasPresharedKey = true
	}

	allowedIPNet, err := ps.getAvailableIP(&ps.device.net)
	if err != nil {
		return id, err
	}

	peer := entity.Peer{
		PrivateKey:                  privateKey,
		PublicKey:                   publicKey,
		PresharedKey:                presharedKey,
		PersistentKeepaliveInterval: persistentKeepAlive,
		AllowedIPs:                  []net.IPNet{*allowedIPNet},
		HasPresharedKey:             hasPresharedKey,
		IsEnabled:                   true, // Enabled by default
		Description:                 description,
	}

	peerConfig := wgtypes.PeerConfig{
		PublicKey:                   peer.PublicKey,
		PresharedKey:                &peer.PresharedKey,
		PersistentKeepaliveInterval: &peer.PersistentKeepaliveInterval,
		AllowedIPs:                  peer.AllowedIPs,
	}

	if err := ps.wgclient.ConfigureDevice(peerConfig); err != nil {
		return id, fmt.Errorf("peer service: failed to configure device: %w", err)
	}

	id, err = ps.storage.Create(ctx, peer.ToPersistedPeer())
	if err != nil {
		return id, err
	}

	ps.device.incPeerCount()

	return id, nil
}

func (ps *PeerService) Update(
	ctx context.Context,
	id uuid.UUID,
	addPresharedKey bool,
	persistentKeepAlive time.Duration,
	description string,
	updateMask []string,
) error {
	persistedPeer, err := ps.storage.Get(ctx, id)
	if err != nil {
		return err
	}

	peer, err := persistedPeer.ToPeer()
	if err != nil {
		return err
	}

	var shouldConfigureDev bool

	// TODO: looks like wrong usage of field mask, should be refactored
	if containsString(updateMask, "add_preshared_key") && addPresharedKey != peer.HasPresharedKey {
		peer.HasPresharedKey = false
		peer.PresharedKey = wgtypes.Key{} // non-nil zero value key clears preshared key
		shouldConfigureDev = true
	}

	if containsString(updateMask, "persistent_keep_alive") {
		peer.PersistentKeepaliveInterval = persistentKeepAlive
		shouldConfigureDev = true
	}

	if containsString(updateMask, "description") {
		peer.Description = description
	}

	if shouldConfigureDev {
		peerConfig := wgtypes.PeerConfig{
			PublicKey:                   peer.PublicKey,
			PresharedKey:                &peer.PresharedKey,
			PersistentKeepaliveInterval: &peer.PersistentKeepaliveInterval,
			UpdateOnly:                  true,
		}

		if err := ps.wgclient.ConfigureDevice(peerConfig); err != nil {
			return fmt.Errorf("peer service: failed to configure device: %w", err)
		}
	}

	return ps.storage.Update(ctx, id, peer.ToPersistedPeer())
}

func (ps *PeerService) Delete(ctx context.Context, id uuid.UUID) error {
	persistedPeer, err := ps.storage.Delete(ctx, id)
	if err != nil {
		return err
	}

	peer, err := persistedPeer.ToPeer()
	if err != nil {
		return err
	}

	for _, addr := range peer.AllowedIPs {
		ps.reservedIPs.Remove(addr.IP)
	}

	peerConfig := wgtypes.PeerConfig{
		PublicKey: peer.PublicKey,
		Remove:    true,
	}

	if err = ps.wgclient.ConfigureDevice(peerConfig); err != nil {
		return fmt.Errorf("peer service: failed to configure device: %w", err)
	}

	ps.device.decPeerCount()

	return nil
}

func (ps *PeerService) Get(ctx context.Context, id uuid.UUID) (*entity.Peer, error) {
	persistedPeer, err := ps.storage.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	peer, err := persistedPeer.ToPeer()
	if err != nil {
		return nil, err
	}

	if peer.IsEnabled {
		wgpeer, err := ps.wgclient.GetPeer(peer.PublicKey)
		if err != nil {
			return nil, err
		}

		peer = peer.PopulateDynamicFields(&wgpeer)
	}

	return peer, nil
}

func (ps *PeerService) GetAll(ctx context.Context) ([]*entity.Peer, error) {
	persistedPeers, err := ps.storage.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	wgpeers, err := ps.wgclient.GetPeers()
	if err != nil {
		return nil, err
	}

	peers := make([]*entity.Peer, 0, len(persistedPeers))

	// TODO: O(n^2) time complexity, should be refactored
	for _, persistedPeer := range persistedPeers {
		peer, err := persistedPeer.ToPeer()
		if err != nil {
			return nil, err
		}

		if peer.IsEnabled {
			for i, p := range wgpeers {
				if bytes.Equal(p.PublicKey[:], peer.PublicKey[:]) {
					peer = peer.PopulateDynamicFields(&wgpeers[i])
					break
				}
			}
		}

		peers = append(peers, peer)
	}

	return peers, nil
}

func (ps *PeerService) DownloadConfig(ctx context.Context, id uuid.UUID) ([]byte, error) {
	persistedPeer, err := ps.storage.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	t, err := template.New("config").Funcs(
		template.FuncMap{
			"StringsJoin": strings.Join,
		},
	).Parse(tmpl.ClientConfigTemplate)
	if err != nil {
		return nil, fmt.Errorf("peer service: failed to read peer config template: %w", err)
	}

	var buf bytes.Buffer

	tmplData := tmpl.TmplData{
		// interface data
		InterfacePrivateKey: persistedPeer.PrivateKey,
		InterfaceAddress:    persistedPeer.AllowedIPs,
		InterfaceDNS:        []string{"9.9.9.9", "149.112.112.112"},
		InterfaceMTU:        1384,

		// peer data (device)
		PeerPublicKey:       ps.device.publicKey.String(),
		PeerPresharedKey:    persistedPeer.PresharedKey,
		PeerEndpoint:        ps.device.endpoint.String(),
		PeerAllowedIPs:      []string{"0.0.0.0/0"},
		PersistentKeepalive: int(persistedPeer.PersistentKeepaliveInterval.Seconds()),
	}

	if err := t.Execute(&buf, tmplData); err != nil {
		return nil, err
	}

	return buf.Bytes(), err
}

func (ps *PeerService) Enable(ctx context.Context, id uuid.UUID) error {
	persistedPeer, err := ps.storage.Get(ctx, id)
	if err != nil {
		return err
	}

	peer, err := persistedPeer.ToPeer()
	if err != nil {
		return err
	}

	if !peer.IsEnabled {
		peer.IsEnabled = true
		peerConfig := wgtypes.PeerConfig{
			PublicKey:                   peer.PublicKey,
			PresharedKey:                &peer.PresharedKey,
			PersistentKeepaliveInterval: &peer.PersistentKeepaliveInterval,
			AllowedIPs:                  peer.AllowedIPs,
			UpdateOnly:                  true,
		}

		if err = ps.wgclient.ConfigureDevice(peerConfig); err != nil {
			return fmt.Errorf("peer service: failed to configure device: %w", err)
		}

		if err := ps.storage.Update(ctx, id, peer.ToPersistedPeer()); err != nil {
			return err
		}
	}

	return nil
}

func (ps *PeerService) Disable(ctx context.Context, id uuid.UUID) error {
	persistedPeer, err := ps.storage.Get(ctx, id)
	if err != nil {
		return err
	}

	peer, err := persistedPeer.ToPeer()
	if err != nil {
		return err
	}

	if peer.IsEnabled {
		peer.IsEnabled = false
		peerConfig := wgtypes.PeerConfig{
			PublicKey: peer.PublicKey,
			Remove:    true,
		}

		if err = ps.wgclient.ConfigureDevice(peerConfig); err != nil {
			return fmt.Errorf("peer service: failed to configure device: %w", err)
		}

		if err := ps.storage.Update(ctx, id, peer.ToPersistedPeer()); err != nil {
			return err
		}
	}

	return nil
}

func (ps *PeerService) Setup(
	ctx context.Context,
	deviceName, deviceAddress, deviceEndpoint, peerFolder string,
) error {
	device, err := ps.wgclient.GetDevice()
	if err != nil {
		return err
	}

	// extract and persist device ip+net to later get available ips
	deviceIP, deviceNet, err := net.ParseCIDR(deviceAddress)
	if err != nil {
		return fmt.Errorf("peer service: failed to parse device address %s: %w", deviceAddress, err)
	}

	ps.device = &deviceConfig{
		name: deviceName,
		net:  *deviceNet,
		endpoint: net.UDPAddr{
			IP:   net.ParseIP(deviceEndpoint),
			Port: device.ListenPort,
		},
		publicKey: device.PublicKey,
	}

	peers, err := ps.GetAll(ctx)
	if err != nil {
		return err
	}

	ps.reservedIPs = &reservedIPs{
		ips: make([]net.IP, 0, len(peers)+1),
	}

	// add device ip
	if err := ps.reservedIPs.Add(deviceIP); err != nil {
		return err
	}

	// add all the peer ips
	for _, peer := range peers {
		ps.device.incPeerCount()

		for _, addr := range peer.AllowedIPs {
			if err := ps.reservedIPs.Add(addr.IP); err != nil {
				return err
			}
		}
	}

	return nil
}
