package entity

import "errors"

var (
	ErrRunOutOfAddresses = errors.New("run out of addresses")
	ErrIPNotFound        = errors.New("ip not found among reserved ips")
)
