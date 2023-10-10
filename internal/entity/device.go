package entity

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type Device struct {
	ID                  uuid.UUID
	Name                string
	Description         string
	Type                wgtypes.DeviceType
	PrivateKey          wgtypes.Key
	PublicKey           wgtypes.Key
	FirewallMark        int
	MaxPeersCount       int
	CurrentPeersCount   int
	Endpoint            string
	Address             string
	Table               string
	MTU                 int
	DNS                 string
	PersistentKeepAlive time.Duration
	PreUp               string
	PreDown             string
	PostUp              string
	PostDown            string
}

func (d *Device) IsValid() []*errdetails.BadRequest_FieldViolation {
	errors := make([]*errdetails.BadRequest_FieldViolation, 0)

	if d.Endpoint == "" {
		errors = append(errors, &errdetails.BadRequest_FieldViolation{
			Field:       "endpoint",
			Description: "endpoint should not be empty",
		})
	}

	if d.Address == "" {
		errors = append(errors, &errdetails.BadRequest_FieldViolation{
			Field:       "address",
			Description: "address should not be empty",
		})
	}

	_, err := net.ResolveUDPAddr("udp", d.Endpoint)
	if err != nil {
		errors = append(errors, &errdetails.BadRequest_FieldViolation{
			Field:       "endpoint",
			Description: "endpoint is not a valid IPv4 or IPv6 with port",
		})
	}

	if d.Address == "" {
		errors = append(errors, &errdetails.BadRequest_FieldViolation{
			Field:       "address",
			Description: "address should not be empty",
		})
	}

	_, _, err = net.ParseCIDR(d.Address)
	if err != nil {
		errors = append(errors, &errdetails.BadRequest_FieldViolation{
			Field:       "address",
			Description: "address is not a valid CIDR IP address",
		})
	}

	if d.DNS != "" {
		for _, addr := range strings.Split(d.DNS, ",") {
			ip := net.ParseIP(strings.TrimSpace(addr))
			if ip == nil {
				errors = append(errors, &errdetails.BadRequest_FieldViolation{
					Field:       "dns",
					Description: fmt.Sprintf("wrong DNS address: %s", addr),
				})
			}
		}
	}

	if len(d.Description) > 40 {
		errors = append(errors, &errdetails.BadRequest_FieldViolation{
			Field:       "description",
			Description: "description should be 40 characters max",
		})
	}

	return errors
}

func (d *Device) PopulateDynamicFields(wgdevice *wgtypes.Device) (*Device, error) {
	d.Type = wgdevice.Type
	d.FirewallMark = wgdevice.FirewallMark
	d.CurrentPeersCount = len(wgdevice.Peers)

	_, deviceNet, err := net.ParseCIDR(d.Address)
	if err != nil {
		return d, fmt.Errorf("device service: %w", err)
	}
	d.MaxPeersCount = computeMaxPeers(deviceNet.Mask)

	return d, nil
}

// computeMaxPeers computes max peers from mask.
func computeMaxPeers(mask net.IPMask) int {
	ones, _ := mask.Size()
	return (2 << (32 - ones - 1)) - 3
}
