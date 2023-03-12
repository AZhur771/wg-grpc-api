package grpcserver

import (
	"context"
	"errors"
	"time"

	"github.com/AZhur771/wg-grpc-api/internal/app"
	"github.com/AZhur771/wg-grpc-api/internal/entity"
	deviceservice "github.com/AZhur771/wg-grpc-api/internal/service/device"
	peerstorage "github.com/AZhur771/wg-grpc-api/internal/storage/peer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func withUnaryServerInterceptor(logger app.Logger, tokens []string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		err := authorize(ctx, tokens)
		if err != nil {
			logFields(logger, info, start, err)
			return nil, err
		}

		h, err := handler(ctx, req)
		logFields(logger, info, start, err)

		return h, err
	}
}

func logFields(logger app.Logger, info *grpc.UnaryServerInfo, start time.Time, err error) {
	latency := time.Since(start)
	fields := []zapcore.Field{
		zap.String("grpc.code", getCode(err).String()),
		zap.String("grpc.method", info.FullMethod),
		zap.Float32("grpc.latency", durationToMilliseconds(latency)),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		logger.Error("GRPC REQUEST", fields...)
	} else {
		logger.Info("GRPC REQUEST", fields...)
	}
}

func authorize(ctx context.Context, tokens []string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["x-api-key"]
	if len(values) == 0 {
		return status.Error(codes.Unauthenticated, "x-api-key is not provided")
	}

	if contains(tokens, values[0]) {
		return nil
	}

	return status.Error(codes.PermissionDenied, "no permission to access this RPC")
}

func durationToMilliseconds(duration time.Duration) float32 {
	return float32(duration.Nanoseconds()/1000) / 1000
}

func getCode(err error) codes.Code {
	if errors.Is(err, peerstorage.ErrPeerNotFound) || errors.Is(err, entity.ErrIPNotFound) {
		return codes.NotFound
	}

	if errors.Is(err, entity.ErrRunOutOfAddresses) {
		return codes.ResourceExhausted
	}

	if errors.Is(err, deviceservice.ErrInvalidPeer) {
		return codes.InvalidArgument
	}

	return status.Code(err)
}

func contains(tokens []string, token string) bool {
	for _, t := range tokens {
		if token == t {
			return true
		}
	}

	return false
}
