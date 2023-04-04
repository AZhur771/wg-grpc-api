package entity

import (
	"fmt"
	"net"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// Peer extends wgtypes.Peer.
type Peer struct {
	ID                          uuid.UUID     `json:"id"`
	Name                        string        `json:"name"`
	Email                       string        `json:"email"`
	PrivateKey                  WgKey         `json:"privateKey"`
	PublicKey                   WgKey         `json:"publicKey"`
	PresharedKey                WgKey         `json:"presharedKey"`
	Endpoint                    *net.UDPAddr  `json:"-"`
	PersistentKeepaliveInterval time.Duration `json:"persistentKeepAliveInterval"`
	LastHandshakeTime           time.Time     `json:"-"`
	ReceiveBytes                int64         `json:"-"`
	TransmitBytes               int64         `json:"-"`
	AllowedIPs                  []string      `json:"allowedIps"`
	ProtocolVersion             int           `json:"-"`
	HasPresharedKey             bool          `json:"hasPresharedKey"`
	IsEnabled                   bool          `json:"isEnabled"`
	IsActive                    bool          `json:"-"`
	Description                 string        `json:"description"`
	Tags                        []string      `json:"tags"`
}

func (p *Peer) IsValid() error {
	if len(p.Name) < 1 || len(p.Name) > 20 {
		return fmt.Errorf("peer: name should be between 1 and 20 characters")
	}

	if p.Email != "" {
		_, err := mail.ParseAddress(p.Email)
		if err != nil {
			return fmt.Errorf("peer: email %s is invalid", p.Email)
		}
	}

	if len(p.Description) > 40 {
		return fmt.Errorf("peer: description should be 40 characters max")
	}

	for _, tag := range p.Tags {
		if len(tag) > 20 {
			return fmt.Errorf("peer: tag should be 20 characters max")
		}
	}

	return nil
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

func (p *Peer) ToPeerConfig() (*wgtypes.PeerConfig, error) {
	allowedIPs, err := parseAllowedIPs(p.AllowedIPs)
	if err != nil {
		return nil, fmt.Errorf("peer: %w", err)
	}

	wgPresharedKey := wgtypes.Key(p.PresharedKey)

	return &wgtypes.PeerConfig{
		PublicKey:                   wgtypes.Key(p.PublicKey),
		PresharedKey:                &wgPresharedKey,
		Endpoint:                    p.Endpoint,
		PersistentKeepaliveInterval: &p.PersistentKeepaliveInterval,
		AllowedIPs:                  allowedIPs,
	}, nil
}

func parseAllowedIPs(allowedIPs []string) ([]net.IPNet, error) {
	res := make([]net.IPNet, 0, len(allowedIPs))

	for _, allowedIP := range allowedIPs {
		_, ipnet, err := net.ParseCIDR(allowedIP)
		if err != nil {
			return nil, err
		}
		res = append(res, *ipnet)
	}

	return res, nil
}
