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
)

type PeersImpl struct {
	Ctx     context.Context
	Logger  *zap.Logger
	Service app.PeerService

	wgpb.UnimplementedPeerServiceServer
}

func NewPeersImpl(ctx context.Context, logger *zap.Logger, service app.PeerService) *PeersImpl {
	return &PeersImpl{
		Ctx:     ctx,
		Logger:  logger,
		Service: service,
	}
}

func (p *PeersImpl) Add(ctx context.Context, req *wgpb.AddPeerRequest) (*wgpb.EntityIdRequest, error) {
	deviceID, err := uuid.Parse(req.GetDeviceId())
	if err != nil {
		return nil, err
	}

	peer, err := p.Service.Add(ctx,
		dto.AddPeerDTO{
			DeviceID:            deviceID,
			Name:                req.GetName(),
			Email:               req.GetEmail(),
			Description:         req.GetDescription(),
			AddPresharedKey:     req.GetAddPresharedKey(),
			PersistentKeepAlive: time.Duration(req.GetPersistentKeepAlive()) * time.Second,
			MTU:                 int(req.GetMtu()),
			DNS:                 req.GetDns(),
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
		Id: peer.ID.String(),
	}, nil
}

func (p *PeersImpl) Update(ctx context.Context, req *wgpb.UpdatePeerRequest) (*empty.Empty, error) {
	peer := req.GetPeer()
	fmask, err := fieldmask_utils.MaskFromPaths(req.FieldMask.Paths, mapNames)
	if err != nil {
		return nil, err
	}

	ID, err := uuid.Parse(peer.GetId())
	if err != nil {
		return nil, err
	}

	_, err = p.Service.Update(ctx,
		dto.UpdatePeerDTO{
			ID:                  ID,
			Name:                peer.GetName(),
			Email:               peer.GetEmail(),
			Description:         peer.GetDescription(),
			AddPresharedKey:     peer.GetAddPresharedKey(),
			RemovePresharedKey:  peer.GetRemovePresharedKey(),
			DNS:                 peer.GetDns(),
			MTU:                 int(peer.GetMtu()),
			PersistentKeepAlive: time.Duration(peer.GetPersistentKeepAlive()) * time.Second,
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

func (p *PeersImpl) Remove(ctx context.Context, req *wgpb.EntityIdRequest) (*empty.Empty, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	err = p.Service.Remove(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (p *PeersImpl) Get(ctx context.Context, req *wgpb.EntityIdRequest) (*wgpb.Peer, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	peer, err := p.Service.Get(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, err
	}

	return mapEntityPeerToPbPeer(peer), nil
}

func (p *PeersImpl) GetAll(ctx context.Context, req *wgpb.GetPeersRequest) (*wgpb.GetPeersResponse, error) {
	deviceIDStr := req.GetDeviceId()

	var deviceID uuid.UUID
	var err error

	if deviceIDStr != "" {
		deviceID, err = uuid.Parse(deviceIDStr)

		if err != nil {
			return nil, err
		}
	} else {
		deviceID = uuid.Nil
	}

	resp, err := p.Service.GetAll(ctx, dto.GetPeersRequestDTO{
		Skip:     int(req.GetSkip()),
		Limit:    int(req.GetLimit()),
		Search:   req.GetSearch(),
		DeviceID: deviceID,
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

	peerspb := make([]*wgpb.PeerAbridged, 0, len(resp.Peers))
	for _, peer := range resp.Peers {
		peerspb = append(peerspb, mapEntityPeerToPbPeerAbridged(peer))
	}

	return &wgpb.GetPeersResponse{
		Peers:   peerspb,
		Total:   int32(resp.Total),
		HasNext: resp.HasNext,
	}, nil
}

func (p *PeersImpl) Enable(ctx context.Context, req *wgpb.EntityIdRequest) (*empty.Empty, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	if err := p.Service.Enable(ctx, id); errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (p *PeersImpl) Disable(ctx context.Context, req *wgpb.EntityIdRequest) (*empty.Empty, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	if err := p.Service.Disable(ctx, id); errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (p *PeersImpl) DownloadConfig(ctx context.Context, req *wgpb.EntityIdRequest) (*wgpb.DownloadFileResponse, error) {
	ID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	file, err := p.Service.DownloadConfig(ctx, ID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, err
	}

	return &wgpb.DownloadFileResponse{
		Name: file.Name,
		Size: file.Size,
		Data: file.Data,
	}, nil
}

func (p *PeersImpl) DownloadQRCode(ctx context.Context, req *wgpb.EntityIdRequest) (*wgpb.DownloadFileResponse, error) {
	ID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	qr, err := p.Service.DownloadQRCode(ctx, ID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, err
	}

	return &wgpb.DownloadFileResponse{
		Name: qr.Name,
		Size: qr.Size,
		Data: qr.Data,
	}, nil
}
