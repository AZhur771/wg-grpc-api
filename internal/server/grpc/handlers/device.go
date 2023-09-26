package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	wgpb "github.com/AZhur771/wg-grpc-api/gen"
	"github.com/AZhur771/wg-grpc-api/internal/app"
	"github.com/AZhur771/wg-grpc-api/internal/dto"
	"github.com/AZhur771/wg-grpc-api/internal/entity"
	deviceservice "github.com/AZhur771/wg-grpc-api/internal/service/device"
	peerservice "github.com/AZhur771/wg-grpc-api/internal/service/peer"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DeviceImpl struct {
	Ctx     context.Context
	Logger  *zap.Logger
	Service app.DeviceService

	wgpb.UnimplementedDeviceServiceServer
}

func NewDeviceImpl(ctx context.Context, logger *zap.Logger, service app.DeviceService) *DeviceImpl {
	return &DeviceImpl{
		Ctx:     ctx,
		Logger:  logger,
		Service: service,
	}
}

func (d *DeviceImpl) AddDevice(ctx context.Context, req *wgpb.AddDeviceRequest) (*wgpb.EntityIdRequest, error) {
	dev, err := d.Service.Add(ctx,
		dto.AddDeviceDTO{
			Description:         req.GetDescription(),
			Endpoint:            req.GetEndpoint(),
			FirewallMark:        int(req.GetFirewallMark()),
			Address:             req.GetAddress(),
			Table:               req.GetTable(),
			MTU:                 int(req.GetMtu()),
			DNS:                 req.GetDns(),
			PersistentKeepAlive: time.Duration(req.GetPersistentKeepAlive()) * time.Second,
			PreUp:               req.GetPreUp(),
			PreDown:             req.GetPreDown(),
			PostUp:              req.GetPostUp(),
			PostDown:            req.GetPostDown(),
		},
	)

	errInvalidDevice := &deviceservice.ErrInvalidDevice{}

	if errors.As(err, errInvalidDevice) {
		st := status.New(codes.InvalidArgument, err.Error())
		st.WithDetails(errInvalidDevice.Details())
		return nil, st.Err()
	} else if err != nil {
		return nil, err
	}

	return &wgpb.EntityIdRequest{
		Id: dev.ID.String(),
	}, nil
}

func (d *DeviceImpl) UpdateDevice(ctx context.Context, req *wgpb.UpdateDeviceRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	_, err = d.Service.Update(ctx,
		dto.UpdateDeviceDTO{
			ID:                  id,
			Description:         req.GetDescription(),
			Endpoint:            req.GetEndpoint(),
			FirewallMark:        int(req.GetFirewallMark()),
			Address:             req.GetAddress(),
			Table:               req.GetTable(),
			MTU:                 int(req.GetMtu()),
			DNS:                 req.GetDns(),
			PersistentKeepAlive: time.Duration(req.GetPersistentKeepAlive()) * time.Second,
			PreUp:               req.GetPreUp(),
			PreDown:             req.GetPreDown(),
			PostUp:              req.GetPostUp(),
			PostDown:            req.GetPostDown(),
		},
	)

	errInvalidDevice := deviceservice.ErrInvalidDevice{}

	if errors.As(err, &errInvalidDevice) {
		st := status.New(codes.InvalidArgument, err.Error())
		st.WithDetails(errInvalidDevice.Details())
		return nil, st.Err()
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (d *DeviceImpl) RemoveDevice(ctx context.Context, req *wgpb.EntityIdRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	if err := d.Service.Remove(ctx, id); errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (d *DeviceImpl) GetDevice(ctx context.Context, req *wgpb.EntityIdRequest) (*wgpb.Device, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	device, err := d.Service.Get(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, err
	}

	return mapEntityDeviceToPbDeivce(device), nil
}

func (d *DeviceImpl) GetDevices(ctx context.Context, req *wgpb.GetDevicesRequest) (*wgpb.GetDevicesResponse, error) {
	resp, err := d.Service.GetAll(ctx, dto.GetDevicesRequestDTO{
		Skip:   int(req.GetSkip()),
		Limit:  int(req.GetLimit()),
		Search: req.GetSearch(),
	})
	if errors.Is(err, peerservice.ErrInvalidPaginationParams) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	} else if err != nil {
		return nil, err
	}

	devicespb := make([]*wgpb.Device, 0, len(resp.Devices))
	for _, dev := range resp.Devices {
		devicespb = append(devicespb, mapEntityDeviceToPbDeivce(dev))
	}

	return &wgpb.GetDevicesResponse{
		Devices: devicespb,
		Total:   int32(resp.Total),
		HasNext: resp.HasNext,
	}, nil
}

func mapEntityDeviceToPbDeivce(dev *entity.Device) *wgpb.Device {
	return &wgpb.Device{
		Id:                  dev.ID.String(),
		Name:                dev.Name,
		Description:         dev.Description,
		Type:                dev.Type.String(),
		PublicKey:           dev.PublicKey.String(),
		FirewallMark:        int32(dev.FirewallMark),
		MaxPeersCount:       int32(dev.MaxPeersCount),
		CurrentPeersCount:   int32(dev.CurrentPeersCount),
		Endpoint:            fmt.Sprintf("%s:%d", dev.Endpoint, dev.ListenPort),
		Address:             dev.Address,
		Mtu:                 int32(dev.MTU),
		Dns:                 dev.DNS,
		Table:               dev.Table,
		PersistentKeepAlive: int32(dev.PersistentKeepAlive),
		PreUp:               dev.PreUp,
		PreDown:             dev.PreDown,
		PostUp:              dev.PostDown,
		PostDown:            dev.PostDown,
	}
}
