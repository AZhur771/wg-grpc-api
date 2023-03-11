//go:generate mockgen -source=interfaces.go -destination=./mocks/app_mocks.go -package=app_mocks
package app

import (
	"context"

	"github.com/AZhur771/wg-grpc-api/internal/dto"
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

type PeerService interface {
	Add(ctx context.Context, addPeerDTO dto.AddPeerDTO) (*entity.Peer, error)
	Update(ctx context.Context, updatePeerDTO dto.UpdatePeerDTO) (*entity.Peer, error)
	Remove(ctx context.Context, id uuid.UUID) error
	Get(ctx context.Context, id uuid.UUID) (*entity.Peer, error)
	GetAll(ctx context.Context, getPeersDTO dto.GetPeersRequestDTO) (dto.GetPeersResponseDTO, error)
	Enable(ctx context.Context, id uuid.UUID) error
	Disable(ctx context.Context, id uuid.UUID) error
	DownloadConfig(ctx context.Context, id uuid.UUID) (dto.DownloadFileDTO, error)
	DownloadQRCode(ctx context.Context, id uuid.UUID) (dto.DownloadFileDTO, error)
}

type DeviceService interface {
	// sync device with peers
	Setup(ctx context.Context, name string, endpoint string, address string) error
	GetDevice() (*entity.Device, error)

	// peer configuration
	AddPeer(peer wgtypes.PeerConfig) error
	UpdatePeer(peer wgtypes.PeerConfig) error
	RemovePeer(peer wgtypes.PeerConfig) error
	GetPeer(publicKey wgtypes.Key) (wgtypes.Peer, error)
	GetPeers() ([]wgtypes.Peer, error)
}

type PeerStorage interface {
	Add(ctx context.Context, peer *entity.Peer) (*entity.Peer, error)
	Update(ctx context.Context, peer *entity.Peer) (*entity.Peer, error)
	Remove(ctx context.Context, id uuid.UUID) (*entity.Peer, error)
	Get(ctx context.Context, id uuid.UUID) (*entity.Peer, error)
	GetAll(ctx context.Context, skip, limit int) ([]*entity.Peer, error)
	CountAll(ctx context.Context) (int, error)
}
