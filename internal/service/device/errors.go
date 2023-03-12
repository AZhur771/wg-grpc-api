package deviceservice

import "errors"

var (
	ErrPeerNotConfigured = errors.New("peer not configured")
	ErrInvalidPeer       = errors.New("invalid peer")
)
