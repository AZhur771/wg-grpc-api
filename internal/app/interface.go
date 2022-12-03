//go:generate mockgen -source=interface.go -destination=./mocks/app_mocks.go -package=app_mocks
package app

import (
	"context"
	"time"

	"github.com/AZhur771/wg-grpc-api/internal/entity"
	"github.com/google/uuid"
	"go.uber.org/zap/zapcore"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type Logger interface {
	Debug(string, ...zapcore.Field)
	Info(string, ...zapcore.Field)
	Warn(string, ...zapcore.Field)
	Error(string, ...zapcore.Field)
}

type Wg interface {
	GetDevice() (*wgtypes.Device, error)
	ConfigureDevice(config wgtypes.PeerConfig) error
	GetPeer(publicKey wgtypes.Key) (wgtypes.Peer, error)
	GetPeers() ([]wgtypes.Peer, error)
}

type PeerService interface {
	Add(ctx context.Context, addPresharedKey bool,
		persistentKeepAlive time.Duration, description string) (uuid.UUID, error)
	Update(ctx context.Context, id uuid.UUID,
		addPresharedKey bool, persistentKeepAlive time.Duration, description string, updateMask []string) error
	Delete(ctx context.Context, id uuid.UUID) error
	Get(ctx context.Context, id uuid.UUID) (*entity.Peer, error)
	GetAll(ctx context.Context) ([]*entity.Peer, error)
	Enable(ctx context.Context, id uuid.UUID) error
	Disable(ctx context.Context, id uuid.UUID) error
	DownloadConfig(ctx context.Context, id uuid.UUID) ([]byte, error)
	Setup(ctx context.Context, deviceName, deviceAddress, deviceEndpoint, peerFolder string) error
}

type PeerStorage interface {
	Create(ctx context.Context, peer *entity.PersistedPeer) (uuid.UUID, error)
	Update(ctx context.Context, id uuid.UUID, peer *entity.PersistedPeer) error
	Delete(ctx context.Context, id uuid.UUID) (*entity.PersistedPeer, error)
	Get(ctx context.Context, id uuid.UUID) (*entity.PersistedPeer, error)
	GetAll(ctx context.Context) ([]*entity.PersistedPeer, error)
}
