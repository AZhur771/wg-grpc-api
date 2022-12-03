package grpcserver

import (
	"context"
	"errors"
	"fmt"

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
	updateMask := updatePeerRequest.GetUpdateMask()
	updateMask.Normalize()

	updatePeerData := updatePeerRequest.GetPeer()

	if !updateMask.IsValid(updatePeerData) {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("wrong update mask values: %v",
			updatePeerRequest.GetUpdateMask().GetPaths()))
	}

	id, err := uuid.Parse(updatePeerData.GetId())
	if err != nil {
		return nil, err
	}

	err = s.service.Update(ctx,
		id,
		updatePeerData.GetAddPresharedKey(),
		updatePeerData.GetPersistentKeepAlive().AsDuration(),
		updatePeerData.GetDescription(),
		updateMask.GetPaths(),
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
	getPeersRequest *empty.Empty,
) (*peerpb.GetPeersResponse, error) {
	peers, err := s.service.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	peerspb := make([]*peerpb.Peer, 0, len(peers))
	for _, peer := range peers {
		peerspb = append(peerspb, mapEntityPeerToPbPeer(peer))
	}

	return &peerpb.GetPeersResponse{
		Peers: peerspb,
	}, nil
}

func (s *ServerImpl) StreamPeers(
	streamPeersRequest *empty.Empty,
	stream peerpb.PeerService_StreamPeersServer,
) error {
	// TODO: should specify context
	peers, err := s.service.GetAll(context.TODO())
	if err != nil {
		return err
	}

	for _, peer := range peers {
		if err := stream.Send(mapEntityPeerToPbPeer(peer)); err != nil {
			return err
		}
	}

	return nil
}

func (s *ServerImpl) EnablePeer(
	ctx context.Context,
	peerIDRequest *peerpb.PeerIdRequest,
) (*empty.Empty, error) {
	id, err := uuid.Parse(peerIDRequest.GetId())
	if err != nil {
		return nil, err
	}

	err = s.service.Enable(ctx, id)
	if errors.Is(err, storage.ErrPeerNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *ServerImpl) DisablePeer(
	ctx context.Context,
	peerIDRequest *peerpb.PeerIdRequest,
) (*empty.Empty, error) {
	id, err := uuid.Parse(peerIDRequest.GetId())
	if err != nil {
		return nil, err
	}

	err = s.service.Disable(ctx, id)
	if errors.Is(err, storage.ErrPeerNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
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
