package grpcserver

import (
	"context"
	"fmt"
	"net"

	wgpb "github.com/AZhur771/wg-grpc-api/gen"
	"github.com/AZhur771/wg-grpc-api/internal/app"
	"github.com/AZhur771/wg-grpc-api/internal/certs"
	"github.com/AZhur771/wg-grpc-api/internal/server/grpc/handlers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	server *grpc.Server
	lis    net.Listener
}

func New(ctx context.Context, logger app.Logger, peerService app.PeerService, deviceService app.DeviceService, addr string, tokens []string, certPath, keyPath string) (*Server, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("grpc server: %w", err)
	}

	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(withUnaryServerInterceptor(logger, tokens)),
	}

	if certPath != "" || keyPath != "" {
		tlsCredentials, err := certs.LoadTLSCredentials(certPath, keyPath)
		if err != nil {
			return nil, fmt.Errorf("grpc server: %w", err)
		}

		serverOptions = append(serverOptions, grpc.Creds(tlsCredentials))
	}

	grpcsrv := grpc.NewServer(serverOptions...)

	srv := &Server{
		lis:    lis,
		server: grpcsrv,
	}

	peers := handlers.NewPeersImpl(ctx, logger, peerService)
	device := handlers.NewDeviceImpl(ctx, logger, deviceService)

	wgpb.RegisterDeviceServiceServer(grpcsrv, device)
	wgpb.RegisterPeerServiceServer(grpcsrv, peers)
	reflection.Register(grpcsrv)

	return srv, nil
}

func (s *Server) Start() error {
	return s.server.Serve(s.lis)
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}
