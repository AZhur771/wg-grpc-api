package wg

import (
	"bytes"

	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type Wg struct {
	name string
	ctrl *wgctrl.Client
}

func New(name string) (*Wg, error) {
	ctrl, err := wgctrl.New()
	if err != nil {
		return nil, err
	}

	return &Wg{name, ctrl}, nil
}

func (w *Wg) GetDevice() (*wgtypes.Device, error) {
	return w.ctrl.Device(w.name)
}

func (w *Wg) ConfigureDevice(config wgtypes.PeerConfig) error {
	return w.ctrl.ConfigureDevice(
		w.name,
		wgtypes.Config{
			Peers: []wgtypes.PeerConfig{config},
		},
	)
}

func (w *Wg) GetPeer(publicKey wgtypes.Key) (wgtypes.Peer, error) {
	device, err := w.ctrl.Device(w.name)
	if err != nil {
		return wgtypes.Peer{}, ErrDeviceNotConfigured
	}

	for _, p := range device.Peers {
		if bytes.Equal(p.PublicKey[:], publicKey[:]) {
			return p, nil
		}
	}

	return wgtypes.Peer{}, ErrPeerNotConfigured
}

func (w *Wg) GetPeers() ([]wgtypes.Peer, error) {
	device, err := w.ctrl.Device(w.name)
	if err != nil {
		return nil, ErrDeviceNotConfigured
	}

	return device.Peers, err
}
