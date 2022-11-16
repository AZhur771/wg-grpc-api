package grpcserver

import (
	"context"

	peerpb "github.com/AZhur771/wg-grpc-api/api/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerImpl) AddPeer(context.Context, *peerpb.AddPeerRequest) (*peerpb.AddPeerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddPeer not implemented")
}

func (s *ServerImpl) RemovePeer(context.Context, *peerpb.RemovePeerRequest) (*peerpb.RemovePeerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemovePeer not implemented")
}

func (s *ServerImpl) UpdatePeer(context.Context, *peerpb.UpdatePeerRequest) (*peerpb.UpdatePeerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePeer not implemented")
}

func (s *ServerImpl) GetPeer(context.Context, *peerpb.GetPeerRequest) (*peerpb.Peer, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPeer not implemented")
}

func (s *ServerImpl) GetPeers(context.Context, *peerpb.GetPeersRequest) (*peerpb.GetPeersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPeers not implemented")
}
