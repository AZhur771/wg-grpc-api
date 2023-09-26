package dto

import (
	"time"

	"github.com/AZhur771/wg-grpc-api/internal/entity"
	"github.com/google/uuid"
)

type AddPeerDTO struct {
	DeviceID            uuid.UUID
	Name                string
	Email               string
	AddPresharedKey     bool
	PersistentKeepAlive time.Duration
	Description         string
	DNS                 string
	MTU                 int
}

type UpdatePeerDTO struct {
	ID                  uuid.UUID
	Name                string
	Email               string
	AddPresharedKey     bool
	RemovePresharedKey  bool
	PersistentKeepAlive time.Duration
	Description         string
	DNS                 string
	MTU                 int
}

type DownloadFileDTO struct {
	Name string
	Size int64
	Data []byte
}

type GetPeersResponseDTO struct {
	Peers   []*entity.Peer
	Total   int
	HasNext bool
}

type GetPeersRequestDTO struct {
	DeviceID uuid.UUID
	Skip     int
	Limit    int
	Search   string
}

func (p *GetPeersRequestDTO) IsValid() bool {
	if p.Skip < 0 || p.Limit < 0 {
		return false
	}

	return true
}
