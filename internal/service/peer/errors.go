package peerservice

import (
	"errors"
)

var (
	ErrInvalidPaginationParams = errors.New("invalid pagination params")
	ErrInvalidPeerData         = errors.New("invalid peer data")
)
