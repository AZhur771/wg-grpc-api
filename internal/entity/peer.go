package entity

import "time"

type Peer struct {
	PublicKey           string
	HasPresharedKey     bool
	Endpoint            string
	PersistentKeepAlive string
	LastHandshake       time.Time
	ReceiveBytes        int64
	TransmitBytes       int64
	AllowedIPs          []string
	ProtocolVersion     int
	IsEnabled           bool
	IsActive            bool
}
