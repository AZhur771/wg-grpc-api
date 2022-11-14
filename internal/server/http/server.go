package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/AZhur771/wg-grpc-api/api/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

const gatewayPrefix string = "/api"

var (
	readHeaderTimeout = 10 * time.Second
	writeTimeout      = 10 * time.Second
	idleTimeout       = 10 * time.Second
)

type ServerImpl struct {
	gateway *http.Server
}

// NewServer returns new grpc Gateway Server.
func NewServer(ctx context.Context, host string, port int, conn *grpc.ClientConn, enableSwagger bool) (*ServerImpl, error) {
	gwmux := runtime.NewServeMux()

	if err := peerpb.RegisterPeerServiceHandler(ctx, gwmux, conn); err != nil {
		return nil, err
	}

	oaHandler, err := GetOpenAPIHandler()
	if err != nil {
		return nil, err
	}

	gwSrv := &http.Server{
		Addr: fmt.Sprintf("%s:%d", host, port),
		Handler: loggingMiddleware(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.HasPrefix(r.URL.Path, gatewayPrefix) {
					gwmux.ServeHTTP(w, r)
					return
				} else if enableSwagger {
					oaHandler.ServeHTTP(w, r)
				}
			})),
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	return &ServerImpl{
		gateway: gwSrv,
	}, nil
}

func (s *ServerImpl) Start() error {
	return s.gateway.ListenAndServe()
}

func (s *ServerImpl) Stop(ctx context.Context) error {
	return s.gateway.Shutdown(ctx)
}
