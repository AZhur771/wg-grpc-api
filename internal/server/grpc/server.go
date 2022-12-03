package grpcserver

import (
	"context"
	"net"

	peerpb "github.com/AZhur771/wg-grpc-api/gen"
	"github.com/AZhur771/wg-grpc-api/internal/app"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ServerImpl struct {
	ctx     context.Context
	logger  app.Logger
	service app.PeerService

	server *grpc.Server
	lis    net.Listener
	peerpb.UnimplementedPeerServiceServer
}

func NewServer(ctx context.Context, logger app.Logger, service app.PeerService, addr string) (*ServerImpl, error) {
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
