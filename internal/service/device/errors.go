package deviceservice

import (
	"errors"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

var (
	ErrPeerNotConfigured       = errors.New("peer not configured")
	ErrInvalidPaginationParams = errors.New("invalid pagination params")
)

type ErrInvalidDevice struct {
	err     error
	details []*errdetails.BadRequest_FieldViolation
}

func NewErrInvalidDevice(err error, details []*errdetails.BadRequest_FieldViolation) ErrInvalidDevice {
	return ErrInvalidDevice{
		err:     err,
		details: details,
	}
}

func (e ErrInvalidDevice) Error() string {
	return e.err.Error()
}

func (e *ErrInvalidDevice) Details() *errdetails.BadRequest {
	br := &errdetails.BadRequest{}
	br.FieldViolations = append(br.FieldViolations, e.details...)

	return br
}
