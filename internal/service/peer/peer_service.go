package peerservice

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"strings"
	"text/template"

	"github.com/AZhur771/wg-grpc-api/internal/app"
	dt "github.com/AZhur771/wg-grpc-api/internal/dto"
	"github.com/AZhur771/wg-grpc-api/internal/entity"
	"github.com/AZhur771/wg-grpc-api/internal/service/common"
	tmpl "github.com/AZhur771/wg-grpc-api/internal/template"
	"github.com/google/uuid"
	fieldmask_utils "github.com/mennanov/fieldmask-utils"
	qrcode "github.com/skip2/go-qrcode"
	"go.uber.org/zap"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

const defaultLimit = 20

type PeerService struct {
	logger        *zap.Logger
	deviceService app.DeviceService
	peerRepo      app.PeerRepo
	deviceRepo    app.DeviceRepo
}

func NewPeerService(logger *zap.Logger, deviceService app.DeviceService, deviceRepo app.DeviceRepo, peerRepo app.PeerRepo) *PeerService {
	return &PeerService{
		logger:        logger,
		deviceService: deviceService,
		peerRepo:      peerRepo,
		deviceRepo:    deviceRepo,
	}
}

func (ps *PeerService) Add(ctx context.Context, dto dt.AddPeerDTO) (*entity.Peer, error) {
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	publicKey := privateKey.PublicKey()

	device, err := ps.deviceService.Get(ctx, dto.DeviceID)
	if err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	_, ipnet, err := net.ParseCIDR(device.Address)
	if err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	ones, _ := ipnet.Mask.Size()

	peer := &entity.Peer{
		DeviceID:                    dto.DeviceID,
		PrivateKey:                  privateKey,
		PublicKey:                   publicKey,
		PersistentKeepaliveInterval: dto.PersistentKeepAlive,
		Description:                 dto.Description,
		Name:                        dto.Name,
		Email:                       dto.Email,
		DNS:                         dto.DNS,
		MTU:                         dto.MTU,
		IsEnabled:                   true,
	}

	if peer.DNS == "" {
		peer.DNS = "9.9.9.9, 149.112.112.112"
	}

	if peer.MTU == 0 {
		// https://gist.github.com/nitred/f16850ca48c48c79bf422e90ee5b9d95
		peer.MTU = 1384
	}

	if dto.AddPresharedKey {
		peer.HasPresharedKey = true
		presharedKey, err := wgtypes.GenerateKey()
		if err != nil {
			return nil, fmt.Errorf("peer service: %w", err)
		}
		peer.PresharedKey = presharedKey
	}

	tx, err := ps.peerRepo.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}
	defer tx.Rollback()

	addr, err := ps.deviceRepo.GenerateAddress(ctx, tx, device)
	if err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}
	peer.AllowedIPs = []string{fmt.Sprintf("%s/%d", addr, ones)}

	if errors := peer.IsValid(); len(errors) > 0 {
		return nil, common.NewErrInvalidData(fmt.Errorf("peer service: %w", ErrInvalidPeerData), errors)
	}

	peer, err = ps.peerRepo.Add(ctx, tx, peer)
	if err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	peerConfig, err := peer.ToPeerConfig(device)
	if err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	if err := ps.deviceService.ConfigureDevice(device.Name, *peerConfig); err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	wgpeer, err := ps.deviceService.GetConfiguredPeer(device.Name, peer.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	return peer.PopulateDynamicFields(&wgpeer), nil
}

