package deviceservice

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/AZhur771/wg-grpc-api/internal/app"
	dt "github.com/AZhur771/wg-grpc-api/internal/dto"
	"github.com/AZhur771/wg-grpc-api/internal/entity"
	"github.com/AZhur771/wg-grpc-api/internal/service/common"
	tmpl "github.com/AZhur771/wg-grpc-api/internal/template"
	"github.com/google/uuid"
	fieldmask_utils "github.com/mennanov/fieldmask-utils"
	"go.uber.org/zap"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

const defaultLimit = 20

type DeviceService struct {
	logger     *zap.Logger
	ctrl       app.WgCtrl
	deviceRepo app.DeviceRepo
	peerRepo   app.PeerRepo
}

func NewDeviceService(logger *zap.Logger, ctrl app.WgCtrl, deviceRepo app.DeviceRepo, peerRepo app.PeerRepo) *DeviceService {
	return &DeviceService{
		logger:     logger,
		ctrl:       ctrl,
		deviceRepo: deviceRepo,
		peerRepo:   peerRepo,
	}
}

//nolint:gocognit
func (ds *DeviceService) SyncDevices(ctx context.Context) error {
	devices, err := ds.deviceRepo.GetAll(ctx, nil, 0, 0, "")
	if err != nil {
		return fmt.Errorf("device service: %w", err)
	}

	for _, device := range devices {
		if err := ds.setupDevice(device); err != nil {
			return fmt.Errorf("device service: %w", err)
		}

		peers, err := ds.peerRepo.GetAll(ctx, nil, 0, 0, "", device.ID)
		if err != nil {
			return fmt.Errorf("device service: %w", err)
		}

		wgPeers, err := ds.GetConfiguredPeers(device.Name)
		if err != nil {
			return fmt.Errorf("device service: %w", err)
		}

		// sync peers with wg
		for _, peer := range peers {
			if errors := peer.IsValid(); len(errors) > 0 {
				return common.NewErrInvalidData(fmt.Errorf("device service: invalid peer data"), errors)
			}

			peerInSync := false
			for _, wgPeer := range wgPeers {
				if peer.PublicKey == wgPeer.PublicKey {
					peerInSync = true
					break
				}
			}

			if !peerInSync {
				peerConfig, err := peer.ToPeerConfig()
				if err != nil {
					return fmt.Errorf("device service: %w", err)
				}

				if err := ds.ConfigureDevice(device.Name, *peerConfig); err != nil {
					return fmt.Errorf("device service: %w", err)
				}
			}
		}
	}

	return nil
}

func (ds *DeviceService) restartDevice(devName string) error {
	//nolint:staticcheck
	if err := exec.Command("wg-quick", "down", devName).Run(); err != nil {
		// ignore
	}

	return exec.Command("wg-quick", "up", devName).Run()
}

func (ds *DeviceService) setupDevice(dev *entity.Device) error {
	filename := fmt.Sprintf("/etc/wireguard/%s.conf", dev.Name)

	t, err := template.New("config").Parse(tmpl.ServerConfigTemplate)
	if err != nil {
		return fmt.Errorf("device service: %w", err)
	}

	privateKey, err := exec.Command("wg", "genkey").Output()
	if err != nil {
		return fmt.Errorf("device service: %w", err)
	}

	port := strings.Split(dev.Endpoint, ":")[1]

	tmplData := tmpl.ServerConfigTmplData{
		InterfacePrivateKey:          string(privateKey),
		InterfaceAddress:             dev.Address,
		InterfacePort:                port,
		InterfaceMTU:                 dev.MTU,
		InterfaceTable:               dev.Table,
		InterfaceDNS:                 dev.DNS,
		InterfaceFwMark:              dev.FirewallMark,
		InterfacePreUp:               dev.PreUp,
		InterfacePostUp:              dev.PostUp,
		InterfacePreDown:             dev.PreDown,
		InterfacePostDown:            dev.PostDown,
		InterfacePersistentKeepAlive: int(dev.PersistentKeepAlive) / (1000 * 1000),
	}

	var buf bytes.Buffer

	if err := t.Execute(&buf, tmplData); err != nil {
		return fmt.Errorf("setup: %w", err)
	}

	if b, err := ioutil.ReadFile(filename); err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(b))

		writeLine := false

		for scanner.Scan() {
			l := scanner.Text()

			if l == "[Peer]" {
				writeLine = true
			}

			if writeLine {
				buf.WriteString(l)
			}
		}
	} else {
		return fmt.Errorf("device service: %w", err)
	}

	if err := ioutil.WriteFile(filename, buf.Bytes(), 0o600); err != nil {
		return fmt.Errorf("device service: %w", err)
	}

	if err := ds.restartDevice(dev.Name); err != nil {
		return err
	}

	return nil
}

