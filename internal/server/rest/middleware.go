package restserver

import (
	"net/http"
	"strings"
	"time"

	"github.com/AZhur771/wg-grpc-api/internal/app"
	"go.uber.org/zap"
)

type responseObserver struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}

func (o *responseObserver) Write(p []byte) (n int, err error) {
	if !o.wroteHeader {
		o.WriteHeader(http.StatusOK)
	}
	n, err = o.ResponseWriter.Write(p)
	o.written += int64(n)
	return
}

func (o *responseObserver) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	if o.wroteHeader {
		return
	}
	o.wroteHeader = true
	o.status = code
}

func loggingMiddleware(next http.Handler, logger app.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestAccepted := time.Now()
		o := &responseObserver{ResponseWriter: w}
		next.ServeHTTP(o, r)
		latency := time.Since(requestAccepted)
		addr := r.RemoteAddr
		if i := strings.LastIndex(addr, ":"); i != -1 {
			addr = addr[:i]
		}
		logger.Info(
			"HTTP request",
			zap.String("address", addr),
			zap.String("method", r.Method),
			zap.String("URI", r.RequestURI),
			zap.String("protocol", r.Proto),
			zap.Int("status", o.status),
			zap.Duration("latency", latency),
		)
	})
}
