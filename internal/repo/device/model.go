package devicerepo

import (
	"time"

	"github.com/AZhur771/wg-grpc-api/internal/entity"
	"github.com/google/uuid"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type DeviceModel struct {
	ID                  uuid.UUID `db:"id" sql:",type:uuid"`
	Name                string
	PrivateKey          string `db:"private_key"`
	Description         string
	Endpoint            string
	FwMark              int `db:"fw_mark"`
	Address             string
	Mtu                 int
	DNS                 string `db:"dns"`
	PersistentKeepAlive int    `db:"persistent_keep_alive"`
	Tble                string
	PreUp               string `db:"pre_up"`
	PostUp              string `db:"post_up"`
	PreDown             string `db:"pre_down"`
	PostDown            string `db:"post_down"`
}

func NewModel() *DeviceModel {
	return &DeviceModel{}
}

func (d *DeviceModel) FromEntity(dev *entity.Device) *DeviceModel {
	d.ID = dev.ID
	d.Name = dev.Name
	d.PrivateKey = dev.PrivateKey.String()
	d.Description = dev.Description
	d.Endpoint = dev.Endpoint
	d.FwMark = dev.FirewallMark
	d.Address = dev.Address
	d.Mtu = dev.MTU
	d.DNS = dev.DNS
	d.PersistentKeepAlive = int(dev.PersistentKeepAlive) / (1000 * 1000 * 1000)
	d.Tble = dev.Table
	d.PreUp = dev.PreUp
	d.PostUp = dev.PostUp
	d.PreDown = dev.PreDown
	d.PostDown = dev.PostDown

	return d
}

func (d *DeviceModel) ToEntity() (*entity.Device, error) {
	dev := &entity.Device{}

	privateKey, err := wgtypes.ParseKey(d.PrivateKey)
	if err != nil {
		return dev, err
	}
	dev.PrivateKey = privateKey
	dev.PublicKey = privateKey.PublicKey()
	dev.ID = d.ID
	dev.Name = d.Name
	dev.Description = d.Description
	dev.Endpoint = d.Endpoint
	dev.FirewallMark = d.FwMark
	dev.Address = d.Address
	dev.MTU = d.Mtu
	dev.DNS = d.DNS
	dev.PersistentKeepAlive = time.Duration(d.PersistentKeepAlive) * time.Second
	dev.Table = d.Tble
	dev.PreUp = d.PreUp
	dev.PostUp = d.PostUp
	dev.PreDown = d.PreDown
	dev.PostDown = d.PostDown

	return dev, nil
}
