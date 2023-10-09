package deviceservice

import (
	"errors"
)

var (
	ErrPeerNotConfigured       = errors.New("peer not configured")
	ErrInvalidPaginationParams = errors.New("invalid pagination params")
	ErrInvalidDeviceData       = errors.New("invalid device data")
)
