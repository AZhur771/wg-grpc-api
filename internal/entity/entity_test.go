package entity_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/AZhur771/wg-grpc-api/internal/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var (
	uuidMock         = "acda9b63-45ae-4352-995c-82202086cac4"
	privateKeyMock   = "WDhbZ4+4sE8LmIu4tSA1AXINX1ly+d+ZUwzazdiRMFU="
	publicKeyMock    = "MG+IiZS7uPsLMigRQoMch5MD7H2XCqEM+o9QJf1VGD4="
	presharedKeyMock = "MBj+hmcP54YPcXdq8odHD18lHx/7Y/G2e1x4ErLXCE8="

	allowedIPsMock        = []string{"10.0.0.2/32"}
	lastHandshakeTimeMock = time.Date(2018, time.May, 19, 1, 2, 3, 4, time.UTC)
)

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

func generateTestPeer() (*entity.Peer, error) {
	id, err := uuid.Parse(uuidMock)
	if err != nil {
		return nil, err
	}

	privateKey, err := wgtypes.ParseKey(privateKeyMock)
	if err != nil {
		return nil, err
	}

	publicKey, err := wgtypes.ParseKey(publicKeyMock)
	if err != nil {
		return nil, err
	}

	presharedKey, err := wgtypes.ParseKey(presharedKeyMock)
	if err != nil {
		return nil, err
	}

	return &entity.Peer{
		ID:                          id,
		PrivateKey:                  entity.WgKey(privateKey),
		PublicKey:                   entity.WgKey(publicKey),
		PresharedKey:                entity.WgKey(presharedKey),
		PersistentKeepaliveInterval: time.Second * 15,
		AllowedIPs:                  allowedIPsMock,
		HasPresharedKey:             true,
		IsEnabled:                   true,
		Description:                 "some description",
	}, nil
}

func generateTestDevice() (*entity.Device, error) {
	privateKey, err := wgtypes.ParseKey(privateKeyMock)
	if err != nil {
		return nil, err
	}

	publicKey, err := wgtypes.ParseKey(publicKeyMock)
	if err != nil {
		return nil, err
	}

	return &entity.Device{
		Name:          "wg0",
		Type:          wgtypes.LinuxKernel,
		PrivateKey:    privateKey,
		PublicKey:     publicKey,
		ListenPort:    51820,
		MaxPeersCount: 255,
	}, nil
}

func TestEntityPeer_PopulateDynamicFields(t *testing.T) {
	testPeer, err := generateTestPeer()
	require.NoError(t, err)

	endpoint, err := net.ResolveUDPAddr("udp", "203.0.113.0:51823")
	require.NoError(t, err)

	allowedIPs, err := parseAllowedIPs(testPeer.AllowedIPs)
	require.NoError(t, err)

	wgTestPeer := &wgtypes.Peer{
		PublicKey:                   wgtypes.Key(testPeer.PublicKey),
		PresharedKey:                wgtypes.Key(testPeer.PresharedKey),
		Endpoint:                    endpoint,
		PersistentKeepaliveInterval: time.Second * 15,
		LastHandshakeTime:           lastHandshakeTimeMock,
		ReceiveBytes:                111,
		TransmitBytes:               222,
		AllowedIPs:                  allowedIPs,
		ProtocolVersion:             1,
	}

	testPeer = testPeer.PopulateDynamicFields(wgTestPeer)
	require.Equal(t, testPeer.Endpoint, wgTestPeer.Endpoint)
	require.Equal(t, testPeer.LastHandshakeTime, wgTestPeer.LastHandshakeTime)
	require.Equal(t, testPeer.ReceiveBytes, wgTestPeer.ReceiveBytes)
	require.Equal(t, testPeer.TransmitBytes, wgTestPeer.TransmitBytes)
	require.Equal(t, testPeer.ProtocolVersion, wgTestPeer.ProtocolVersion)
}

func TestEntityPeer_IsValid(t *testing.T) {
	testPeer, err := generateTestPeer()
	require.NoError(t, err)

	require.Error(t, testPeer.IsValid(), "peer: name should be between 1 and 20 characters")

	testPeer.Name = "some_name"
	require.NoError(t, testPeer.IsValid())

	testPeer.Email = "invalid.email"
	require.Error(t, testPeer.IsValid(), fmt.Sprintf("peer: email %s is invalid", testPeer.Email))

	testPeer.Email = "valid.email@example.com"
	require.NoError(t, testPeer.IsValid())
}

func TestEntityDevice_PopulateDynamicFields(t *testing.T) {
	testDevice, err := generateTestDevice()
	require.NoError(t, err)

	testPeers := []wgtypes.Peer{
		{},
		{},
	}

	wgTestDevice := &wgtypes.Device{
		Peers: testPeers,
	}

	testDevice = testDevice.PopulateDynamicFields(wgTestDevice)
	require.Equal(t, testDevice.CurrentPeersCount, len(wgTestDevice.Peers))
}

func TestEntityDevice_IsValid(t *testing.T) {
	testDevice, err := generateTestDevice()
	require.NoError(t, err)

	require.Error(t, testDevice.IsValid(), "device: endpoint should not be empty")

	testDevice.Endpoint = "192.0.2"
	require.Error(t, testDevice.IsValid(), "device: endpoint is not a valid IPv4 or IPv6")

	testDevice.Endpoint = "192.0.2.1"
	require.Error(t, testDevice.IsValid(), "device: address should not be empty")

	testDevice.Address = "10.0.0.1"
	require.Error(t, testDevice.IsValid(), "device: address is not a valid CIDR IP address")

	testDevice.Address = "10.0.0.1/24"
	require.NoError(t, testDevice.IsValid())
}

func TestEntityDevice_GetAvailableIP(t *testing.T) {
	testDevice, err := generateTestDevice()
	require.NoError(t, err)

	testDevice.Endpoint = "192.0.2.1"
	testDevice.Address = "10.0.0.1/24"

	err = testDevice.ComputeMaxPeersCount()
	require.NoError(t, err)

	err = testDevice.ComputeInitialReservedIPs()
	require.NoError(t, err)

	testIP, err := testDevice.GetAvailableIP()
	require.NoError(t, err)
	require.Equal(t, testIP.String(), "10.0.0.0/32")

	testIP, err = testDevice.GetAvailableIP()
	require.NoError(t, err)
	require.Equal(t, testIP.String(), "10.0.0.2/32")

	testDevice.Address = "10.0.0.1/32"
	_, err = testDevice.GetAvailableIP()
	require.ErrorIs(t, err, entity.ErrRunOutOfAddresses)
}
