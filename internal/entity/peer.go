package entity

import (
	"fmt"
	"net"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type Peer struct {
	ID                          uuid.UUID
	DeviceID                    uuid.UUID
	Name                        string
	Email                       string
	PrivateKey                  wgtypes.Key
	PublicKey                   wgtypes.Key
	PresharedKey                wgtypes.Key
	Endpoint                    *net.UDPAddr
	PersistentKeepaliveInterval time.Duration
	LastHandshakeTime           time.Time
	ReceiveBytes                int64
	TransmitBytes               int64
	AllowedIPs                  []string
	DNS                         string
	MTU                         int
	ProtocolVersion             int
	HasPresharedKey             bool
	IsEnabled                   bool
	IsActive                    bool
	Description                 string
}

func (p *Peer) IsValid() []*errdetails.BadRequest_FieldViolation {
	errors := make([]*errdetails.BadRequest_FieldViolation, 0)

	if len(p.Name) < 1 || len(p.Name) > 20 {
		errors = append(errors, &errdetails.BadRequest_FieldViolation{
			Field:       "name",
			Description: "name should be between 1 and 20 characters",
		})
	}

	if p.Email != "" {
		_, err := mail.ParseAddress(p.Email)
		if err != nil {
			errors = append(errors, &errdetails.BadRequest_FieldViolation{
				Field:       "email",
				Description: fmt.Sprintf("email %s is invalid", p.Email),
			})
		}
	}

	if p.DNS != "" {
		for _, addr := range strings.Split(p.DNS, ",") {
			ip := net.ParseIP(strings.TrimSpace(addr))
			if ip == nil {
				errors = append(errors, &errdetails.BadRequest_FieldViolation{
					Field:       "dns",
					Description: fmt.Sprintf("wrong DNS: %s", addr),
				})
			}
		}
	}

	if len(p.Description) > 40 {
		errors = append(errors, &errdetails.BadRequest_FieldViolation{
			Field:       "description",
			Description: "description should be 40 characters max",
		})
	}

	return errors
}

func (p *Peer) PopulateDynamicFields(wgpeer *wgtypes.Peer) *Peer {
	p.Endpoint = wgpeer.Endpoint
	p.LastHandshakeTime = wgpeer.LastHandshakeTime
	p.ReceiveBytes = wgpeer.ReceiveBytes
	p.TransmitBytes = wgpeer.TransmitBytes
	p.ProtocolVersion = wgpeer.ProtocolVersion
	p.IsActive = wgpeer.LastHandshakeTime.After(time.Now().Add(-2 * time.Minute))

	return p
}

func (p *Peer) ToPeerConfig(dev *Device) (*wgtypes.PeerConfig, error) {
	allowedIPs, err := parseAllowedIPs(p.AllowedIPs)
	if err != nil {
		return nil, fmt.Errorf("peer: %w", err)
	}

	conf := &wgtypes.PeerConfig{
		PublicKey:    p.PublicKey,
		PresharedKey: &p.PresharedKey,
		Endpoint:     p.Endpoint,
		AllowedIPs:   allowedIPs,
	}

	if dev.PersistentKeepAlive != 0 {
		conf.PersistentKeepaliveInterval = &dev.PersistentKeepAlive
	}

	return conf, nil
}

func parseAllowedIPs(allowedIPs []string) ([]net.IPNet, error) {
	res := make([]net.IPNet, 0, len(allowedIPs))

	for _, allowedIP := range allowedIPs {
		ip, _, err := net.ParseCIDR(allowedIP)
		if err != nil {
			return nil, err
		}
		res = append(res, net.IPNet{
			IP:   ip,
			Mask: []byte{255, 255, 255, 255},
		})
	}

	return res, nil
}
