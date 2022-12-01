package service

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"strings"
	"text/template"
	"time"

	"github.com/AZhur771/wg-grpc-api/internal/entity"
	"github.com/AZhur771/wg-grpc-api/internal/storage"
	tmpl "github.com/AZhur771/wg-grpc-api/internal/template"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type peerService struct {
	logger      *zap.Logger
	wgclient    *wgctrl.Client
	reservedIPs *reservedIPs
	device      *deviceConfig
	storage     storage.PeerStorage
}

func NewPeerService(logger *zap.Logger, wgclient *wgctrl.Client, storage storage.PeerStorage) PeerService {
	return &peerService{
		logger:   logger,
		wgclient: wgclient,
		storage:  storage,
	}
}

func (ps *peerService) getAvailableIP(ipnet *net.IPNet) (*net.IPNet, error) {
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

func (ps *peerService) configureDevice(peerConfig wgtypes.PeerConfig) error {
	return ps.wgclient.ConfigureDevice(
		ps.device.name,
		wgtypes.Config{
			Peers: []wgtypes.PeerConfig{peerConfig},
		},
	)
}

func (ps *peerService) Add(
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

	if err := ps.configureDevice(peerConfig); err != nil {
		return id, fmt.Errorf("peer service: failed to configure device: %w", err)
	}

	id, err = ps.storage.Create(ctx, peer.ToPersistedPeer())
	if err != nil {
		return id, err
	}

	ps.device.incPeerCount()

	return id, nil
}

func (ps *peerService) Update(
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

	if containsString(updateMask, "add_preshared_key") && addPresharedKey != peer.HasPresharedKey {
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

		if err := ps.configureDevice(peerConfig); err != nil {
			return fmt.Errorf("peer service: failed to configure device: %w", err)
		}
	}

	return ps.storage.Update(ctx, id, peer.ToPersistedPeer())
}

func (ps *peerService) Delete(ctx context.Context, id uuid.UUID) error {
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

	if err = ps.configureDevice(peerConfig); err != nil {
		return fmt.Errorf("peer service: failed to configure device: %w", err)
	}

	ps.device.decPeerCount()

	return nil
}

func (ps *peerService) Get(ctx context.Context, id uuid.UUID) (*entity.Peer, error) {
	device, err := ps.wgclient.Device(ps.device.name)
	if err != nil {
		return nil, fmt.Errorf("peer service: failed to get device: %w", err)
	}

	persistedPeer, err := ps.storage.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	peer, err := persistedPeer.ToPeer()
	if err != nil {
		return nil, err
	}

	for i, p := range device.Peers {
		if bytes.Equal(p.PublicKey[:], peer.PublicKey[:]) {
			return peer.PopulateDynamicFields(&device.Peers[i]), nil
		}
	}

	return peer, nil
}

func (ps *peerService) GetAll(ctx context.Context, limit, skip int) (*entity.PaginatedPeers, error) {
	device, err := ps.wgclient.Device(ps.device.name)
	if err != nil {
		return nil, fmt.Errorf("peer service: failed to get device: %w", err)
	}

	persistedPeers, err := ps.storage.GetAll(ctx, limit, skip)
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

		for i, p := range device.Peers {
			if bytes.Equal(p.PublicKey[:], peer.PublicKey[:]) {
				peers = append(peers, peer.PopulateDynamicFields(&device.Peers[i]))
				break
			}
		}
	}

	return &entity.PaginatedPeers{
		Peers:   peers,
		Total:   ps.device.getPeerCount(),
		HasNext: ps.device.getPeerCount() > len(peers)+limit+skip,
	}, nil
}

func (ps *peerService) DownloadConfig(ctx context.Context, id uuid.UUID) ([]byte, error) {
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

	if err := t.Execute(&buf, tmpl.TmplData{
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
	}); err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

func (ps *peerService) Setup(
	ctx context.Context,
	deviceName, deviceAddress, deviceEndpoint, peerFolder string,
) error {
	device, err := ps.wgclient.Device(deviceName)
	if err != nil {
		return fmt.Errorf("peer service: failed to get device: %w", err)
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

	paginatedPeers, err := ps.GetAll(ctx, 0, 0)
	if err != nil {
		return err
	}

	ps.reservedIPs = &reservedIPs{
		ips: make([]net.IP, 0, len(paginatedPeers.Peers)+1),
	}

	// add device ip
	if err := ps.reservedIPs.Add(deviceIP); err != nil {
		return err
	}

	// add all the peer ips
	for _, peer := range paginatedPeers.Peers {
		ps.device.incPeerCount()

		for _, addr := range peer.AllowedIPs {
			if err := ps.reservedIPs.Add(addr.IP); err != nil {
				return err
			}
		}
	}

	return nil
}
