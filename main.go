package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AZhur771/wg-grpc-api/internal/app"
	database "github.com/AZhur771/wg-grpc-api/internal/db"
	devicerepo "github.com/AZhur771/wg-grpc-api/internal/repo/device"
	peerrepo "github.com/AZhur771/wg-grpc-api/internal/repo/peer"
	"github.com/AZhur771/wg-grpc-api/internal/server"
	deviceservice "github.com/AZhur771/wg-grpc-api/internal/service/device"
	peerservice "github.com/AZhur771/wg-grpc-api/internal/service/peer"
	_ "github.com/AZhur771/wg-grpc-api/migrations"
	"github.com/caarlos0/env/v6"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"golang.zx2c4.com/wireguard/wgctrl"
)

const envPrefix = "WG_GRPC_API_"

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

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

	cfg := app.Config{}
	err := env.Parse(&cfg,
		env.Options{
			Prefix: envPrefix,
		},
	)
	logErrorAndExit(err)

	logger, err := initializeLogger(cfg.IsProduction)
	logErrorAndExit(err)
	defer logger.Sync()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	mdb, err := sql.Open("pgx", database.GetConnectionString(cfg))
	logErrorAndExit(err)
	err = goose.Up(mdb, "migrations")
	logErrorAndExit(err)

	db, err := database.New(cfg)
	logErrorAndExit(err)
	deviceRepo := devicerepo.New(db)
	peerRepo := peerrepo.New(db)

	wgclient, err := wgctrl.New()
	logErrorAndExit(err)

	deviceService := deviceservice.NewDeviceService(logger, wgclient, deviceRepo, peerRepo)
	err = deviceService.SyncDevices(ctx)
	logErrorAndExit(err)

	peerService := peerservice.NewPeerService(logger, deviceService, deviceRepo, peerRepo)

	server, err := server.NewServer(ctx, logger, peerService, deviceService, cfg)
	logErrorAndExit(err)

	server.Run(ctx, cancel)

	os.Exit(1) //nolint:gocritic
}
