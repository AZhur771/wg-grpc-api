package httpserver

import (
	"log"
	"net/http"
	"strings"
	"time"
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

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestAccepted := time.Now()
		o := &responseObserver{ResponseWriter: w}
		next.ServeHTTP(o, r)
		latency := time.Since(requestAccepted)
		addr := r.RemoteAddr
		if i := strings.LastIndex(addr, ":"); i != -1 {
			addr = addr[:i]
		}
		log.Printf(
			"%s - %s - %s - %s - %d - %s - %s",
			addr, r.Method, r.RequestURI, r.Proto, o.status, latency, r.UserAgent(),
		)
	})
}
