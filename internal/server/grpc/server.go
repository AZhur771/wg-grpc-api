package grpcserver

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	wgpb "github.com/AZhur771/wg-grpc-api/gen"
	"github.com/AZhur771/wg-grpc-api/internal/app"
	"github.com/AZhur771/wg-grpc-api/internal/certs"
	"github.com/AZhur771/wg-grpc-api/internal/server/grpc/handlers"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Server struct {
	server *grpc.Server
	lis    net.Listener
}

func handlePanic(p interface{}) (err error) {
	return status.Errorf(codes.Unknown, "panic triggered: %v", p)
}

func New(ctx context.Context, logger *zap.Logger, peerService app.PeerService, deviceService app.DeviceService,
	addr, host, caCertPath, certPath, keyPath string, certOpt tls.ClientAuthType,
	tokens []string,
) (*Server, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("grpc server: %w", err)
	}

	recOptions := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(handlePanic),
	}

	logOptions := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	}

	serverOptions := []grpc.ServerOption{}

	if certPath != "" || keyPath != "" {
		tlsCredentials, err := certs.LoadTLSCredentials(host, caCertPath, certPath, keyPath, certOpt)
		if err != nil {
			return nil, fmt.Errorf("grpc server: %w", err)
		}

		serverOptions = append(serverOptions, grpc.Creds(tlsCredentials))
	}

	serverOptions = append(
		serverOptions,
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(interceptorLogger(logger), logOptions...),
			withUnaryServerInterceptor(tokens),
			grpc_recovery.UnaryServerInterceptor(recOptions...),
		),
	)

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
