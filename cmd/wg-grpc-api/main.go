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
	"github.com/AZhur771/wg-grpc-api/internal/service"
	"github.com/AZhur771/wg-grpc-api/internal/storage"
	"github.com/AZhur771/wg-grpc-api/pkg/redis"
	"github.com/AZhur771/wg-grpc-api/pkg/wg"
	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
)

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

type config struct {
	IsProduction bool   `env:"PRODUCTION"`
	Host         string `env:"HOST" envDefault:"localhost"`
	Port         int    `env:"PORT" envDefault:"3000"`
	Gateway      bool   `env:"GATEWAY"`
	GatewayPort  int    `env:"GATEWAY_PORT" envDefault:"3001"`
	ServeSwagger bool   `env:"SWAGGER"`
	Device       string `env:"DEVICE" envDefault:"wg0"`
	Address      string `env:"ADDRESS,required"`
	Endpoint     string `env:"ENDPOINT,required"`
	PeerFolder   string `env:"PEER_FOLDER" envDefault:"${HOME}/wg_peers" envExpand:"true"`

	RedisHost     string `env:"REDIS_HOST" envDefault:"localhost"`
	RedisPort     int    `env:"REDIS_PORT" envDefault:"6379"`
	RedisPassword string `env:"REDIS_PASSWORD"`
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
	defer logger.Sync() // flush logger

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	rsAddr := fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort) // Redis address
	grpcAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)         // GRPC API address
	restAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.GatewayPort)  // REST API address

	rsclient, err := redis.New(ctx, rsAddr, cfg.RedisPassword)
	logErrorAndExit(err)

	wgclient, err := wg.New(cfg.Device)
	logErrorAndExit(err)

	peerStorage := storage.NewPeerStorage(logger, rsclient)
	peerService := service.NewPeerService(logger, wgclient, peerStorage)
	logErrorAndExit(peerService.Setup(ctx, cfg.Device, cfg.Address, cfg.Endpoint, cfg.PeerFolder))

	grpcSrv, err := grpcserver.NewServer(ctx, logger, peerService, grpcAddr)
	logErrorAndExit(err)

	var restSrv *restserver.ServerImpl
	var wg sync.WaitGroup

	if cfg.Gateway {
		wg.Add(2)
		restSrv, err = restserver.NewServer(ctx, logger, restAddr, grpcAddr, cfg.ServeSwagger)
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
				// TODO: if in docker print correct url
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
		// gracefully stop rest gateway server
		// TODO: add timeout
		if err := restSrv.Stop(context.TODO()); err != nil {
			logger.Error("failed to stop rest gateway server", zap.Error(err))
		}
	}

	// gracefully stop grpc server
	grpcSrv.Stop()

	// wait for all goroutines to stop
	wg.Wait()

	// exit
	os.Exit(1) //nolint:gocritic
}
