package entity_test

import (
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

	allowedIPs := make([]net.IPNet, 0, len(allowedIPsMock))
	for _, addr := range allowedIPsMock {
		_, ipnet, err := net.ParseCIDR(addr)
		if err != nil {
			return nil, err
		}

		allowedIPs = append(allowedIPs, *ipnet)
	}

	return &entity.Peer{
		ID:                          id,
		PrivateKey:                  privateKey,
		PublicKey:                   publicKey,
		PresharedKey:                presharedKey,
		PersistentKeepaliveInterval: time.Second * 15,
		AllowedIPs:                  allowedIPs,
		HasPresharedKey:             true,
		IsEnabled:                   true,
		Description:                 "some description",
	}, nil
}

func TestPeerToPersistedPeerConversion(t *testing.T) {
	testPeer, err := generateTestPeer()
	require.NoError(t, err)

	persistedPeer := testPeer.ToPersistedPeer()
	require.Equal(t, uuidMock, persistedPeer.ID)
	require.Equal(t, privateKeyMock, persistedPeer.PrivateKey)
	require.Equal(t, publicKeyMock, persistedPeer.PublicKey)
	require.Equal(t, presharedKeyMock, persistedPeer.PresharedKey)
	require.Equal(t, testPeer.PersistentKeepaliveInterval, persistedPeer.PersistentKeepaliveInterval)
	require.Equal(t, allowedIPsMock, persistedPeer.AllowedIPs)
	require.Equal(t, testPeer.HasPresharedKey, persistedPeer.HasPresharedKey)
	require.Equal(t, testPeer.IsEnabled, persistedPeer.IsEnabled)
	require.Equal(t, testPeer.Description, persistedPeer.Description)
}

func TestPeerPopulateDynamicFields(t *testing.T) {
	testPeer, err := generateTestPeer()
	require.NoError(t, err)

	endpoint, err := net.ResolveUDPAddr("udp", "203.0.113.0:51823")
	require.NoError(t, err)

	wgTestPeer := &wgtypes.Peer{
		PublicKey:                   testPeer.PublicKey,
		PresharedKey:                testPeer.PresharedKey,
		Endpoint:                    endpoint,
		PersistentKeepaliveInterval: time.Second * 15,
		LastHandshakeTime:           lastHandshakeTimeMock,
		ReceiveBytes:                111,
		TransmitBytes:               222,
		AllowedIPs:                  testPeer.AllowedIPs,
		ProtocolVersion:             1,
	}

	testPeer = testPeer.PopulateDynamicFields(wgTestPeer)

	require.Equal(t, testPeer.Endpoint, wgTestPeer.Endpoint)
	require.Equal(t, testPeer.LastHandshakeTime, wgTestPeer.LastHandshakeTime)
	require.Equal(t, testPeer.ReceiveBytes, wgTestPeer.ReceiveBytes)
	require.Equal(t, testPeer.TransmitBytes, wgTestPeer.TransmitBytes)
	require.Equal(t, testPeer.ProtocolVersion, wgTestPeer.ProtocolVersion)
}

func TestPersistedPeerToPeerConversion(t *testing.T) {
	testPersistedPeer := &entity.PersistedPeer{
		ID:                          uuidMock,
		PrivateKey:                  privateKeyMock,
		PublicKey:                   publicKeyMock,
		PresharedKey:                presharedKeyMock,
		PersistentKeepaliveInterval: time.Second * 15,
		LastHandshakeTime:           lastHandshakeTimeMock,
		AllowedIPs:                  allowedIPsMock,
		HasPresharedKey:             true,
		IsEnabled:                   true,
		Description:                 "some test peer",
	}

	peer, err := testPersistedPeer.ToPeer()
	require.NoError(t, err)

	require.Equal(t, testPersistedPeer.ID, peer.ID.String())
	require.Equal(t, testPersistedPeer.PrivateKey, peer.PrivateKey.String())
	require.Equal(t, testPersistedPeer.PublicKey, peer.PublicKey.String())
	require.Equal(t, testPersistedPeer.PresharedKey, peer.PresharedKey.String())
	require.Equal(t, testPersistedPeer.PersistentKeepaliveInterval, peer.PersistentKeepaliveInterval)
	require.Equal(t, testPersistedPeer.HasPresharedKey, peer.HasPresharedKey)
	require.Equal(t, testPersistedPeer.IsEnabled, peer.IsEnabled)
	require.Equal(t, testPersistedPeer.Description, peer.Description)
}
