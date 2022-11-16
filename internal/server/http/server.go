package restserver

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"time"

	"github.com/AZhur771/wg-grpc-api/third_party"
	"go.uber.org/zap"

	peerpb "github.com/AZhur771/wg-grpc-api/api/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

const defaultGatewayPrefix string = "/api/"
const defaultSwaggerPrefix string = "/swagger-ui/"

var (
	readHeaderTimeout = 10 * time.Second
	writeTimeout      = 10 * time.Second
	idleTimeout       = 10 * time.Second
)

type ServerImpl struct {
	logger  *zap.Logger
	gateway *http.Server
}

// serve swagger ui dist from third_party/openapiv2 folder
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
func NewServer(
	ctx context.Context,
	logger *zap.Logger,
	conn *grpc.ClientConn,
	host string,
	port int,
	enableSwagger bool,
) (*ServerImpl, error) {
	mux := http.NewServeMux()

	if enableSwagger {
		if err := serveSwagger(mux); err != nil {
			logger.Error("failed to serve swagger ui dist", zap.Error(err))
			return nil, err
		}
	}

	gwmux := runtime.NewServeMux()
	if err := peerpb.RegisterPeerServiceHandler(ctx, gwmux, conn); err != nil {
		logger.Error("failed to register gateway handler", zap.Error(err))
		return nil, err
	}

	mux.Handle(fmt.Sprintf("/%s", defaultGatewayPrefix), gwmux)

	gwSrv := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", host, port),
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
