package common

import "google.golang.org/genproto/googleapis/rpc/errdetails"

type ErrInvalidData struct {
	err     error
	details []*errdetails.BadRequest_FieldViolation
}

func NewErrInvalidData(err error, details []*errdetails.BadRequest_FieldViolation) ErrInvalidData {
	return ErrInvalidData{
		err:     err,
		details: details,
	}
}

func (e ErrInvalidData) Error() string {
	return e.err.Error()
}

func (e *ErrInvalidData) Details() *errdetails.BadRequest {
	br := &errdetails.BadRequest{}
	br.FieldViolations = append(br.FieldViolations, e.details...)

	return br
}
