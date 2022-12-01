package service

import (
	"net"
	"sync"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type deviceConfig struct {
	name      string
	net       net.IPNet
	endpoint  net.UDPAddr
	publicKey wgtypes.Key

	sync.RWMutex
	totalPeers int
}

func (d *deviceConfig) incPeerCount() {
	d.Lock()
	defer d.Unlock()

	d.totalPeers++
}

func (d *deviceConfig) decPeerCount() {
	d.Lock()
	defer d.Unlock()

	d.totalPeers--
}

func (d *deviceConfig) getPeerCount() int {
	d.RLock()
	defer d.RUnlock()

	return d.totalPeers
}