func (ps *PeerService) Update(ctx context.Context, dto dt.UpdatePeerDTO, mask fieldmask_utils.Mask) (*entity.Peer, error) {
	peer, err := ps.peerRepo.Get(ctx, nil, dto.ID)
	if err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	device, err := ps.deviceService.Get(ctx, peer.DeviceID)
	if err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	fieldmask_utils.StructToStruct(mask, dto, peer)

	if peer.DNS == "" {
		peer.DNS = "9.9.9.9, 149.112.112.112"
	}

	if peer.MTU == 0 {
		// https://gist.github.com/nitred/f16850ca48c48c79bf422e90ee5b9d95
		peer.MTU = 1384
	}

	if peer.HasPresharedKey && dto.RemovePresharedKey {
		peer.HasPresharedKey = false
		peer.PresharedKey = wgtypes.Key{} // non-nil zero value key clears preshared key
	}

	if !peer.HasPresharedKey && dto.AddPresharedKey {
		peer.HasPresharedKey = true
		presharedKey, err := wgtypes.GenerateKey()
		if err != nil {
			return nil, fmt.Errorf("peer service: %w", err)
		}
		peer.PresharedKey = presharedKey
	}

	if errors := peer.IsValid(); len(errors) > 0 {
		return nil, common.NewErrInvalidData(fmt.Errorf("peer service: %w", ErrInvalidPeerData), errors)
	}

	peerConfig, err := peer.ToPeerConfig(device)
	if err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}
	peerConfig.UpdateOnly = true

	if err := ps.deviceService.ConfigureDevice(device.Name, *peerConfig); err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	peer, err = ps.peerRepo.Update(ctx, nil, peer)
	if err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	wgpeer, err := ps.deviceService.GetConfiguredPeer(device.Name, peer.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	return peer.PopulateDynamicFields(&wgpeer), nil
}

func (ps *PeerService) Remove(ctx context.Context, id uuid.UUID) error {
	peer, err := ps.peerRepo.Get(ctx, nil, id)
	if err != nil {
		return fmt.Errorf("peer service: %w", err)
	}

	device, err := ps.deviceService.Get(ctx, peer.DeviceID)
	if err != nil {
		return fmt.Errorf("peer service: %w", err)
	}

	if err := ps.peerRepo.Remove(ctx, nil, id); err != nil {
		return fmt.Errorf("peer service: %w", err)
	}

	if peer.IsEnabled {
		peerConfig, err := peer.ToPeerConfig(device)
		if err != nil {
			return fmt.Errorf("peer service: %w", err)
		}
		peerConfig.Remove = true

		device, err := ps.deviceService.Get(ctx, peer.DeviceID)
		if err != nil {
			return fmt.Errorf("peer service: %w", err)
		}

		if err := ps.deviceService.ConfigureDevice(device.Name, *peerConfig); err != nil {
			return fmt.Errorf("peer service: %w", err)
		}
	}

	return nil
}

func (ps *PeerService) Get(ctx context.Context, id uuid.UUID) (*entity.Peer, error) {
	peer, err := ps.peerRepo.Get(ctx, nil, id)
	if err != nil {
		return nil, fmt.Errorf("peer service: %w", err)
	}

	if peer.IsEnabled {
		device, err := ps.deviceService.Get(ctx, peer.DeviceID)
		if err != nil {
			return nil, fmt.Errorf("peer service: %w", err)
		}

		wgpeer, err := ps.deviceService.GetConfiguredPeer(device.Name, peer.PublicKey)
		if err != nil {
			return nil, fmt.Errorf("peer service: %w", err)
		}

		peer = peer.PopulateDynamicFields(&wgpeer)
	}

	return peer, nil
}

func (ps *PeerService) GetAll(ctx context.Context, dto dt.GetPeersRequestDTO) (dt.GetPeersResponseDTO, error) {
	resp := dt.GetPeersResponseDTO{}

	if errors := dto.IsValid(); len(errors) > 0 {
		return resp, common.NewErrInvalidData(fmt.Errorf("peer service: %w", ErrInvalidPaginationParams), errors)
	}

	total, err := ps.peerRepo.Count(ctx, nil, dto.DeviceID)
	if err != nil {
		return resp, fmt.Errorf("peer service: %w", err)
	}

	if dto.Limit == 0 {
		dto.Limit = defaultLimit
	}

	peers, err := ps.peerRepo.GetAll(ctx, nil, dto.Skip, dto.Limit, dto.Search, dto.DeviceID)
	if err != nil {
		return resp, fmt.Errorf("peer service: %w", err)
	}

	resp.Total = total
	resp.Peers = peers
	resp.HasNext = (dto.Skip + dto.Limit) < total

	return resp, nil
}

