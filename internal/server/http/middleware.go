package restserver

import (
	"net/http"
	"strings"
	"time"

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

func loggingMiddleware(next http.Handler, logger *zap.Logger) http.Handler {
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

func authTokenMiddleware(tokens ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Token "))

			if !stringInSlice(token, tokens) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func stringInSlice(s string, vv []string) bool {
	for _, v := range vv {
		if v == s {
			return true
		}
	}

	return false
}
