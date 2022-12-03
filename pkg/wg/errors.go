package wg

import "errors"

var (
	ErrPeerNotConfigured   = errors.New("peer not configured")
	ErrDeviceNotConfigured = errors.New("device not configured")
)
