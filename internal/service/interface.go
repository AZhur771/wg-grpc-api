package service

import (
	"context"
	"time"

	"github.com/AZhur771/wg-grpc-api/internal/entity"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interface.go -destination=./mocks/peer_service_test.go -package=service_test
type PeerService interface {
	Add(ctx context.Context, addPresharedKey bool,
		persistentKeepAlive time.Duration, description string) (uuid.UUID, error)
	Update(ctx context.Context, id uuid.UUID,
		addPresharedKey bool, persistentKeepAlive time.Duration, description string, updateMask []string) error
	Delete(ctx context.Context, id uuid.UUID) error
	Get(ctx context.Context, id uuid.UUID) (*entity.Peer, error)
	GetAll(ctx context.Context, limit, skip int) (*entity.PaginatedPeers, error)
	DownloadConfig(ctx context.Context, id uuid.UUID) ([]byte, error)
	Setup(ctx context.Context, deviceName, deviceAddress, deviceEndpoint, peerFolder string) error
}
