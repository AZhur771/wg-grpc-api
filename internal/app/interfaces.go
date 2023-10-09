//go:generate mockgen -source=interfaces.go -destination=./mocks/app_mocks.go -package=app_mocks
package app

import (
	"context"
	"database/sql"

	"github.com/AZhur771/wg-grpc-api/internal/dto"
	"github.com/AZhur771/wg-grpc-api/internal/entity"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	fieldmask_utils "github.com/mennanov/fieldmask-utils"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type WgCtrl interface {
	Close() error
	Device(name string) (*wgtypes.Device, error)
	Devices() ([]*wgtypes.Device, error)
	ConfigureDevice(name string, cfg wgtypes.Config) error
}

type PeerService interface {
	Add(ctx context.Context, dt dto.AddPeerDTO) (*entity.Peer, error)
	Update(ctx context.Context, dt dto.UpdatePeerDTO, mask fieldmask_utils.Mask) (*entity.Peer, error)
	Remove(ctx context.Context, id uuid.UUID) error
	Get(ctx context.Context, id uuid.UUID) (*entity.Peer, error)
	GetAll(ctx context.Context, dt dto.GetPeersRequestDTO) (dto.GetPeersResponseDTO, error)
	Enable(ctx context.Context, id uuid.UUID) error
	Disable(ctx context.Context, id uuid.UUID) error
	DownloadConfig(ctx context.Context, id uuid.UUID) (dto.DownloadFileDTO, error)
	DownloadQRCode(ctx context.Context, id uuid.UUID) (dto.DownloadFileDTO, error)
}

type DeviceService interface {
	Add(ctx context.Context, dt dto.AddDeviceDTO) (*entity.Device, error)
	Update(ctx context.Context, dt dto.UpdateDeviceDTO, mask fieldmask_utils.Mask) (*entity.Device, error)
	Remove(ctx context.Context, id uuid.UUID) error
	Get(ctx context.Context, id uuid.UUID) (*entity.Device, error)
	GetAll(ctx context.Context, dt dto.GetDevicesRequestDTO) (dto.GetDevicesResponseDTO, error)
	ConfigureDevice(device string, config wgtypes.PeerConfig) error
	GetConfiguredPeer(dev string, publicKey wgtypes.Key) (wgtypes.Peer, error)
	GetConfiguredPeers(dev string) ([]wgtypes.Peer, error)
}

type PeerRepo interface {
	Add(ctx context.Context, tx *sqlx.Tx, peer *entity.Peer) (*entity.Peer, error)
	Update(ctx context.Context, tx *sqlx.Tx, peer *entity.Peer) (*entity.Peer, error)
	Remove(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) error
	Get(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) (*entity.Peer, error)
	GetAll(ctx context.Context, tx *sqlx.Tx, skip, limit int, search string, deviceID uuid.UUID) ([]*entity.Peer, error)
	Count(ctx context.Context, tx *sqlx.Tx, deviceID uuid.UUID) (int, error)
	BeginTxx(ctx context.Context, options *sql.TxOptions) (*sqlx.Tx, error)
}

type DeviceRepo interface {
	Add(ctx context.Context, tx *sqlx.Tx, dev *entity.Device) (*entity.Device, error)
	Update(ctx context.Context, tx *sqlx.Tx, dev *entity.Device) (*entity.Device, error)
	Remove(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) error
	Get(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) (*entity.Device, error)
	GetAll(ctx context.Context, tx *sqlx.Tx, skip, limit int, search string) ([]*entity.Device, error)
	Count(ctx context.Context, tx *sqlx.Tx) (int, error)
	GenerateAddress(ctx context.Context, tx *sqlx.Tx, dev *entity.Device) (string, error)
	BeginTxx(ctx context.Context, options *sql.TxOptions) (*sqlx.Tx, error)
}
