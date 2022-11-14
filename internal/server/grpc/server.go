package grpcserver

import (
	"context"
	"fmt"
	peerpb "github.com/AZhur771/wg-grpc-api/api/proto"
	"google.golang.org/grpc"
	"net"
)

type ServerImpl struct {
	ctx    context.Context
	server *grpc.Server
	lis    net.Listener
	peerpb.UnimplementedPeerServiceServer
}

func NewServer(ctx context.Context, host string, port int, server *grpc.Server) (*ServerImpl, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}

	srvImpl := &ServerImpl{
		ctx:    ctx,
		lis:    lis,
		server: server,
	}

	peerpb.RegisterPeerServiceServer(server, srvImpl)

	return srvImpl, nil
}

func (s *ServerImpl) Start() error {
	return s.server.Serve(s.lis)
}

func (s *ServerImpl) Stop() {
	s.server.GracefulStop()
}
