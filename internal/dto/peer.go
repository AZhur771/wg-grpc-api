package dto

import (
	"time"

	"github.com/AZhur771/wg-grpc-api/internal/entity"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
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

func (p *GetPeersRequestDTO) IsValid() []*errdetails.BadRequest_FieldViolation {
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
