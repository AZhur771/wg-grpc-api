package handlers

import (
	"context"
	"database/sql"
	"errors"
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
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (p *PeersImpl) AddPeer(ctx context.Context, req *wgpb.AddPeerRequest) (*wgpb.EntityIdRequest, error) {
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

	return &wgpb.EntityIdRequest{
		Id: peer.ID.String(),
	}, nil
}

func (p *PeersImpl) UpdatePeer(ctx context.Context, req *wgpb.UpdatePeerRequest) (*empty.Empty, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	_, err = p.Service.Update(ctx,
		dto.UpdatePeerDTO{
			ID:                  id,
			Name:                req.GetName(),
			Email:               req.GetEmail(),
			Description:         req.GetDescription(),
			AddPresharedKey:     req.GetAddPresharedKey(),
			RemovePresharedKey:  req.GetRemovePresharedKey(),
			PersistentKeepAlive: time.Duration(req.GetPersistentKeepAlive()) * time.Second,
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

func (p *PeersImpl) RemovePeer(ctx context.Context, req *wgpb.EntityIdRequest) (*empty.Empty, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	if err := p.Service.Remove(ctx, id); errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (p *PeersImpl) GetPeer(ctx context.Context, req *wgpb.EntityIdRequest) (*wgpb.Peer, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	peer, err := p.Service.Get(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, err
	}

	return mapEntityPeerToPbPeer(peer), nil
}

func (p *PeersImpl) GetPeers(ctx context.Context, req *wgpb.GetPeersRequest) (*wgpb.GetPeersResponse, error) {
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
	if errors.Is(err, peerservice.ErrInvalidPaginationParams) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	} else if err != nil {
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

func (p *PeersImpl) EnablePeer(ctx context.Context, req *wgpb.EntityIdRequest) (*empty.Empty, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	if err := p.Service.Enable(ctx, id); errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (p *PeersImpl) DisablePeer(ctx context.Context, req *wgpb.EntityIdRequest) (*empty.Empty, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	if err := p.Service.Disable(ctx, id); errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (p *PeersImpl) DownloadPeerConfig(ctx context.Context, req *wgpb.EntityIdRequest) (*wgpb.DownloadFileResponse, error) {
	ID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	file, err := p.Service.DownloadConfig(ctx, ID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, err
	}

	return &wgpb.DownloadFileResponse{
		Name: file.Name,
		Size: file.Size,
		Data: file.Data,
	}, nil
}

func (p *PeersImpl) DownloadPeerQRCode(ctx context.Context, req *wgpb.EntityIdRequest) (*wgpb.DownloadFileResponse, error) {
	ID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	qr, err := p.Service.DownloadQRCode(ctx, ID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, err
	}

	return &wgpb.DownloadFileResponse{
		Name: qr.Name,
		Size: qr.Size,
		Data: qr.Data,
	}, nil
}

func mapEntityPeerToPbPeer(peer *entity.Peer) *wgpb.Peer {
	return &wgpb.Peer{
		Id:                  peer.ID.String(),
		DeviceId:            peer.DeviceID.String(),
		Name:                peer.Name,
		Email:               peer.Email,
		PublicKey:           peer.PublicKey.String(),
		Endpoint:            peer.Endpoint.String(),
		PersistentKeepAlive: int32(peer.PersistentKeepaliveInterval.Seconds()),
		AllowedIps:          peer.AllowedIPs,
		ProtocolVersion:     uint32(peer.ProtocolVersion),
		ReceiveBytes:        peer.ReceiveBytes,
		TransmitBytes:       peer.TransmitBytes,
		LastHandshake:       timestamppb.New(peer.LastHandshakeTime),
		HasPresharedKey:     peer.HasPresharedKey,
		IsEnabled:           peer.IsEnabled,
		IsActive:            peer.IsActive,
		Description:         peer.Description,
		Dns:                 peer.DNS,
		Mtu:                 int32(peer.MTU),
	}
}

func mapEntityPeerToPbPeerAbridged(peer *entity.Peer) *wgpb.PeerAbridged {
	return &wgpb.PeerAbridged{
		Id:                  peer.ID.String(),
		DeviceId:            peer.DeviceID.String(),
		Name:                peer.Name,
		Email:               peer.Email,
		PublicKey:           peer.PublicKey.String(),
		PersistentKeepAlive: int32(peer.PersistentKeepaliveInterval.Seconds()),
		AllowedIps:          peer.AllowedIPs,
		HasPresharedKey:     peer.HasPresharedKey,
		IsEnabled:           peer.IsEnabled,
		Description:         peer.Description,
		Dns:                 peer.DNS,
		Mtu:                 int32(peer.MTU),
	}
}
