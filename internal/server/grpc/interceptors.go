package grpcserver

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func getUnaryServerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		h, err := handler(ctx, req)
		latency := time.Since(start)

		logger.Info(
			"GRPC request",
			zap.String("method", info.FullMethod),
			zap.Duration("latency", latency),
			zap.Error(err),
		)

		return h, err
	}
}
