package restserver

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"time"

	wgpb "github.com/AZhur771/wg-grpc-api/gen"
	"github.com/AZhur771/wg-grpc-api/internal/certs"
	"github.com/AZhur771/wg-grpc-api/third_party"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	defaultGatewayPrefix string = "/api/"
	defaultSwaggerPrefix string = "/swagger-ui/"
)

var (
	readHeaderTimeout = 10 * time.Second
	writeTimeout      = 10 * time.Second
	idleTimeout       = 10 * time.Second
)

type Server struct {
	logger  *zap.Logger
	gateway *http.Server
}

// serve swagger ui dist from third_party/openapiv2 folder.
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

func New(ctx context.Context, logger *zap.Logger, swagger bool,
	addr, grpcAddr, host, caCertPath, certPath, keyPath string,
	certOpt tls.ClientAuthType,
) (*Server, error) {
	mux := http.NewServeMux()

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

	gwmux := runtime.NewServeMux(gatewayOptions...)

	if swagger {
		if err := serveSwagger(mux); err != nil {
			logger.Error("failed to serve swagger ui dist", zap.Error(err))
			return nil, err
		}
	}

	mux.Handle(defaultGatewayPrefix, gwmux)

	opts := make([]grpc.DialOption, 0, 1)

	if certPath != "" || keyPath != "" {
		tlsCredentials, err := certs.LoadTLSCredentials(host, caCertPath, certPath, keyPath, certOpt)
		if err != nil {
			return nil, fmt.Errorf("rest gateway server: %w", err)
		}

		opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	if err := wgpb.RegisterPeerServiceHandlerFromEndpoint(ctx, gwmux, grpcAddr, opts); err != nil {
		logger.Error("failed to register peer gateway handler", zap.Error(err))
		return nil, err
	}

	if err := wgpb.RegisterDeviceServiceHandlerFromEndpoint(ctx, gwmux, grpcAddr, opts); err != nil {
		logger.Error("failed to register device gateway handler", zap.Error(err))
		return nil, err
	}

	gwSrv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	return &Server{
		logger:  logger,
		gateway: gwSrv,
	}, nil
}

func (s *Server) Start() error {
	return s.gateway.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.gateway.Shutdown(ctx)
}
