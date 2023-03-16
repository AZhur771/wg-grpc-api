package peerservice

import "errors"

var (
	ErrInvalidPeer             = errors.New("invalid peer")
	ErrInvalidPaginationParams = errors.New("invalid pagination params")
)