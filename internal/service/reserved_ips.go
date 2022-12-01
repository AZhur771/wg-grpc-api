package service

import (
	"errors"
	"net"
	"sync"
)

var ErrIPNotFound = errors.New("ip not found among reserved ips")

type reservedIPs struct {
	sync.RWMutex
	ips []net.IP
}

func (r *reservedIPs) Add(ip net.IP) error {
	r.Lock()
	defer r.Unlock()

	r.ips = append(r.ips, ip)

	return nil
}

func (r *reservedIPs) Remove(ip net.IP) error {
	r.Lock()
	defer r.Unlock()

	idx := -1

	for i, rip := range r.ips {
		if rip.Equal(ip) {
			idx = i
			break
		}
	}

	if idx == -1 {
		return ErrIPNotFound
	}

	r.ips = append(r.ips[:idx], r.ips[idx+1:]...)
	return nil
}

func (r *reservedIPs) Contains(ip net.IP) bool {
	r.RLock()
	defer r.RUnlock()

	for _, rip := range r.ips {
		if rip.Equal(ip) {
			return true
		}
	}

	return false
}
