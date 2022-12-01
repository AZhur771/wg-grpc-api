package storage

import (
	"context"

	"github.com/AZhur771/wg-grpc-api/internal/entity"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interface.go -destination=./mocks/peer_storage_test.go -package=storage_test
type PeerStorage interface {
	Create(ctx context.Context, peer *entity.PersistedPeer) (uuid.UUID, error)
	Update(ctx context.Context, id uuid.UUID, peer *entity.PersistedPeer) error
	Delete(ctx context.Context, id uuid.UUID) (*entity.PersistedPeer, error)
	Get(ctx context.Context, id uuid.UUID) (*entity.PersistedPeer, error)
	GetAll(ctx context.Context, limit, skip int) ([]*entity.PersistedPeer, error)
}
