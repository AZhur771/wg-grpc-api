package main

import (
	"context"
	"flag"
	"fmt"
	grpcclient "github.com/AZhur771/wg-grpc-api/internal/client/grpc"
	grpcserver "github.com/AZhur771/wg-grpc-api/internal/server/grpc"
	httpserver "github.com/AZhur771/wg-grpc-api/internal/server/http"
	"google.golang.org/grpc"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AZhur771/wg-grpc-api/internal/config"
)

var configFile string

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

func init() {
	flag.StringVar(&configFile, "config", "/etc/wgapi/config.yaml", "Path to configuration file")
}

func logErrorAndExit(msg string, err error) {
	if err != nil {
		log.Fatalf(msg, err)
	}
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		fmt.Printf("Release: %s\nBuild date: %s\nGit hash: %s", release, buildDate, gitHash)
		return
	}

	cfg, err := config.ParseConfig(configFile)

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)
	defer cancel()

	grpcSrv, err := grpcserver.NewServer(ctx, cfg.Server.Host, cfg.Server.Port, grpc.NewServer())
	logErrorAndExit("failed to create grpc server: %s", err)

	var httpSrv *httpserver.ServerImpl

	go func() {
		if err := grpcSrv.Start(); err != nil {
			// log error
			cancel()
		}
	}()

	if cfg.Server.EnableGateway {
		clientConn, err := grpcclient.NewClientConn(ctx, cfg.Server.Host, cfg.Server.Port)
		logErrorAndExit("failed to create grpc client conn: %s", err)

		httpSrv, err = httpserver.NewServer(
			ctx,
			cfg.Server.Host,
			cfg.Server.GatewayPort,
			clientConn,
			cfg.Server.EnableSwagger,
		)
		logErrorAndExit("failed to create http gateway server: %s", err)

		go func() {
			if err := httpSrv.Start(); err != nil {
				// log error
				cancel()
			}
		}()
	}

	<-ctx.Done()

	if cfg.Server.EnableGateway {
		ctx, cancel = context.WithTimeout(
			context.Background(),
			time.Duration(cfg.Server.ShutdownTimeout)*time.Millisecond,
		)
		defer cancel()

		if err := httpSrv.Stop(ctx); err != nil {
			// log error
		}
	}

	// TODO: check that stop methods can be invoked without errors + choose logger

	grpcSrv.Stop()
	os.Exit(1)
}
