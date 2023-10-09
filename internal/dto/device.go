package dto

import (
	"time"

	"github.com/AZhur771/wg-grpc-api/internal/entity"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
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

func (p *GetDevicesRequestDTO) IsValid() []*errdetails.BadRequest_FieldViolation {
	errors := make([]*errdetails.BadRequest_FieldViolation, 0)

	if p.Skip < 0 {
		errors = append(errors, &errdetails.BadRequest_FieldViolation{
			Field:       "skip",
			Description: "skip should not be less than zero",
		})
	}

	if p.Limit < 0 {
		errors = append(errors, &errdetails.BadRequest_FieldViolation{
			Field:       "limit",
			Description: "limit should not be less than zero",
		})
	}

	if p.Limit > 100 {
		errors = append(errors, &errdetails.BadRequest_FieldViolation{
			Field:       "limit",
			Description: "limit should not be more than 100",
		})
	}

	return errors
}
