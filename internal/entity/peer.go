package entity

import (
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// Peer extends wgtypes.Peer.
type Peer struct {
	// Unique identifier
	ID uuid.UUID

	// PrivateKey used to compute PublicKey
	PrivateKey wgtypes.Key

	// PublicKey is the public key of a peer, computed from its private key.
	//
	// PublicKey is always present in a Peer.
	PublicKey wgtypes.Key

	// PresharedKey is an optional preshared key which may be used as an
	// additional layer of security for peer communications.
	//
	// A zero-value Key means no preshared key is configured.
	PresharedKey wgtypes.Key

	// Endpoint is the most recent source address used for communication by
	// this Peer.
	Endpoint *net.UDPAddr

	// PersistentKeepaliveInterval specifies how often an "empty" packet is sent
	// to a peer to keep a connection alive.
	//
	// A value of 0 indicates that persistent keepalives are disabled.
	PersistentKeepaliveInterval time.Duration

	// LastHandshakeTime indicates the most recent time a handshake was performed
	// with this peer.
	//
	// A zero-value time.Time indicates that no handshake has taken place with
	// this peer.
	LastHandshakeTime time.Time

	// ReceiveBytes indicates the number of bytes received from this peer.
	ReceiveBytes int64

	// TransmitBytes indicates the number of bytes transmitted to this peer.
	TransmitBytes int64

	// AllowedIPs specifies which IPv4 and IPv6 addresses this peer is allowed
	// to communicate on.
	//
	// 0.0.0.0/0 indicates that all IPv4 addresses are allowed, and ::/0
	// indicates that all IPv6 addresses are allowed.
	AllowedIPs []net.IPNet

	// ProtocolVersion specifies which version of the WireGuard protocol is used
	// for this Peer.
	//
	// A value of 0 indicates that the most recent protocol version will be used.
	ProtocolVersion int

	// HasPresharedKey indicates whether peer has been configured with preshared key
	HasPresharedKey bool

	// IsEnabled indicates whether peer can receive/transmit bytes
	IsEnabled bool

	// IsActive indicates whether peer is currently connected to the vpn server
	// Peer config is considered active if LastHandshake happened less than two minutes ago
	IsActive bool

	// Short description
	Description string
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

func (p *Peer) ToPersistedPeer() *PersistedPeer {
	persistedPeer := &PersistedPeer{
		ID:                          p.ID.String(),
		PrivateKey:                  p.PrivateKey.String(),
		PublicKey:                   p.PublicKey.String(),
		PersistentKeepaliveInterval: p.PersistentKeepaliveInterval,
		AllowedIPs:                  convertAllowedIPsToStringSlice(p.AllowedIPs),
		HasPresharedKey:             p.HasPresharedKey,
		IsEnabled:                   p.IsEnabled,
		Description:                 p.Description,
	}

	if p.HasPresharedKey {
		persistedPeer.PresharedKey = p.PresharedKey.String()
	}

	return persistedPeer
}

// PersistedPeer used to store peer data in db.
type PersistedPeer struct {
	ID                          string        `json:"id"`
	PrivateKey                  string        `json:"privateKey"`
	PublicKey                   string        `json:"publicKey"`
	PresharedKey                string        `json:"presharedKey,omitempty"`
	PersistentKeepaliveInterval time.Duration `json:"persistentKeepaliveInterval,omitempty"`
	LastHandshakeTime           time.Time     `json:"lastHandshakeTime,omitempty"`
	AllowedIPs                  []string      `json:"allowedIPs,omitempty"`
	HasPresharedKey             bool          `json:"hasPresharedKey"`
	IsEnabled                   bool          `json:"isEnabled"`
	Description                 string        `json:"description,omitempty"`
}

func (p *PersistedPeer) ToPeer() (*Peer, error) {
	id, err := uuid.Parse(p.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse uuid: %w", err)
	}

	privateKey, err := wgtypes.ParseKey(p.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	publicKey, err := wgtypes.ParseKey(p.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	var presharedKey wgtypes.Key

	if p.HasPresharedKey {
		presharedKey, err = wgtypes.ParseKey(p.PresharedKey)
		if err != nil {
			return nil, fmt.Errorf("failed to parse preshared key: %w", err)
		}
	}

	allowedIPs, err := parseAllowedIPs(p.AllowedIPs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse allowed ips: %w", err)
	}

	return &Peer{
		ID:                          id,
		PrivateKey:                  privateKey,
		PublicKey:                   publicKey,
		PresharedKey:                presharedKey,
		PersistentKeepaliveInterval: p.PersistentKeepaliveInterval,
		AllowedIPs:                  allowedIPs,
		IsEnabled:                   p.IsEnabled,
		HasPresharedKey:             p.HasPresharedKey,
		Description:                 p.Description,
	}, nil
}

// TODO: duplicated code, should be removed.
func convertAllowedIPsToStringSlice(data []net.IPNet) []string {
	res := make([]string, 0, len(data))

	for _, e := range data {
		res = append(res, e.String())
	}

	return res
}

func parseAllowedIPs(data []string) ([]net.IPNet, error) {
	res := make([]net.IPNet, 0, len(data))

	for _, e := range data {
		_, ipnet, err := net.ParseCIDR(e)
		if err != nil {
			return nil, err
		}
		res = append(res, *ipnet)
	}

	return res, nil
}
