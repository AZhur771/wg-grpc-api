package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/AZhur771/wg-grpc-api/internal/entity"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type peerStorage struct {
	logger *zap.Logger
	rdb    *redis.Client
}

func (ps peerStorage) generateKey(id string) string {
	return "wg:" + id
}

func (ps peerStorage) Create(ctx context.Context, p *entity.PersistedPeer) (uuid.UUID, error) {
	id := uuid.New()
	p.ID = id.String()

	v, err := json.Marshal(p)
	if err != nil {
		return id, fmt.Errorf("peer storage: failed to marshal peer with id %s: %w", p.ID, err)
	}

	return id, ps.rdb.Set(ctx, ps.generateKey(p.ID), v, 0).Err()
}

func (ps peerStorage) Update(ctx context.Context, id uuid.UUID, p *entity.PersistedPeer) error {
	p.ID = id.String()
	v, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("peer storage: failed to marshal peer with id %s: %w", p.ID, err)
	}

	return ps.rdb.Set(ctx, ps.generateKey(p.ID), v, 0).Err()
}

func (ps peerStorage) Delete(ctx context.Context, id uuid.UUID) (*entity.PersistedPeer, error) {
	p, err := ps.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := ps.rdb.Del(ctx, ps.generateKey(p.ID)).Err(); err != nil {
		return nil, fmt.Errorf("peer storage: failed to delete peer with id %s: %w", p.ID, err)
	}

	return p, nil
}

func (ps peerStorage) Get(ctx context.Context, id uuid.UUID) (*entity.PersistedPeer, error) {
	result, err := ps.rdb.Get(ctx, ps.generateKey(id.String())).Result()
	if err != nil && errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("peer storage: redis error while getting peer with id %s: %w", id, err)
	}

	if result == "" {
		return nil, ErrPeerNotFound
	}

	p := &entity.PersistedPeer{}
	if err := json.Unmarshal([]byte(result), p); err != nil {
		return nil, fmt.Errorf("peer storage: failed to unmarshal peer with id %s: %w", id.String(), err)
	}

	return p, nil
}

func (ps peerStorage) GetAll(ctx context.Context, skip, limit int) ([]*entity.PersistedPeer, error) {
	peers := make([]*entity.PersistedPeer, 0)

	iter := ps.rdb.Scan(ctx, uint64(skip), "wg:*", int64(limit)).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		id := strings.Replace(key, "wg:", "", 1)
		result, err := ps.rdb.Get(ctx, key).Result()
		if err != nil && errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("peer storage: redis error while getting peer with id %s: %w", id, err)
		}

		p := &entity.PersistedPeer{}
		if err := json.Unmarshal([]byte(result), p); err != nil {
			return nil, fmt.Errorf("peer storage: failed to unmarshal peer with id %s: %w", id, err)
		}

		peers = append(peers, p)
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("peer storage: failed to scan db: %w", err)
	}

	return peers, nil
}

func NewPeerStorage(logger *zap.Logger, rdb *redis.Client) PeerStorage {
	return &peerStorage{
		logger: logger,
		rdb:    rdb,
	}
}
