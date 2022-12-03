package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/AZhur771/wg-grpc-api/internal/app"
	"github.com/AZhur771/wg-grpc-api/internal/entity"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
)

type PeerStorage struct {
	logger app.Logger

	rdb *redis.Client
}

func NewPeerStorage(logger app.Logger, rdb *redis.Client) *PeerStorage {
	return &PeerStorage{
		logger: logger,
		rdb:    rdb,
	}
}

func (ps PeerStorage) generateKey(id string) string {
	return "wg:" + id
}

func (ps PeerStorage) Create(ctx context.Context, p *entity.PersistedPeer) (uuid.UUID, error) {
	id := uuid.New()
	p.ID = id.String()

	v, err := json.Marshal(p)
	if err != nil {
		return id, fmt.Errorf("peer storage: failed to marshal peer with id %s: %w", p.ID, err)
	}

	return id, ps.rdb.Set(ctx, ps.generateKey(p.ID), v, 0).Err()
}

func (ps PeerStorage) Update(ctx context.Context, id uuid.UUID, p *entity.PersistedPeer) error {
	p.ID = id.String()
	v, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("peer storage: failed to marshal peer with id %s: %w", p.ID, err)
	}

	return ps.rdb.Set(ctx, ps.generateKey(p.ID), v, 0).Err()
}

func (ps PeerStorage) Delete(ctx context.Context, id uuid.UUID) (*entity.PersistedPeer, error) {
	p, err := ps.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := ps.rdb.Del(ctx, ps.generateKey(p.ID)).Err(); err != nil {
		return nil, fmt.Errorf("peer storage: failed to delete peer with id %s: %w", p.ID, err)
	}

	return p, nil
}

func (ps PeerStorage) Get(ctx context.Context, id uuid.UUID) (*entity.PersistedPeer, error) {
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

func (ps PeerStorage) GetAll(ctx context.Context) ([]*entity.PersistedPeer, error) {
	peers := make([]*entity.PersistedPeer, 0)

	iter := ps.rdb.Scan(ctx, 0, "wg:*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		result, err := ps.rdb.Get(ctx, key).Result()
		if err != nil && errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("peer storage: redis error while getting peer with id %s: %w",
				strings.Replace(key, "wg:", "", 1), err)
		}

		p := &entity.PersistedPeer{}
		if err := json.Unmarshal([]byte(result), p); err != nil {
			return nil, fmt.Errorf("peer storage: failed to unmarshal peer with id %s: %w",
				strings.Replace(key, "wg:", "", 1), err)
		}

		peers = append(peers, p)
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("peer storage: failed to scan db: %w", err)
	}

	return peers, nil
}
