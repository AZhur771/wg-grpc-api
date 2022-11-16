package grpcserver

import (
	"context"
	"fmt"
	"net"

	peerpb "github.com/AZhur771/wg-grpc-api/api/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ServerImpl struct {
	ctx    context.Context
	logger *zap.Logger
	server *grpc.Server
	lis    net.Listener
	peerpb.UnimplementedPeerServiceServer
}

func NewServer(ctx context.Context, logger *zap.Logger, host string, port int) (*ServerImpl, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(getUnaryServerInterceptor(logger)),
	)

	srvImpl := &ServerImpl{
		logger: logger,
		ctx:    ctx,
		lis:    lis,
		server: srv,
	}

	peerpb.RegisterPeerServiceServer(srv, srvImpl)

	return srvImpl, nil
}

func (s *ServerImpl) Start() error {
	return s.server.Serve(s.lis)
}

func (s *ServerImpl) Stop() {
	s.server.GracefulStop()
}
