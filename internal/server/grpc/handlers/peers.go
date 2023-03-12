package handlers

import (
	"context"
	"errors"

	wgpb "github.com/AZhur771/wg-grpc-api/gen"
	"github.com/AZhur771/wg-grpc-api/internal/app"
	"github.com/AZhur771/wg-grpc-api/internal/dto"
	"github.com/AZhur771/wg-grpc-api/internal/entity"
	peerservice "github.com/AZhur771/wg-grpc-api/internal/service/peer"
	peerstorage "github.com/AZhur771/wg-grpc-api/internal/storage/peer"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PeersImpl struct {
	Ctx     context.Context
	Logger  app.Logger
	Service app.PeerService

	wgpb.UnimplementedPeerServiceServer
}

func NewPeersImpl(ctx context.Context, logger app.Logger, service app.PeerService) *PeersImpl {
	return &PeersImpl{
		Ctx:     ctx,
		Logger:  logger,
		Service: service,
	}
}

func (p *PeersImpl) AddPeer(ctx context.Context, addPeerRequest *wgpb.AddPeerRequest) (*wgpb.PeerIdRequest, error) {
	peer, err := p.Service.Add(ctx,
		dto.AddPeerDTO{
			Name:                addPeerRequest.GetName(),
			Email:               addPeerRequest.GetEmail(),
			Description:         addPeerRequest.GetDescription(),
			Tags:                addPeerRequest.GetTags(),
			AddPresharedKey:     addPeerRequest.GetAddPresharedKey(),
			PersistentKeepAlive: addPeerRequest.GetPersistentKeepAlive().AsDuration(),
		},
	)
	if err != nil {
		return nil, err
	}

	return &wgpb.PeerIdRequest{
		Id: peer.ID.String(),
	}, nil
}

func (p *PeersImpl) RemovePeer(ctx context.Context, peerIDRequest *wgpb.PeerIdRequest) (*empty.Empty, error) {
	id, err := uuid.Parse(peerIDRequest.GetId())
	if err != nil {
		return nil, err
	}

	if err := p.Service.Remove(ctx, id); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (p *PeersImpl) UpdatePeer(ctx context.Context, updatePeerRequest *wgpb.UpdatePeerRequest) (*empty.Empty, error) {
	ID, err := uuid.Parse(updatePeerRequest.GetId())
	if err != nil {
		return nil, err
	}

	_, err = p.Service.Update(ctx,
		dto.UpdatePeerDTO{
			ID:                  ID,
			Name:                updatePeerRequest.GetName(),
			Email:               updatePeerRequest.GetEmail(),
			Description:         updatePeerRequest.GetDescription(),
			Tags:                updatePeerRequest.GetTags(),
			AddPresharedKey:     updatePeerRequest.GetAddPresharedKey(),
			RemovePresharedKey:  updatePeerRequest.GetRemovePresharedKey(),
			PersistentKeepAlive: updatePeerRequest.PersistentKeepAlive.AsDuration(),
		},
	)

	if errors.Is(err, peerstorage.ErrPeerNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (p *PeersImpl) GetPeer(ctx context.Context, peerIDRequest *wgpb.PeerIdRequest) (*wgpb.Peer, error) {
	id, err := uuid.Parse(peerIDRequest.GetId())
	if err != nil {
		return nil, err
	}

	peer, err := p.Service.Get(ctx, id)
	if errors.Is(err, peerstorage.ErrPeerNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, err
	}

	return mapEntityPeerToPbPeer(peer), nil
}

func (p *PeersImpl) GetPeers(ctx context.Context, getPeersRequest *wgpb.GetPeersRequest) (*wgpb.GetPeersResponse, error) {
	getPeersResponse, err := p.Service.GetAll(ctx, dto.GetPeersRequestDTO{
		Skip:  int(getPeersRequest.GetSkip()),
		Limit: int(getPeersRequest.GetLimit()),
	})

	if errors.Is(err, peerservice.ErrInvalidPaginationParams) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err != nil {
		return nil, err
	}

	peerspb := make([]*wgpb.Peer, 0, len(getPeersResponse.Peers))
	for _, peer := range getPeersResponse.Peers {
		peerspb = append(peerspb, mapEntityPeerToPbPeer(peer))
	}

	return &wgpb.GetPeersResponse{
		Peers:   peerspb,
		Total:   int32(getPeersResponse.Total),
		HasNext: getPeersResponse.HasNext,
	}, nil
}

func (p *PeersImpl) EnablePeer(ctx context.Context, peerIDRequest *wgpb.PeerIdRequest) (*empty.Empty, error) {
	ID, err := uuid.Parse(peerIDRequest.GetId())
	if err != nil {
		return nil, err
	}

	if err := p.Service.Enable(ctx, ID); err != nil {
		if errors.Is(err, peerstorage.ErrPeerNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, err
	}

	return &empty.Empty{}, nil
}

func (p *PeersImpl) DisablePeer(ctx context.Context, peerIDRequest *wgpb.PeerIdRequest) (*empty.Empty, error) {
	ID, err := uuid.Parse(peerIDRequest.GetId())
	if err != nil {
		return nil, err
	}

	if err := p.Service.Disable(ctx, ID); err != nil {
		if errors.Is(err, peerstorage.ErrPeerNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, err
	}

	return &empty.Empty{}, nil
}

func (p *PeersImpl) DownloadPeerConfig(ctx context.Context, peerIDRequest *wgpb.PeerIdRequest) (*wgpb.DownloadFileResponse, error) {
	ID, err := uuid.Parse(peerIDRequest.GetId())
	if err != nil {
		return nil, err
	}

	downloadFileDTO, err := p.Service.DownloadConfig(ctx, ID)
	if errors.Is(err, peerstorage.ErrPeerNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, err
	}

	return &wgpb.DownloadFileResponse{
		Name: downloadFileDTO.Name,
		Size: downloadFileDTO.Size,
		Data: downloadFileDTO.Data,
	}, nil
}

func (p *PeersImpl) DownloadPeerQRCode(ctx context.Context, peerIDRequest *wgpb.PeerIdRequest) (*wgpb.DownloadFileResponse, error) {
	ID, err := uuid.Parse(peerIDRequest.GetId())
	if err != nil {
		return nil, err
	}

	downloadFileDTO, err := p.Service.DownloadQRCode(ctx, ID)
	if errors.Is(err, peerstorage.ErrPeerNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, err
	}

	return &wgpb.DownloadFileResponse{
		Name: downloadFileDTO.Name,
		Size: downloadFileDTO.Size,
		Data: downloadFileDTO.Data,
	}, nil
}

func mapEntityPeerToPbPeer(peer *entity.Peer) *wgpb.Peer {
	return &wgpb.Peer{
		Id:                  peer.ID.String(),
		Name:                peer.Name,
		Email:               peer.Email,
		PublicKey:           peer.PublicKey.String(),
		Endpoint:            peer.Endpoint.String(),
		PersistentKeepAlive: durationpb.New(peer.PersistentKeepaliveInterval),
		AllowedIps:          peer.AllowedIPs,
		ProtocolVersion:     uint32(peer.ProtocolVersion),
		ReceiveBytes:        peer.ReceiveBytes,
		TransmitBytes:       peer.TransmitBytes,
		LastHandshake:       timestamppb.New(peer.LastHandshakeTime),
		HasPresharedKey:     peer.HasPresharedKey,
		IsEnabled:           peer.IsEnabled,
		IsActive:            peer.IsActive,
		Tags:                peer.Tags,
		Description:         peer.Description,
	}
}
