package device_service

import "errors"

var (
	ErrPeerNotConfigured = errors.New("peer not configured")
	ErrInvalidPeer       = errors.New("invalid peer")
)
