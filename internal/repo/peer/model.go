package peerrepo

import (
	"time"

	"github.com/AZhur771/wg-grpc-api/internal/entity"
	"github.com/google/uuid"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type PeerModel struct {
	ID                  uuid.UUID `db:"id" sql:",type:uuid"`
	DeviceID            uuid.UUID `db:"device_id" sql:",type:uuid"`
	PrivateKey          string    `db:"private_key"`
	PresharedKey        string    `db:"preshared_key"`
	Name                string
	Description         string
	Email               string
	DNS                 string `db:"dns"`
	Mtu                 int
	PersistentKeepAlive int  `db:"persistent_keep_alive"`
	IsEnabled           bool `db:"is_enabled"`
	AllowedIPs          []string
}

func NewModel() *PeerModel {
	return &PeerModel{}
}

func (p *PeerModel) FromEntity(peer *entity.Peer) *PeerModel {
	p.ID = peer.ID
	p.DeviceID = peer.DeviceID
	p.PrivateKey = peer.PrivateKey.String()
	p.Name = peer.Name
	p.Description = peer.Description
	p.Email = peer.Email
	p.DNS = peer.DNS
	p.Mtu = peer.MTU
	p.PersistentKeepAlive = int(peer.PersistentKeepaliveInterval)
	p.IsEnabled = peer.IsEnabled

	if peer.HasPresharedKey {
		p.PresharedKey = peer.PresharedKey.String()
	}

	return p
}

func (p *PeerModel) ToEntity() (*entity.Peer, error) {
	peer := &entity.Peer{}

	privateKey, err := wgtypes.ParseKey(p.PrivateKey)
	if err != nil {
		return peer, err
	}
	peer.PrivateKey = privateKey

	if p.PresharedKey != "" {
		presharedKey, err := wgtypes.ParseKey(p.PresharedKey)
		if err != nil {
			return peer, err
		}
		peer.PresharedKey = presharedKey
		peer.HasPresharedKey = true
	}

	peer.PublicKey = privateKey.PublicKey()
	peer.ID = p.ID
	peer.DeviceID = p.DeviceID
	peer.Name = p.Name
	peer.Description = p.Description
	peer.Email = p.Email
	peer.MTU = p.Mtu
	peer.DNS = p.DNS
	peer.PersistentKeepaliveInterval = time.Duration(p.PersistentKeepAlive)
	peer.IsEnabled = p.IsEnabled

	return peer, nil
}
