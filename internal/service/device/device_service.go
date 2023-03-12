package deviceservice

import (
	"bytes"
	"context"
	"fmt"
	"net"

	"github.com/AZhur771/wg-grpc-api/internal/app"
	"github.com/AZhur771/wg-grpc-api/internal/entity"
	"go.uber.org/zap"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type DeviceService struct {
	logger  app.Logger
	ctrl    app.WgCtrl
	storage app.PeerStorage

	device *entity.Device
}

func NewDeviceService(logger app.Logger, ctrl app.WgCtrl, storage app.PeerStorage) *DeviceService {
	return &DeviceService{
		logger:  logger,
		ctrl:    ctrl,
		storage: storage,
	}
}

//nolint:gocognit
func (ds *DeviceService) syncPeers(ctx context.Context) error {
	peers, err := ds.storage.GetAll(ctx, 0, 0)
	if err != nil {
		return fmt.Errorf("peer service: %w", err)
	}

	wgPeers, err := ds.GetPeers()
	if err != nil {
		return fmt.Errorf("peer service: %w", err)
	}

	for _, wgPeer := range wgPeers {
		for _, allowedIP := range wgPeer.AllowedIPs {
			ds.device.ReserveIP(allowedIP.IP)
		}
	}

	// sync peers with wg
	for _, peer := range peers {
		if err := peer.IsValid(); err != nil {
			return fmt.Errorf("peer service: %w", err)
		}

		peerInSync := false
		for _, wgPeer := range wgPeers {
			if peer.PublicKey == wgPeer.PublicKey {
				peerInSync = true
				break
			}
		}

		//nolint:nestif
		if !peerInSync {
			validPeer := true
			tmpAllowedIps := make([]net.IP, 0, len(peer.AllowedIPs))
			for _, allowedIP := range peer.AllowedIPs {
				netip, _, err := net.ParseCIDR(allowedIP)
				if err != nil {
					return fmt.Errorf("peer service: %w", err)
				}

				if ds.device.IsReservedIP(netip) {
					tmpAllowedIps = append(tmpAllowedIps, netip)
					validPeer = false
					break
				}
			}

			if validPeer {
				if err := ds.device.ReserveManyIPs(tmpAllowedIps); err != nil {
					return fmt.Errorf("peer service: %w", err)
				}

				peerConfig, err := peer.ToPeerConfig()
				if err != nil {
					return fmt.Errorf("peer service: %w", err)
				}

				if err := ds.AddPeer(*peerConfig); err != nil {
					return fmt.Errorf("peer service: %w", err)
				}
			} else {
				ds.logger.Error("Invalid peer", zap.String("ID", peer.ID.String()), zap.Error(ErrInvalidPeer))
			}
		}
	}

	return nil
}

func (ds *DeviceService) Setup(ctx context.Context, name string, endpoint string, address string) error {
	wgdevice, err := ds.ctrl.Device(name)
	if err != nil {
		return fmt.Errorf("device service: %w", err)
	}

	ds.device = &entity.Device{
		Name:         name,
		Endpoint:     endpoint,
		Address:      address,
		Type:         wgdevice.Type,
		PrivateKey:   wgdevice.PrivateKey,
		PublicKey:    wgdevice.PublicKey,
		ListenPort:   wgdevice.ListenPort,
		FirewallMark: wgdevice.FirewallMark,
	}

	if err = ds.device.IsValid(); err != nil {
		return fmt.Errorf("device service: %w", err)
	}

	if err = ds.device.ComputeMaxPeersCount(); err != nil {
		return fmt.Errorf("device service: %w", err)
	}

	if err = ds.device.ComputeInitialReservedIPs(); err != nil {
		return fmt.Errorf("device service: %w", err)
	}

	// sync peers with wg
	return ds.syncPeers(ctx)
}

func (ds *DeviceService) GetDevice() (*entity.Device, error) {
	wgdevice, err := ds.ctrl.Device(ds.device.Name)
	if err != nil {
		return nil, fmt.Errorf("device service: %w", err)
	}

	return ds.device.PopulateDynamicFields(wgdevice), nil
}

func (ds *DeviceService) ConfigureDevice(config wgtypes.PeerConfig) error {
	return ds.ctrl.ConfigureDevice(
		ds.device.Name,
		wgtypes.Config{
			Peers: []wgtypes.PeerConfig{config},
		},
	)
}

func (ds *DeviceService) AddPeer(config wgtypes.PeerConfig) error {
	return ds.ConfigureDevice(config)
}

func (ds *DeviceService) UpdatePeer(config wgtypes.PeerConfig) error {
	config.UpdateOnly = true
	return ds.ConfigureDevice(config)
}

func (ds *DeviceService) RemovePeer(config wgtypes.PeerConfig) error {
	config.Remove = true
	return ds.ConfigureDevice(config)
}

func (ds *DeviceService) GetPeer(publicKey wgtypes.Key) (wgtypes.Peer, error) {
	device, err := ds.ctrl.Device(ds.device.Name)
	if err != nil {
		return wgtypes.Peer{}, fmt.Errorf("device service: %w", err)
	}

	for _, p := range device.Peers {
		if bytes.Equal(p.PublicKey[:], publicKey[:]) {
			return p, nil
		}
	}

	return wgtypes.Peer{}, fmt.Errorf("wg: %w", ErrPeerNotConfigured)
}

func (ds *DeviceService) GetPeers() ([]wgtypes.Peer, error) {
	device, err := ds.ctrl.Device(ds.device.Name)
	if err != nil {
		return nil, fmt.Errorf("device service: %w", err)
	}

	return device.Peers, nil
}
