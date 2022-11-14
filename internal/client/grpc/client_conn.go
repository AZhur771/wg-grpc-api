package grpcclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClientConn(ctx context.Context, host string, port int) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		ctx,
		fmt.Sprintf("%s:%d", host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}
