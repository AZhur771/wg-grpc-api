package peerservice

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/AZhur771/wg-grpc-api/internal/app"
	"github.com/AZhur771/wg-grpc-api/internal/dto"
	"github.com/AZhur771/wg-grpc-api/internal/entity"
	tmpl "github.com/AZhur771/wg-grpc-api/internal/template"
	"github.com/google/uuid"
	qrcode "github.com/skip2/go-qrcode"
	"go.uber.org/zap"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

const defaultLimit = 20

type PeerService struct {
	logger        app.Logger
	deviceService app.DeviceService
	storage       app.PeerStorage
}

func NewPeerService(logger app.Logger, deviceService app.DeviceService, storage app.PeerStorage) *PeerService {
	return &PeerService{
		logger:        logger,
		deviceService: deviceService,
		storage:       storage,
	}
}

func (ps *PeerService) Add(ctx context.Context, addPeerDTO dto.AddPeerDTO) (*entity.Peer, error) {
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	publicKey := privateKey.PublicKey()

	device, err := ps.deviceService.GetDevice()
	if err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	allowedIPNet, err := device.GetAvailableIP()
	if err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	peer := &entity.Peer{
		PrivateKey:                  entity.WgKey(privateKey),
		PublicKey:                   entity.WgKey(publicKey),
		PersistentKeepaliveInterval: addPeerDTO.PersistentKeepAlive,
		AllowedIPs:                  []string{allowedIPNet.String()},
		Tags:                        addPeerDTO.Tags,
		Description:                 addPeerDTO.Description,
		Name:                        addPeerDTO.Name,
		Email:                       addPeerDTO.Email,
		IsEnabled:                   true, // enabled by default
	}

	if addPeerDTO.AddPresharedKey {
		peer.HasPresharedKey = true
		presharedKey, err := wgtypes.GenerateKey()
		if err != nil {
			return nil, fmt.Errorf("peer service: %w", err)
		}
		peer.PresharedKey = entity.WgKey(presharedKey)
	}

	if err := peer.IsValid(); err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	peerConfig, err := peer.ToPeerConfig()
	if err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	if err := ps.deviceService.AddPeer(*peerConfig); err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	peer, err = ps.storage.Add(ctx, peer)
	if err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	wgpeer, err := ps.deviceService.GetPeer(publicKey)
	if err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	return peer.PopulateDynamicFields(&wgpeer), nil
}

func (ps *PeerService) Update(ctx context.Context, updatePeerDTO dto.UpdatePeerDTO) (*entity.Peer, error) {
	peer, err := ps.storage.Get(ctx, updatePeerDTO.ID)
	if err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	peer.Name = updatePeerDTO.Name
	peer.Email = updatePeerDTO.Email
	peer.Tags = updatePeerDTO.Tags
	peer.Description = updatePeerDTO.Description
	peer.PersistentKeepaliveInterval = updatePeerDTO.PersistentKeepAlive

	if peer.HasPresharedKey && updatePeerDTO.RemovePresharedKey {
		peer.HasPresharedKey = false
		peer.PresharedKey = entity.WgKey(wgtypes.Key{}) // non-nil zero value key clears preshared key
	}

	if !peer.HasPresharedKey && updatePeerDTO.AddPresharedKey {
		peer.HasPresharedKey = true
		presharedKey, err := wgtypes.GenerateKey()
		if err != nil {
			return nil, fmt.Errorf("peer service: %w", err)
		}
		peer.PresharedKey = entity.WgKey(presharedKey)
	}

	if err := peer.IsValid(); err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	peerConfig, err := peer.ToPeerConfig()
	if err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	if err := ps.deviceService.UpdatePeer(*peerConfig); err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	peer, err = ps.storage.Update(ctx, peer)
	if err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	wgpeer, err := ps.deviceService.GetPeer(wgtypes.Key(peer.PublicKey))
	if err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	return peer.PopulateDynamicFields(&wgpeer), nil
}

func (ps *PeerService) Remove(ctx context.Context, id uuid.UUID) error {
	peer, err := ps.storage.Remove(ctx, id)
	if err != nil {
		return fmt.Errorf("peer service: %w", err)
	}

	if peer.IsEnabled {
		peerConfig, err := peer.ToPeerConfig()
		if err != nil {
			return fmt.Errorf("peer service: %w", err)
		}

		if err := ps.deviceService.RemovePeer(*peerConfig); err != nil {
			return fmt.Errorf("peer service: %w", err)
		}
	}

	return nil
}

func (ps *PeerService) Get(ctx context.Context, id uuid.UUID) (*entity.Peer, error) {
	peer, err := ps.storage.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	if peer.IsEnabled {
		wgpeer, err := ps.deviceService.GetPeer(wgtypes.Key(peer.PublicKey))
		if err != nil {
			return nil, fmt.Errorf("peer service: %w", err)
		}

		peer = peer.PopulateDynamicFields(&wgpeer)
	}

	return peer, nil
}

