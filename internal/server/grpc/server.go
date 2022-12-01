package grpcserver

import (
	"context"
	"net"

	peerpb "github.com/AZhur771/wg-grpc-api/gen"
	"github.com/AZhur771/wg-grpc-api/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ServerImpl struct {
	ctx     context.Context
	logger  *zap.Logger
	server  *grpc.Server
	service service.PeerService
	lis     net.Listener
	peerpb.UnimplementedPeerServiceServer
}

func NewServer(ctx context.Context, logger *zap.Logger, service service.PeerService, addr string) (*ServerImpl, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(getUnaryServerInterceptor(logger)),
	)
	reflection.Register(srv)

	srvImpl := &ServerImpl{
		logger:  logger,
		service: service,
		ctx:     ctx,
		lis:     lis,
		server:  srv,
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
