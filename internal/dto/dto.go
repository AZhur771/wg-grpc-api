package dto

import (
	"time"

	"github.com/AZhur771/wg-grpc-api/internal/entity"
	"github.com/google/uuid"
)

type AddPeerDTO struct {
	Name                string
	Email               string
	AddPresharedKey     bool
	PersistentKeepAlive time.Duration
	Tags                []string
	Description         string
}

type UpdatePeerDTO struct {
	ID                  uuid.UUID
	Name                string
	Email               string
	AddPresharedKey     bool
	RemovePresharedKey  bool
	PersistentKeepAlive time.Duration
	Tags                []string
	Description         string
}

type DownloadFileDTO struct {
	Name string
	Size int64
	Data []byte
}

type GetPeersRequestDTO struct {
	Skip  int
	Limit int
}

func (gp *GetPeersRequestDTO) IsValid() bool {
	if gp.Skip < 0 || gp.Limit < 0 {
		return false
	}

	return true
}

type GetPeersResponseDTO struct {
	Peers   []*entity.Peer
	Total   int
	HasNext bool
}
