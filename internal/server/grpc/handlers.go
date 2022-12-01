package grpcserver

import (
	"context"
	"errors"

	peerpb "github.com/AZhur771/wg-grpc-api/gen"
	"github.com/AZhur771/wg-grpc-api/internal/storage"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerImpl) AddPeer(
	ctx context.Context,
	addPeerRequest *peerpb.AddPeerRequest,
) (*peerpb.AddPeerResponse, error) {
	id, err := s.service.Add(ctx,
		addPeerRequest.GetAddPresharedKey(),
		addPeerRequest.GetPersistentKeepAlive().AsDuration(),
		addPeerRequest.GetDescription(),
	)
	if err != nil {
		return nil, err
	}

	return &peerpb.AddPeerResponse{
		Id: id.String(),
	}, nil
}

func (s *ServerImpl) RemovePeer(
	ctx context.Context,
	peerIDRequest *peerpb.PeerIdRequest,
) (*empty.Empty, error) {
	id, err := uuid.Parse(peerIDRequest.GetId())
	if err != nil {
		return nil, err
	}

	if err := s.service.Delete(ctx, id); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *ServerImpl) UpdatePeer(
	ctx context.Context,
	updatePeerRequest *peerpb.UpdatePeerRequest,
) (*empty.Empty, error) {
	updatePeerReq := updatePeerRequest.GetPeer()
	id, err := uuid.Parse(updatePeerReq.GetId())
	if err != nil {
		return nil, err
	}

	err = s.service.Update(ctx,
		id,
		updatePeerReq.GetAddPresharedKey(),
		updatePeerReq.GetPersistentKeepAlive().AsDuration(),
		updatePeerReq.GetDescription(),
		updatePeerRequest.GetUpdateMask().GetPaths(),
	)

	if errors.Is(err, storage.ErrPeerNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *ServerImpl) GetPeer(
	ctx context.Context,
	peerIDRequest *peerpb.PeerIdRequest,
) (*peerpb.Peer, error) {
	id, err := uuid.Parse(peerIDRequest.GetId())
	if err != nil {
		return nil, err
	}

	peer, err := s.service.Get(ctx, id)
	if errors.Is(err, storage.ErrPeerNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, err
	}

	return mapEntityPeerToPbPeer(peer), nil
}

func (s *ServerImpl) GetPeers(
	ctx context.Context,
	getPeersRequest *peerpb.GetPeersRequest,
) (*peerpb.GetPeersResponse, error) {
	paginatedPeers, err := s.service.GetAll(
		ctx,
		int(getPeersRequest.GetSkip()),
		int(getPeersRequest.GetLimit()),
	)
	if err != nil {
		return nil, err
	}

	peerspb := make([]*peerpb.Peer, 0, len(paginatedPeers.Peers))
	for _, peer := range paginatedPeers.Peers {
		peerspb = append(peerspb, mapEntityPeerToPbPeer(peer))
	}

	return &peerpb.GetPeersResponse{
		Peers:   peerspb,
		Total:   int32(paginatedPeers.Total),
		HasNext: paginatedPeers.HasNext,
	}, nil
}

func (s *ServerImpl) EnablePeer(
	ctx context.Context,
	peerIDRequest *peerpb.PeerIdRequest,
) (*empty.Empty, error) {
	// TODO implement me
	panic("implement me")
}

func (s *ServerImpl) DisablePeer(
	ctx context.Context,
	peerIDRequest *peerpb.PeerIdRequest,
) (*empty.Empty, error) {
	// TODO implement me
	panic("implement me")
}

func (s *ServerImpl) DownloadPeerConfig(
	ctx context.Context,
	peerIDRequest *peerpb.PeerIdRequest,
) (*peerpb.DownloadPeerConfigResponse, error) {
	id, err := uuid.Parse(peerIDRequest.GetId())
	if err != nil {
		return nil, err
	}

	peerConfig, err := s.service.DownloadConfig(ctx, id)
	if errors.Is(err, storage.ErrPeerNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, err
	}

	return &peerpb.DownloadPeerConfigResponse{
		Config: peerConfig,
	}, nil
}
