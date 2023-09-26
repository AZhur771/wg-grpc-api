package peerservice

import (
	"errors"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

var ErrInvalidPaginationParams = errors.New("invalid pagination params")

type ErrInvalidPeer struct {
	err     error
	details []*errdetails.BadRequest_FieldViolation
}

func NewErrInvalidPeer(err error, details []*errdetails.BadRequest_FieldViolation) ErrInvalidPeer {
	return ErrInvalidPeer{
		err:     err,
		details: details,
	}
}

func (e ErrInvalidPeer) Error() string {
	return e.err.Error()
}

func (e *ErrInvalidPeer) Details() *errdetails.BadRequest {
	br := &errdetails.BadRequest{}
	br.FieldViolations = append(br.FieldViolations, e.details...)

	return br
}
