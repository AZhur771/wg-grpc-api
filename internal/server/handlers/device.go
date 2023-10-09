package handlers

import (
	"context"
	"database/sql"
	"errors"
	"time"

	wgpb "github.com/AZhur771/wg-grpc-api/gen"
	"github.com/AZhur771/wg-grpc-api/internal/app"
	"github.com/AZhur771/wg-grpc-api/internal/dto"
	"github.com/AZhur771/wg-grpc-api/internal/service/common"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	fieldmask_utils "github.com/mennanov/fieldmask-utils"
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

func NewDeviceImpl(ctx context.Context, logger *zap.Logger, service app.DeviceService) wgpb.DeviceServiceServer {
	return &DeviceImpl{
		Ctx:     ctx,
		Logger:  logger,
		Service: service,
	}
}

func (d *DeviceImpl) Add(ctx context.Context, req *wgpb.AddDeviceRequest) (*wgpb.EntityIdRequest, error) {
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

	errInvalidData := &common.ErrInvalidData{}

	if errors.As(err, errInvalidData) {
		st := status.New(codes.InvalidArgument, err.Error())
		st, err = st.WithDetails(errInvalidData.Details())
		if err != nil {
			return nil, err
		}
		return nil, st.Err()
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, err
	}

	return &wgpb.EntityIdRequest{
		Id: dev.ID.String(),
	}, nil
}

func (d *DeviceImpl) Update(ctx context.Context, req *wgpb.UpdateDeviceRequest) (*emptypb.Empty, error) {
	device := req.GetDevice()

	fmask, err := fieldmask_utils.MaskFromPaths(req.FieldMask.Paths, mapNames)
	if err != nil {
		return nil, err
	}

	ID, err := uuid.Parse(device.GetId())
	if err != nil {
		return nil, err
	}

	_, err = d.Service.Update(ctx,
		dto.UpdateDeviceDTO{
			ID:                  ID,
			Description:         device.GetDescription(),
			Endpoint:            device.GetEndpoint(),
			FirewallMark:        int(device.GetFirewallMark()),
			Address:             device.GetAddress(),
			Table:               device.GetTable(),
			MTU:                 int(device.GetMtu()),
			DNS:                 device.GetDns(),
			PersistentKeepAlive: time.Duration(device.GetPersistentKeepAlive()) * time.Second,
			PreUp:               device.GetPreUp(),
			PreDown:             device.GetPreDown(),
			PostUp:              device.GetPostUp(),
			PostDown:            device.GetPostDown(),
		},
		fmask,
	)

	errInvalidData := &common.ErrInvalidData{}

	if errors.As(err, errInvalidData) {
		st := status.New(codes.InvalidArgument, err.Error())
		st, err = st.WithDetails(errInvalidData.Details())
		if err != nil {
			return nil, err
		}
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

func (d *DeviceImpl) Remove(ctx context.Context, req *wgpb.EntityIdRequest) (*emptypb.Empty, error) {
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

func (d *DeviceImpl) Get(ctx context.Context, req *wgpb.EntityIdRequest) (*wgpb.Device, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	device, err := d.Service.Get(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, err
	}

	return mapEntityDeviceToPbDeivce(device), nil
}

func (d *DeviceImpl) GetAll(ctx context.Context, req *wgpb.GetDevicesRequest) (*wgpb.GetDevicesResponse, error) {
	resp, err := d.Service.GetAll(ctx, dto.GetDevicesRequestDTO{
		Skip:   int(req.GetSkip()),
		Limit:  int(req.GetLimit()),
		Search: req.GetSearch(),
	})

	errInvalidData := &common.ErrInvalidData{}

	if errors.As(err, errInvalidData) {
		st := status.New(codes.InvalidArgument, err.Error())
		st, err = st.WithDetails(errInvalidData.Details())
		if err != nil {
			return nil, err
		}
		return nil, st.Err()
	}

	if err != nil {
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
