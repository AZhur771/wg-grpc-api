package handlers

import (
	"context"

	wgpb "github.com/AZhur771/wg-grpc-api/gen"
	"github.com/AZhur771/wg-grpc-api/internal/app"
	"github.com/AZhur771/wg-grpc-api/internal/entity"
	"github.com/golang/protobuf/ptypes/empty"
)

type DeviceImpl struct {
	Ctx     context.Context
	Logger  app.Logger
	Service app.DeviceService

	wgpb.UnimplementedDeviceServiceServer
}

func NewDeviceImpl(ctx context.Context, logger app.Logger, service app.DeviceService) *DeviceImpl {
	return &DeviceImpl{
		Ctx:     ctx,
		Logger:  logger,
		Service: service,
	}
}

func (d *DeviceImpl) GetDevice(ctx context.Context, getDeviceRequest *empty.Empty) (*wgpb.Device, error) {
	device, err := d.Service.GetDevice()
	if err != nil {
		return nil, err
	}

	return mapEntityDeviceToPbDeivce(device), nil
}

func mapEntityDeviceToPbDeivce(device *entity.Device) *wgpb.Device {
	return &wgpb.Device{
		Name:              device.Name,
		Type:              device.Type.String(),
		PublicKey:         device.PublicKey.String(),
		ListenPort:        int32(device.ListenPort),
		FirewallMark:      int32(device.FirewallMark),
		MaxPeersCount:     int32(device.MaxPeersCount),
		CurrentPeersCount: int32(device.CurrentPeersCount),
		Endpoint:          device.Endpoint,
		Address:           device.Address,
	}
}
