package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	grpcserver "github.com/AZhur771/wg-grpc-api/internal/server/grpc"
	restserver "github.com/AZhur771/wg-grpc-api/internal/server/rest"
	deviceservice "github.com/AZhur771/wg-grpc-api/internal/service/device"
	peerservice "github.com/AZhur771/wg-grpc-api/internal/service/peer"
	"github.com/AZhur771/wg-grpc-api/internal/storage"
	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
	"golang.zx2c4.com/wireguard/wgctrl"
)

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

type config struct {
	IsProduction bool     `env:"PRODUCTION"`
	Host         string   `env:"HOST" envDefault:"localhost"`
	Port         int      `env:"PORT" envDefault:"3000"`
	Gateway      bool     `env:"GATEWAY"`
	GatewayPort  int      `env:"GATEWAY_PORT" envDefault:"3001"`
	ServeSwagger bool     `env:"SWAGGER"`
	Device       string   `env:"DEVICE" envDefault:"wg0"`
	Address      string   `env:"ADDRESS,required"`
	Endpoint     string   `env:"ENDPOINT,required"`
	Tokens       []string `env:"TOKENS" envSeparator:","`

	ServerCert string `env:"SERVER_CERT"`
	ServerKey  string `env:"SERVER_KEY"`

	// PeerFolder indicates folder where to store peer configs
	PeerFolder string `env:"PEER_FOLDER" envDefault:"${HOME}/wg_peers" envExpand:"true"`
}

func logErrorAndExit(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func initializeLogger(isProd bool) (*zap.Logger, error) {
	if isProd {
		return zap.NewProduction()
	}

	return zap.NewDevelopment()
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		if err := json.NewEncoder(os.Stdout).Encode(struct {
			Release   string
			BuildDate string
			GitHash   string
		}{
			Release:   release,
			BuildDate: buildDate,
			GitHash:   gitHash,
		}); err != nil {
			fmt.Printf("error while decoding version info: %v\n", err)
		}
		return
	}

	cfg := config{}
	err := env.Parse(&cfg,
		env.Options{
			Prefix: "WG_GRPC_API_",
		},
	)
	logErrorAndExit(err)

	logger, err := initializeLogger(cfg.IsProduction)
	logErrorAndExit(err)
	defer logger.Sync()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	grpcAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)        // GRPC API address
	restAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.GatewayPort) // REST API address
	logErrorAndExit(err)

	// Wireguard controller
	wgclient, err := wgctrl.New()
	logErrorAndExit(err)

	// Storage
	storage := storage.NewStorage(logger, cfg.PeerFolder)

	// Device service
	deviceService := deviceservice.NewDeviceService(logger, wgclient, storage)
	logErrorAndExit(deviceService.Setup(ctx, cfg.Device, cfg.Endpoint, cfg.Address))

	// Peer service
	peerService := peerservice.NewPeerService(logger, deviceService, storage)

	// GRPC server
	grpcSrv, err := grpcserver.New(ctx, logger, peerService, deviceService, grpcAddr, cfg.Tokens, cfg.ServerCert, cfg.ServerKey)
	logErrorAndExit(err)

	var restSrv *restserver.Server
	var wg sync.WaitGroup

	if cfg.Gateway {
		wg.Add(2)
		// REST gateway
		restSrv, err = restserver.New(ctx, logger, restAddr, grpcAddr, cfg.ServeSwagger, cfg.ServerCert, cfg.ServerKey)
		logErrorAndExit(err)
	} else {
		wg.Add(1)
	}

	go func() {
		defer wg.Done()
		log.Println("grpc server is up and running")
		if err := grpcSrv.Start(); err != nil {
			logger.Error("failed to start grpc server", zap.Error(err))
			cancel()
		}
	}()

	if cfg.Gateway {
		go func() {
			defer wg.Done()
			log.Println("rest gateway server is up and running")
			if cfg.ServeSwagger {
				log.Printf("swagger docs available at http://%s/swagger-ui\n", restAddr)
			}

			if err := restSrv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Error("failed to start rest gateway server", zap.Error(err))
				cancel()
			}
		}()
	}

	// wait for signal to stop servers
	<-ctx.Done()

	if cfg.Gateway {
		// TODO: add timeout
		if err := restSrv.Stop(context.TODO()); err != nil {
			logger.Error("failed to stop rest gateway server", zap.Error(err))
		}
	}

	grpcSrv.Stop()

	wg.Wait()

	os.Exit(1) //nolint:gocritic
}
