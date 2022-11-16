package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	grpcclient "github.com/AZhur771/wg-grpc-api/internal/client/grpc"
	grpcserver "github.com/AZhur771/wg-grpc-api/internal/server/grpc"
	restserver "github.com/AZhur771/wg-grpc-api/internal/server/http"
	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
)

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

type config struct {
	Host          string        `env:"HOST" envDefault:"localhost"`
	Port          int           `env:"PORT" envDefault:"3000"`
	GatewayPort   int           `env:"GATEWAY_PORT" envDefault:"3001"`
	EnableSwagger bool          `env:"SWAGGER"`
	IsProduction  bool          `env:"PRODUCTION"`
	Timeout       time.Duration `env:"TIMEOUT" envDefault:"3000s"`
	PeerFolder    string        `env:"PEER_FOLDER" envDefault:"${HOME}/wg_peers" envExpand:"true"`
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
		fmt.Printf("Release: %s\nBuild date: %s\nGit hash: %s", release, buildDate, gitHash)
		return
	}

	cfg := config{}
	err := env.Parse(&cfg, env.Options{
		Prefix: "WG_GRPC_API_",
	})
	logErrorAndExit(err)

	logger, err := initializeLogger(cfg.IsProduction)
	logErrorAndExit(err)
	defer logger.Sync() // flush logger

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)
	defer cancel()

	grpcSrv, err := grpcserver.NewServer(
		ctx,
		logger,
		cfg.Host,
		cfg.Port,
	)
	logErrorAndExit(err)

	clientConn, err := grpcclient.NewClientConn(
		ctx,
		cfg.Host,
		cfg.Port,
	)
	logErrorAndExit(err)

	httpSrv, err := restserver.NewServer(
		ctx,
		logger,
		clientConn,
		cfg.Host,
		cfg.GatewayPort,
		cfg.EnableSwagger,
	)
	logErrorAndExit(err)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		log.Println("grpc server is up and running")
		if err := grpcSrv.Start(); err != nil {
			logger.Error("failed to start grpc server", zap.Error(err))
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		log.Println("http gateway server is up and running")
		if cfg.EnableSwagger {
			log.Printf("swagger docs available at http://%s:%d/swagger-ui\n", cfg.Host, cfg.GatewayPort)
		}

		if err := httpSrv.Start(); err != nil {
			logger.Error("failed to start http gateway server", zap.Error(err))
			cancel()
		}
	}()

	// wait for signal to stop servers
	<-ctx.Done()

	// gracefully stop http gateway server
	if err := httpSrv.Stop(context.TODO()); err != nil {
		logger.Error("failed to stop http gateway server", zap.Error(err))
	}

	// gracefully stop grpc server
	grpcSrv.Stop()

	// wait for all goroutines to stop
	wg.Wait()

	// exit
	os.Exit(1)
}
