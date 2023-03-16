package peerstorage

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/AZhur771/wg-grpc-api/internal/app"
	"github.com/AZhur771/wg-grpc-api/internal/entity"
	"github.com/google/uuid"
	"github.com/mailru/easyjson"
)

type Storage struct {
	logger app.Logger
	dir    string
}

func NewPeerStorage(logger app.Logger, dir string) *Storage {
	return &Storage{
		logger: logger,
		dir:    dir,
	}
}

func (ps *Storage) createOrUpdateEntityFile(peer *entity.Peer) (*entity.Peer, error) {
	entityFile, err := os.Create(path.Join(ps.dir, fmt.Sprintf("%s.json", peer.ID.String())))
	if err != nil {
		return nil, fmt.Errorf("storage: %w", err)
	}
	defer entityFile.Close()

	rawBytes, err := easyjson.Marshal(peer)
	if err != nil {
		return nil, fmt.Errorf("storage: %w", err)
	}

	if _, err := entityFile.Write(rawBytes); err != nil {
		return nil, fmt.Errorf("storage: %w", err)
	}

	return peer, nil
}

func (ps *Storage) Add(ctx context.Context, peer *entity.Peer) (*entity.Peer, error) {
	peer.ID = uuid.New()

	return ps.createOrUpdateEntityFile(peer)
}

func (ps *Storage) Update(ctx context.Context, peer *entity.Peer) (*entity.Peer, error) {
	_, err := ps.Get(ctx, peer.ID)
	if err != nil {
		return nil, fmt.Errorf("storage: %w", ErrPeerNotFound)
	}

	return ps.createOrUpdateEntityFile(peer)
}

func (ps *Storage) Remove(ctx context.Context, id uuid.UUID) (*entity.Peer, error) {
	peer, err := ps.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("storage: %w", ErrPeerNotFound)
	}

	if err := os.Remove(path.Join(ps.dir, fmt.Sprintf("%s.json", id.String()))); err != nil {
		return nil, fmt.Errorf("storage: %w", err)
	}

	return peer, nil
}

func (ps *Storage) Get(ctx context.Context, id uuid.UUID) (*entity.Peer, error) {
	bytes, err := os.ReadFile(path.Join(ps.dir, fmt.Sprintf("%s.json", id.String())))
	if err != nil {
		return nil, fmt.Errorf("storage: %w", err)
	}

	peer := &entity.Peer{}
	if err := easyjson.Unmarshal(bytes, peer); err != nil {
		return nil, fmt.Errorf("storage: %w", err)
	}

	return peer, nil
}

func (ps *Storage) sliceFiles(files []fs.FileInfo, skip, limit int) []fs.FileInfo {
	if skip >= len(files) {
		return make([]fs.FileInfo, 0)
	}

	if limit == 0 || skip+limit >= len(files) {
		return files[skip:]
	}

	return files[skip : skip+limit]
}

func (ps *Storage) GetAll(ctx context.Context, skip, limit int) ([]*entity.Peer, error) {
	files, err := ioutil.ReadDir(ps.dir)
	if err != nil {
		return nil, fmt.Errorf("storage: %w", err)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().After(files[j].ModTime())
	})

	files = ps.sliceFiles(files, skip, limit)

	peers := make([]*entity.Peer, 0, len(files))

	for _, file := range files {
		id, err := uuid.Parse(strings.Replace(file.Name(), ".json", "", 1))
		if err != nil {
			return nil, fmt.Errorf("storage: %w", err)
		}

		peer, err := ps.Get(ctx, id)
		if err != nil {
			return nil, err
		}

		peers = append(peers, peer)
	}

	return peers, nil
}

func (ps *Storage) CountAll(ctx context.Context) (int, error) {
	files, err := ioutil.ReadDir(ps.dir)
	if err != nil {
		return 0, fmt.Errorf("storage: %w", err)
	}

	return len(files), nil
}
