package entity

import (
	"fmt"
	"net"
	"sync"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// Device extends wgtypes.Device (omitting peers).
type Device struct {
	Name              string
	Type              wgtypes.DeviceType
	PrivateKey        wgtypes.Key
	PublicKey         wgtypes.Key
	ListenPort        int
	FirewallMark      int
	MaxPeersCount     int
	CurrentPeersCount int
	Endpoint          string
	Address           string

	sync.RWMutex
	reservedIPs []net.IP
}

func (d *Device) IsValid() error {
	if d.Endpoint == "" {
		return fmt.Errorf("device: endpoint should not be empty")
	}

	res := net.ParseIP(d.Endpoint)
	if res == nil {
		return fmt.Errorf("device: endpoint is not a valid IPv4 or IPv6")
	}

	if d.Address == "" {
		return fmt.Errorf("device: address should not be empty")
	}

	_, _, err := net.ParseCIDR(d.Address)
	if err != nil {
		return fmt.Errorf("device: address is not a valid CIDR IP address")
	}

	return nil
}

func (d *Device) ComputeMaxPeersCount() error {
	_, deviceNet, err := net.ParseCIDR(d.Address)
	if err != nil {
		return fmt.Errorf("device service: %w", err)
	}

	d.MaxPeersCount = computeMaxPeers(deviceNet.Mask)

	return nil
}

func (d *Device) ComputeInitialReservedIPs() error {
	deviceIP, deviceNet, err := net.ParseCIDR(d.Address)
	if err != nil {
		return fmt.Errorf("device service: %w", err)
	}

	// add device address
	if err := d.ReserveIP(deviceIP); err != nil {
		return fmt.Errorf("peer service: %w", err)
	}

	// add broadcast address
	if err := d.ReserveIP(getBroadcastAddr(deviceNet)); err != nil {
		return fmt.Errorf("peer service: %w", err)
	}

	return nil
}

func (d *Device) PopulateDynamicFields(wgdevice *wgtypes.Device) *Device {
	d.CurrentPeersCount = len(wgdevice.Peers)

	return d
}

func (d *Device) ReserveIP(ip net.IP) error {
	d.Lock()
	defer d.Unlock()

	d.reservedIPs = append(d.reservedIPs, ip)

	return nil
}

func (d *Device) ReserveManyIPs(ips []net.IP) error {
	d.Lock()
	defer d.Unlock()

	d.reservedIPs = append(d.reservedIPs, ips...)

	return nil
}

func (d *Device) ReleaseIP(ip net.IP) error {
	d.Lock()
	defer d.Unlock()

	idx := -1

	for i, rip := range d.reservedIPs {
		if rip.Equal(ip) {
			idx = i
			break
		}
	}

	if idx == -1 {
		return ErrIPNotFound
	}

	d.reservedIPs = append(d.reservedIPs[:idx], d.reservedIPs[idx+1:]...)
	return nil
}

func (d *Device) IsReservedIP(ip net.IP) bool {
	d.RLock()
	defer d.RUnlock()

	for _, rip := range d.reservedIPs {
		if rip.Equal(ip) {
			return true
		}
	}

	return false
}

func (d *Device) GetAvailableIP() (*net.IPNet, error) {
	_, deviceIPNet, err := net.ParseCIDR(d.Address)
	if err != nil {
		return nil, fmt.Errorf("device service: %w", err)
	}

	for ip := deviceIPNet.IP.Mask(deviceIPNet.Mask); deviceIPNet.Contains(ip); incAddr(ip) {
		// check that incremented ip is not reserved
		isReserved := d.IsReservedIP(ip)
		if isReserved {
			continue
		}

		if err := d.ReserveIP(ip); err != nil {
			return nil, fmt.Errorf("device service: %w", err)
		}

		return &net.IPNet{
			IP:   ip,
			Mask: net.IPv4Mask(255, 255, 255, 255),
		}, nil
	}

	return nil, fmt.Errorf("device service: %w", ErrRunOutOfAddresses)
}

// incAddr allows getting next ip address.
func incAddr(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// broadcastAddr returns the last address in the given network, or the broadcast address.
func getBroadcastAddr(n *net.IPNet) net.IP {
	// TODO: leave link to original solution
	var broadcast net.IP
	if len(n.IP) == 4 {
		broadcast = net.ParseIP("0.0.0.0").To4()
	} else {
		broadcast = net.ParseIP("::")
	}
	for i := 0; i < len(n.IP); i++ {
		broadcast[i] = n.IP[i] | ^n.Mask[i]
	}
	return broadcast
}

// computeMaxPeers computes max peers from mask
func computeMaxPeers(mask net.IPMask) int {
	ones, _ := mask.Size()
	return (2 << (32 - ones - 1)) - 1
}