func (ps *PeerService) Enable(ctx context.Context, id uuid.UUID) error {
	peer, err := ps.peerRepo.Get(ctx, nil, id)
	if err != nil {
		return fmt.Errorf("peer service: %w", err)
	}

	device, err := ps.deviceService.Get(ctx, peer.DeviceID)
	if err != nil {
		return fmt.Errorf("peer service: %w", err)
	}

	if !peer.IsEnabled {
		peerConfig, err := peer.ToPeerConfig(device)
		if err != nil {
			return fmt.Errorf("peer service: %w", err)
		}

		device, err := ps.deviceService.Get(ctx, peer.DeviceID)
		if err != nil {
			return fmt.Errorf("peer service: %w", err)
		}

		if err = ps.deviceService.ConfigureDevice(device.Name, *peerConfig); err != nil {
			return fmt.Errorf("peer service: %w", err)
		}
		peer.IsEnabled = true

		if _, err := ps.peerRepo.Update(ctx, nil, peer); err != nil {
			return fmt.Errorf("peer service: %w", err)
		}
	}

	return nil
}

func (ps *PeerService) Disable(ctx context.Context, id uuid.UUID) error {
	peer, err := ps.peerRepo.Get(ctx, nil, id)
	if err != nil {
		return fmt.Errorf("peer service: %w", err)
	}

	device, err := ps.deviceService.Get(ctx, peer.DeviceID)
	if err != nil {
		return fmt.Errorf("peer service: %w", err)
	}

	if peer.IsEnabled {
		peerConfig, err := peer.ToPeerConfig(device)
		if err != nil {
			return fmt.Errorf("peer service: %w", err)
		}
		peerConfig.Remove = true

		device, err := ps.deviceService.Get(ctx, peer.DeviceID)
		if err != nil {
			return fmt.Errorf("peer service: %w", err)
		}

		if err = ps.deviceService.ConfigureDevice(device.Name, *peerConfig); err != nil {
			return fmt.Errorf("peer service: %w", err)
		}
		peer.IsEnabled = false

		if _, err := ps.peerRepo.Update(ctx, nil, peer); err != nil {
			return fmt.Errorf("peer service: %w", err)
		}
	}

	return nil
}

func (ps *PeerService) DownloadConfig(ctx context.Context, id uuid.UUID) (dt.DownloadFileDTO, error) {
	downloadFileDTO := dt.DownloadFileDTO{
		Name: fmt.Sprintf("%s.conf", id.String()),
	}

	peer, err := ps.peerRepo.Get(ctx, nil, id)
	if err != nil {
		return downloadFileDTO, fmt.Errorf("peer service: %w", err)
	}

	device, err := ps.deviceService.Get(ctx, peer.DeviceID)
	if err != nil {
		return downloadFileDTO, fmt.Errorf("peer service: %w", err)
	}

	t, err := template.New("config").Funcs(
		template.FuncMap{
			"StringsJoin": strings.Join,
		},
	).Parse(tmpl.ConfigTemplate)
	if err != nil {
		return downloadFileDTO, fmt.Errorf("peer service: %w", err)
	}

	var buf bytes.Buffer

	tmplData := tmpl.ConfigTmplData{
		// interface data
		InterfacePrivateKey: peer.PrivateKey.String(),
		InterfaceAddress:    peer.AllowedIPs,
		InterfaceDNS:        peer.DNS,
		InterfaceMTU:        peer.MTU,

		InterfacePeers: []tmpl.PeerConfigTmplData{
			{
				PeerPublicKey:           device.PublicKey.String(),
				PeerPresharedKey:        peer.PresharedKey.String(),
				PeerEndpoint:            device.Endpoint,
				PeerAllowedIPs:          []string{"0.0.0.0/0"},
				PeerPersistentKeepalive: int(peer.PersistentKeepaliveInterval) / (1000 * 1000 * 1000),
			},
		},
	}

	if err := t.Execute(&buf, tmplData); err != nil {
		return downloadFileDTO, fmt.Errorf("peer service: %w", err)
	}

	downloadFileDTO.Data = buf.Bytes()
	downloadFileDTO.Size = int64(buf.Len())

	return downloadFileDTO, nil
}

func (ps *PeerService) DownloadQRCode(ctx context.Context, id uuid.UUID) (dt.DownloadFileDTO, error) {
	downloadFileDTO := dt.DownloadFileDTO{
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
