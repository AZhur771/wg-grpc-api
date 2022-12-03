package restserver

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"time"

	peerpb "github.com/AZhur771/wg-grpc-api/gen"
	"github.com/AZhur771/wg-grpc-api/internal/app"
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

type ServerImpl struct {
	logger  app.Logger
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

// NewServer returns new grpc Gateway Server.
func NewServer(ctx context.Context, logger app.Logger, addr, grpcAddr string, swagger bool) (*ServerImpl, error) {
	mux := http.NewServeMux()
	gwmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			EmitUnpopulated: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}))

	if swagger {
		if err := serveSwagger(mux); err != nil {
			logger.Error("failed to serve swagger ui dist", zap.Error(err))
			return nil, err
		}
	}

	mux.Handle(defaultGatewayPrefix, gwmux)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := peerpb.RegisterPeerServiceHandlerFromEndpoint(ctx, gwmux, grpcAddr, opts); err != nil {
		logger.Error("failed to register gateway handler", zap.Error(err))
		return nil, err
	}

	gwSrv := &http.Server{
		Addr:              addr,
		Handler:           loggingMiddleware(mux, logger),
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	return &ServerImpl{
		logger:  logger,
		gateway: gwSrv,
	}, nil
}

func (s *ServerImpl) Start() error {
	return s.gateway.ListenAndServe()
}

func (s *ServerImpl) Stop(ctx context.Context) error {
	return s.gateway.Shutdown(ctx)
}
