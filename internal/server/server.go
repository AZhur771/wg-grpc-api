package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/AZhur771/wg-grpc-api/internal/app"
	grpcserver "github.com/AZhur771/wg-grpc-api/internal/server/grpc"
	restserver "github.com/AZhur771/wg-grpc-api/internal/server/rest"
	deviceservice "github.com/AZhur771/wg-grpc-api/internal/service/device"
	peerservice "github.com/AZhur771/wg-grpc-api/internal/service/peer"
	"go.uber.org/zap"
)

type Server struct {
	logger   *zap.Logger
	grpcSrv  *grpcserver.Server
	restSrv  *restserver.Server
	gateway  bool
	swagger  bool
	grpcAddr string
	restAddr string
}

func NewServer(ctx context.Context, logger *zap.Logger,
	peerService *peerservice.PeerService, deviceService *deviceservice.DeviceService, cfg app.Config,
) (*Server, error) {
	grpcAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)        // GRPC API address
	restAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.GatewayPort) // REST API address

	// GRPC server
	grpcSrv, err := grpcserver.New(ctx, logger, peerService, deviceService, grpcAddr,
		cfg.Host, cfg.CaCert, cfg.Cert, cfg.Key, tls.ClientAuthType(cfg.CertOpt), cfg.Tokens)
	if err != nil {
		return nil, err
	}

	var restSrv *restserver.Server

	if cfg.Gateway {
		// REST gateway
		restSrv, err = restserver.New(ctx, logger, cfg.ServeSwagger, restAddr, grpcAddr,
			cfg.Host, cfg.CaCert, cfg.Cert, cfg.Key, tls.ClientAuthType(cfg.CertOpt))
		if err != nil {
			return nil, err
		}
	}

	return &Server{
		gateway:  cfg.Gateway,
		swagger:  cfg.ServeSwagger,
		logger:   logger,
		grpcSrv:  grpcSrv,
		restSrv:  restSrv,
		restAddr: restAddr,
		grpcAddr: grpcAddr,
	}, nil
}

func (s *Server) Run(ctx context.Context, cancel context.CancelFunc) {
	var wg sync.WaitGroup

	if s.gateway {
		wg.Add(2)
	} else {
		wg.Add(1)
	}

	go func() {
		defer wg.Done()
		log.Printf("grpc server is up and running at %s", s.grpcAddr)
		if err := s.grpcSrv.Start(); err != nil {
			s.logger.Error("failed to start grpc server", zap.Error(err))
			cancel()
		}
	}()

	if s.gateway {
		go func() {
			defer wg.Done()
			log.Printf("rest gateway server is up and running at %s", s.restAddr)
			if s.swagger {
				log.Printf("swagger docs available at http://%s/swagger-ui\n", s.restAddr)
			}

			if err := s.restSrv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				s.logger.Error("failed to start rest gateway server", zap.Error(err))
				cancel()
			}
		}()
	}

	// wait for signal to stop servers
	<-ctx.Done()

	if s.gateway {
		if err := s.restSrv.Stop(context.Background()); err != nil {
			s.logger.Error("failed to stop rest gateway server", zap.Error(err))
		}
	}

	s.grpcSrv.Stop()

	wg.Wait()
}
