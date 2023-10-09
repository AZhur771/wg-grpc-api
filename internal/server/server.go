package server

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	wgpb "github.com/AZhur771/wg-grpc-api/gen"
	"github.com/AZhur771/wg-grpc-api/internal/app"
	"github.com/AZhur771/wg-grpc-api/internal/certs"
	"github.com/AZhur771/wg-grpc-api/internal/server/handlers"
	deviceservice "github.com/AZhur771/wg-grpc-api/internal/service/device"
	peerservice "github.com/AZhur771/wg-grpc-api/internal/service/peer"
	"github.com/AZhur771/wg-grpc-api/third_party"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	defaultGatewayPrefix string = "/api/"
	defaultSwaggerPrefix string = "/swagger-ui/"
)

type Server struct {
	logger  *zap.Logger
	grpcSrv *grpc.Server
	mux     *http.ServeMux
	addr    string
	cert    string
	key     string
	swagger bool
	tls     bool
}

func handlePanic(p interface{}) (err error) {
	return status.Errorf(codes.Unknown, "panic triggered: %v", p)
}

func serveSwagger(mux *http.ServeMux) error {
	mime.AddExtensionType(".svg", "image/svg+xml")

	subFS, err := fs.Sub(third_party.OpenAPIV2, "openapiv2")
	if err != nil {
		return err
	}

	// Expose files in third_party/openapiv2/ on <host>/swagger-ui/
	fileServer := http.FileServer(http.FS(subFS))

	mux.Handle(defaultSwaggerPrefix, http.StripPrefix(defaultSwaggerPrefix, fileServer))

	// Allow downloading swagger.json file at <host>/swagger-ui/swagger.json
	mux.HandleFunc(
		fmt.Sprintf("%s/swagger.json", defaultSwaggerPrefix),
		func(w http.ResponseWriter, req *http.Request) {
			f, err := subFS.Open("peer_service.swagger.json")
			if err != nil {
				w.WriteHeader(505)
				return
			}
			defer f.Close()

			io.Copy(w, f)
		},
	)

	return nil
}

func grpcHandlerFunc(srv *http2.Server, grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), srv)
}

func NewServer(ctx context.Context, logger *zap.Logger,
	peerService *peerservice.PeerService, deviceService *deviceservice.DeviceService, cfg app.Config,
) (*Server, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port) // gRPC/REST API address

	logOptions := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	}

	recOptions := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(handlePanic),
	}

	grpcOptions := make([]grpc.ServerOption, 0)

	if cfg.Cert != "" || cfg.Key != "" {
		tlsCredentials, err := certs.LoadTLSCredentials(cfg.Cert, cfg.Key)
		if err != nil {
			return nil, fmt.Errorf("grpc server: %w", err)
		}

		grpcOptions = append(grpcOptions, grpc.Creds(tlsCredentials))
	}

	grpcOptions = append(
		grpcOptions,
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(interceptorLogger(logger), logOptions...),
			withUnaryServerInterceptor(cfg.Tokens),
			grpc_recovery.UnaryServerInterceptor(recOptions...),
		),
	)

	grpcSrv := grpc.NewServer(grpcOptions...)

	device := handlers.NewDeviceImpl(ctx, logger, deviceService)
	peers := handlers.NewPeersImpl(ctx, logger, peerService)

	wgpb.RegisterDeviceServiceServer(grpcSrv, device)
	wgpb.RegisterPeerServiceServer(grpcSrv, peers)
	reflection.Register(grpcSrv)

	gatewayOptions := []runtime.ServeMuxOption{
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	}

	mux := http.NewServeMux()
	gwmux := runtime.NewServeMux(gatewayOptions...)

	if cfg.ServeSwagger {
		if err := serveSwagger(mux); err != nil {
			logger.Error("failed to serve swagger ui dist", zap.Error(err))
			return nil, err
		}
	}

	mux.Handle(defaultGatewayPrefix, gwmux)

	grpcDialOpts := make([]grpc.DialOption, 0, 1)

	if cfg.CaCert != "" {
		tlsCredentials, err := certs.LoadCATLSCredentials(cfg.CaCert)
		if err != nil {
			return nil, fmt.Errorf("rest gateway server: %w", err)
		}
		grpcDialOpts = append(grpcDialOpts, grpc.WithTransportCredentials(tlsCredentials))
	} else {
		grpcDialOpts = append(grpcDialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	if err := wgpb.RegisterDeviceServiceHandlerFromEndpoint(ctx, gwmux, addr, grpcDialOpts); err != nil {
		logger.Error("failed to register peer gateway handler", zap.Error(err))
		return nil, err
	}

	if err := wgpb.RegisterPeerServiceHandlerFromEndpoint(ctx, gwmux, addr, grpcDialOpts); err != nil {
		logger.Error("failed to register device gateway handler", zap.Error(err))
		return nil, err
	}

	return &Server{
		swagger: cfg.ServeSwagger,
		tls:     cfg.Cert != "" && cfg.Key != "",
		addr:    addr,
		cert:    cfg.Cert,
		key:     cfg.Key,
		logger:  logger,
		grpcSrv: grpcSrv,
		mux:     mux,
	}, nil
}

func (s *Server) Run(ctx context.Context, stop context.CancelFunc) {
	var wg sync.WaitGroup
	wg.Add(1)

	srvh2 := &http2.Server{}
	srvh := &http.Server{
		Addr:              s.addr,
		Handler:           grpcHandlerFunc(srvh2, s.grpcSrv, s.mux),
		ReadHeaderTimeout: time.Second * 10,
	}

	if err := http2.ConfigureServer(srvh, srvh2); err != nil {
		log.Fatalf("failed to configure server: %s\n", err)
	}

	go func() {
		log.Printf("server is up and running at %s", s.addr)
		defer wg.Done()

		if s.swagger {
			var proto string
			if s.tls {
				proto = "https"
			} else {
				proto = "http"
			}
			log.Printf("swagger docs available at %s://%s/swagger-ui\n", proto, s.addr)
		}

		if s.tls {
			if err := srvh.ListenAndServeTLS(s.cert, s.key); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}

		if err := srvh.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()
	log.Print("received signal")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srvh.Shutdown(ctx); err != nil {
		//nolint:gocritic
		log.Fatalf("server shutdown failed: %v\n", err)
	}

	wg.Wait()
	os.Exit(1)
}