func (ds *DeviceService) Add(ctx context.Context, dto dt.AddDeviceDTO) (*entity.Device, error) {
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, fmt.Errorf("device service: %w", err)
	}

	dev := &entity.Device{
		PrivateKey:          privateKey,
		PublicKey:           privateKey.PublicKey(),
		Description:         dto.Description,
		Endpoint:            dto.Endpoint,
		Address:             dto.Address,
		FirewallMark:        dto.FirewallMark,
		PersistentKeepAlive: dto.PersistentKeepAlive,
		MTU:                 dto.MTU,
		DNS:                 dto.DNS,
		Table:               dto.Table,
		PostUp:              dto.PostUp,
		PostDown:            dto.PostDown,
		PreUp:               dto.PreUp,
		PreDown:             dto.PreDown,
	}

	if errors := dev.IsValid(); len(errors) > 0 {
		return nil, common.NewErrInvalidData(fmt.Errorf("device service: %w", ErrInvalidDeviceData), errors)
	}

	dev, err = ds.deviceRepo.Add(ctx, nil, dev)
	if err != nil {
		return nil, fmt.Errorf("device service: %w", err)
	}

	if err := ds.setupDevice(dev); err != nil {
		return nil, fmt.Errorf("device service: %w", err)
	}

	if err := ds.restartDevice(dev.Name); err != nil {
		return nil, fmt.Errorf("device service: %w", err)
	}

	wgdev, err := ds.ctrl.Device(dev.Name)
	if err != nil {
		return nil, fmt.Errorf("device service: %w", err)
	}

	return dev.PopulateDynamicFields(wgdev)
}

func (ds *DeviceService) Update(ctx context.Context, dto dt.UpdateDeviceDTO, mask fieldmask_utils.Mask) (*entity.Device, error) {
	dev, err := ds.deviceRepo.Get(ctx, nil, dto.ID)
	if err != nil {
		return nil, fmt.Errorf("device service: %w", err)
	}

	fieldmask_utils.StructToStruct(mask, dto, dev)

	if errors := dev.IsValid(); len(errors) > 0 {
		return nil, common.NewErrInvalidData(fmt.Errorf("device service: %w", ErrInvalidDeviceData), errors)
	}

	if err := ds.setupDevice(dev); err != nil {
		return nil, fmt.Errorf("device service: %w", err)
	}

	if err := ds.restartDevice(dev.Name); err != nil {
		return nil, fmt.Errorf("device service: %w", err)
	}

	dev, err = ds.deviceRepo.Update(ctx, nil, dev)
	if err != nil {
		return nil, fmt.Errorf("device service: %w", err)
	}

	wgdev, err := ds.ctrl.Device(dev.Name)
	if err != nil {
		return nil, fmt.Errorf("device service: %w", err)
	}

	return dev.PopulateDynamicFields(wgdev)
}

func (ds *DeviceService) Remove(ctx context.Context, id uuid.UUID) error {
	dev, err := ds.deviceRepo.Get(ctx, nil, id)
	if err != nil {
		return fmt.Errorf("device service: %w", err)
	}

	if err := ds.deviceRepo.Remove(ctx, nil, id); err != nil {
		return fmt.Errorf("device service: %w", err)
	}

	//nolint:gosec
	if err := exec.Command("wg-quick", "down", dev.Name).Run(); err != nil {
		ds.logger.Error(fmt.Sprintf("'wg-quick down %s' errored", dev.Name), zap.Error(err))
	}

	filename := fmt.Sprintf("/etc/wireguard/%s.conf", dev.Name)
	if err := os.Remove(filename); err != nil {
		return err
	}

	return nil
}

func (ds *DeviceService) Get(ctx context.Context, id uuid.UUID) (*entity.Device, error) {
	dev, err := ds.deviceRepo.Get(ctx, nil, id)
	if err != nil {
		return nil, fmt.Errorf("device service: %w", err)
	}

	wgdev, err := ds.ctrl.Device(dev.Name)
	if err != nil {
		return nil, fmt.Errorf("device service: %w", err)
	}

	return dev.PopulateDynamicFields(wgdev)
}

func (ds *DeviceService) GetAll(ctx context.Context, dto dt.GetDevicesRequestDTO) (dt.GetDevicesResponseDTO, error) {
	resp := dt.GetDevicesResponseDTO{}

	if errors := dto.IsValid(); len(errors) > 0 {
		return resp, common.NewErrInvalidData(fmt.Errorf("device service: %w", ErrInvalidPaginationParams), errors)
	}

	total, err := ds.deviceRepo.Count(ctx, nil)
	if err != nil {
		return resp, fmt.Errorf("device service: %w", err)
	}

	if dto.Limit == 0 {
		dto.Limit = defaultLimit
	}

	devices, err := ds.deviceRepo.GetAll(ctx, nil, dto.Skip, dto.Limit, dto.Search)
	if err != nil {
		return resp, fmt.Errorf("device service: %w", err)
	}

	for i, dev := range devices {
		wgdev, err := ds.ctrl.Device(dev.Name)
		if err != nil {
			ds.logger.Error("failed to get device", zap.Error(err))
		}

		if devices[i], err = dev.PopulateDynamicFields(wgdev); err != nil {
			return resp, fmt.Errorf("device service: %w", err)
		}
	}

	resp.Total = total
	resp.Devices = devices
	resp.HasNext = (dto.Skip + dto.Limit) < total

	return resp, nil
}

func (ds *DeviceService) ConfigureDevice(device string, config wgtypes.PeerConfig) error {
	return ds.ctrl.ConfigureDevice(
		device,
		wgtypes.Config{
			Peers: []wgtypes.PeerConfig{config},
		},
	)
}

func (ds *DeviceService) GetConfiguredPeer(dev string, publicKey wgtypes.Key) (wgtypes.Peer, error) {
	device, err := ds.ctrl.Device(dev)
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

func (ds *DeviceService) GetConfiguredPeers(dev string) ([]wgtypes.Peer, error) {
	device, err := ds.ctrl.Device(dev)
	if err != nil {
		return nil, fmt.Errorf("device service: %w", err)
	}

	return device.Peers, nil
}
