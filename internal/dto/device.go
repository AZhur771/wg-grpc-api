package dto

import (
	"time"

	"github.com/AZhur771/wg-grpc-api/internal/entity"
	"github.com/google/uuid"
)

type AddDeviceDTO struct {
	Description         string
	Endpoint            string
	FirewallMark        int
	Address             string
	Table               string
	MTU                 int
	DNS                 string
	PersistentKeepAlive time.Duration
	PreUp               string
	PreDown             string
	PostUp              string
	PostDown            string
}

type UpdateDeviceDTO struct {
	ID                  uuid.UUID
	Description         string
	Endpoint            string
	FirewallMark        int
	Address             string
	Table               string
	MTU                 int
	DNS                 string
	PersistentKeepAlive time.Duration
	PreUp               string
	PreDown             string
	PostUp              string
	PostDown            string
}

type GetDevicesResponseDTO struct {
	Devices []*entity.Device
	Total   int
	HasNext bool
}

type GetDevicesRequestDTO struct {
	Skip   int
	Limit  int
	Search string
}

func (p *GetDevicesRequestDTO) IsValid() bool {
	if p.Skip < 0 || p.Limit < 0 {
		return false
	}

	return true
}