func (ps *PeerService) GetAll(ctx context.Context, getPeersRequest dto.GetPeersRequestDTO) (dto.GetPeersResponseDTO, error) {
	getPeersResponse := dto.GetPeersResponseDTO{}

	if !getPeersRequest.IsValid() {
		return getPeersResponse, fmt.Errorf("peer service: %w", ErrInvalidPaginationParams)
	}

	totalPeers, err := ps.storage.CountAll(ctx)
	if err != nil {
		return getPeersResponse, fmt.Errorf("peer service: %w", err)
	}

	if getPeersRequest.Limit == 0 {
		getPeersRequest.Limit = defaultLimit // Default limit
	}

	peers, err := ps.storage.GetAll(ctx, getPeersRequest.Skip, getPeersRequest.Limit)
	if err != nil {
		return getPeersResponse, fmt.Errorf("peer service: %w", err)
	}

	wgPeers, err := ps.deviceService.GetPeers()
	if err != nil {
		return getPeersResponse, fmt.Errorf("peer service: %w", err)
	}

	wgPeersMap := make(map[wgtypes.Key]wgtypes.Peer)
	for _, wgPeer := range wgPeers {
		wgPeersMap[wgPeer.PublicKey] = wgPeer
	}

	for idx, peer := range peers {
		if peer.IsEnabled {
			wgPeer, ok := wgPeersMap[wgtypes.Key(peer.PublicKey)]
			if ok {
				peers[idx] = peer.PopulateDynamicFields(&wgPeer)
			} else {
				ps.logger.Error("Invalid peer", zap.String("ID", peer.ID.String()), zap.Error(ErrInvalidPeer))
			}
		}
	}

	getPeersResponse.Total = totalPeers
	getPeersResponse.Peers = peers
	getPeersResponse.HasNext = (getPeersRequest.Skip + getPeersRequest.Limit) < totalPeers

	return getPeersResponse, nil
}

func (ps *PeerService) Enable(ctx context.Context, id uuid.UUID) error {
	peer, err := ps.storage.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("peer service: %w", err)
	}

	if !peer.IsEnabled {
		peerConfig, err := peer.ToPeerConfig()
		if err != nil {
			return fmt.Errorf("peer service: %w", err)
		}

		if err = ps.deviceService.AddPeer(*peerConfig); err != nil {
			return fmt.Errorf("peer service: %w", err)
		}
		peer.IsEnabled = true

		if _, err := ps.storage.Update(ctx, peer); err != nil {
			return fmt.Errorf("peer service: %w", err)
		}
	}

	return nil
}

func (ps *PeerService) Disable(ctx context.Context, id uuid.UUID) error {
	peer, err := ps.storage.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("peer service: %w", err)
	}

	if peer.IsEnabled {
		peerConfig, err := peer.ToPeerConfig()
		if err != nil {
			return fmt.Errorf("peer service: %w", err)
		}

		if err = ps.deviceService.RemovePeer(*peerConfig); err != nil {
			return fmt.Errorf("peer service: %w", err)
		}
		peer.IsEnabled = false

		if _, err := ps.storage.Update(ctx, peer); err != nil {
			return fmt.Errorf("peer service: %w", err)
		}
	}

	return nil
}

func (ps *PeerService) DownloadConfig(ctx context.Context, id uuid.UUID) (dto.DownloadFileDTO, error) {
	downloadFileDTO := dto.DownloadFileDTO{
		Name: fmt.Sprintf("%s.conf", id.String()),
	}

	device, err := ps.deviceService.GetDevice()
	if err != nil {
		return downloadFileDTO, fmt.Errorf("peer service: %w", err)
	}

	peer, err := ps.storage.Get(ctx, id)
	if err != nil {
		return downloadFileDTO, fmt.Errorf("peer service: %w", err)
	}

	t, err := template.New("config").Funcs(
		template.FuncMap{
			"StringsJoin": strings.Join,
		},
	).Parse(tmpl.ClientConfigTemplate)
	if err != nil {
		return downloadFileDTO, fmt.Errorf("peer service: %w", err)
	}

	var buf bytes.Buffer

	tmplData := tmpl.TmplData{
		// interface data
		InterfacePrivateKey: peer.PrivateKey.String(),
		InterfaceAddress:    peer.AllowedIPs,
		InterfaceDNS:        []string{"9.9.9.9", "149.112.112.112"}, // TODO: allow configure DNS
		InterfaceMTU:        1384,                                   // TODO: allow configure MTU

		// peer data (device)
		PeerPublicKey:       device.PublicKey.String(),
		PeerEndpoint:        fmt.Sprintf("%s:%d", device.Endpoint, device.ListenPort),
		PeerAllowedIPs:      []string{"0.0.0.0/0"},
		PersistentKeepalive: int(peer.PersistentKeepaliveInterval.Seconds()),
	}

	if !peer.PresharedKey.IsEmpty() {
		tmplData.PeerPresharedKey = peer.PresharedKey.String()
	}

	if err := t.Execute(&buf, tmplData); err != nil {
		return downloadFileDTO, fmt.Errorf("peer service: %w", err)
	}

	downloadFileDTO.Data = buf.Bytes()
	downloadFileDTO.Size = int64(buf.Len())

	return downloadFileDTO, nil
}

func (ps *PeerService) DownloadQRCode(ctx context.Context, id uuid.UUID) (dto.DownloadFileDTO, error) {
	downloadFileDTO := dto.DownloadFileDTO{
		Name: fmt.Sprintf("%s.png", id.String()),
	}

	config, err := ps.DownloadConfig(ctx, id)
	if err != nil {
		return downloadFileDTO, err
	}

	png, err := qrcode.Encode(string(config.Data), qrcode.Medium, 256)
	if err != nil {
		return downloadFileDTO, fmt.Errorf("peer service: %w", err)
	}

	downloadFileDTO.Data = png
	downloadFileDTO.Size = int64(len(png))

	return downloadFileDTO, err
}
