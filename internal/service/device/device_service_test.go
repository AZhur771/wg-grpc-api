package deviceservice_test

import (
	"context"
	"net"
	"testing"
	"time"

	app_mocks "github.com/AZhur771/wg-grpc-api/internal/app/mocks"
	"github.com/AZhur771/wg-grpc-api/internal/entity"
	deviceservice "github.com/AZhur771/wg-grpc-api/internal/service/device"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func generateMockWgDevice() (*wgtypes.Device, error) {
	privateKeyMock, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	device := &wgtypes.Device{
		Name:         "wg0",
		Type:         0,
		PrivateKey:   privateKeyMock,
		PublicKey:    privateKeyMock.PublicKey(),
		ListenPort:   51820,
		FirewallMark: 0,
	}

	return device, nil
}

func generateMockWgPeer() (wgtypes.Peer, error) {
	wgPeer := wgtypes.Peer{}

	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return wgPeer, err
	}

	wgPeer.PublicKey = privateKey.PublicKey()

	wgPeer.PresharedKey, err = wgtypes.GenerateKey()
	if err != nil {
		return wgPeer, err
	}

	_, mockIPNet, _ := net.ParseCIDR("10.0.0.3/32")

	wgPeer.PersistentKeepaliveInterval = time.Second * 15
	wgPeer.LastHandshakeTime = time.Now()
	wgPeer.AllowedIPs = []net.IPNet{*mockIPNet}

	return wgPeer, nil
}

func generateMockWgPeerConfig() (wgtypes.PeerConfig, error) {
	wgPeerConf := wgtypes.PeerConfig{}

	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return wgPeerConf, err
	}

	wgPeerConf.PublicKey = privateKey.PublicKey()

	presharedKey, err := wgtypes.GenerateKey()
	if err != nil {
		return wgPeerConf, err
	}

	wgPeerConf.PresharedKey = &presharedKey

	_, mockIPNet, _ := net.ParseCIDR("10.0.0.3/32")

	persistentKeepaliveInterval := time.Second * 15
	wgPeerConf.PersistentKeepaliveInterval = &persistentKeepaliveInterval
	wgPeerConf.AllowedIPs = []net.IPNet{*mockIPNet}

	return wgPeerConf, nil
}

func TestDeviceService_GetDevice(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLogger := app_mocks.NewMockLogger(mockCtrl)
	mockWgCtrl := app_mocks.NewMockWgCtrl(mockCtrl)
	mockStorage := app_mocks.NewMockPeerStorage(mockCtrl)

	mockWgDevice, err := generateMockWgDevice()
	require.NoError(t, err)

	mockWgCtrl.EXPECT().
		Device(mockWgDevice.Name).
		Return(mockWgDevice, nil).
		AnyTimes()

	mockStorage.EXPECT().
		GetAll(context.Background(), 0, 0).
		Return(make([]*entity.Peer, 0), nil).
		AnyTimes()

	srv := deviceservice.NewDeviceService(
		mockLogger,
		mockWgCtrl,
		mockStorage,
	)

	err = srv.Setup(context.Background(), "wg0", "192.0.2.1", "10.0.0.1/24")
	require.NoError(t, err)

	device, err := srv.GetDevice()
	require.NoError(t, err)
	require.Equal(t, device.Name, "wg0")
	require.Equal(t, device.Address, "10.0.0.1/24")
	require.Equal(t, device.Endpoint, "192.0.2.1")
	require.Equal(t, device.MaxPeersCount, 255)
}

func TestDeviceService_GetPeer(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLogger := app_mocks.NewMockLogger(mockCtrl)
	mockWgCtrl := app_mocks.NewMockWgCtrl(mockCtrl)
	mockStorage := app_mocks.NewMockPeerStorage(mockCtrl)

	mockWgDevice, err := generateMockWgDevice()
	require.NoError(t, err)

	mockWgPeer, err := generateMockWgPeer()
	require.NoError(t, err)

	mockWgDevice.Peers = []wgtypes.Peer{mockWgPeer}

	mockWgCtrl.EXPECT().
		Device(mockWgDevice.Name).
		Return(mockWgDevice, nil).
		AnyTimes()

	mockStorage.EXPECT().
		GetAll(context.Background(), 0, 0).
		Return(make([]*entity.Peer, 0), nil).
		AnyTimes()

	srv := deviceservice.NewDeviceService(
		mockLogger,
		mockWgCtrl,
		mockStorage,
	)

	err = srv.Setup(context.Background(), "wg0", "192.0.2.1", "10.0.0.1/24")
	require.NoError(t, err)

	peer, err := srv.GetPeer(mockWgPeer.PublicKey)
	require.NoError(t, err)
	require.Equal(t, mockWgPeer.PublicKey, peer.PublicKey)
}

func TestDeviceService_GetPeers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLogger := app_mocks.NewMockLogger(mockCtrl)
	mockWgCtrl := app_mocks.NewMockWgCtrl(mockCtrl)
	mockStorage := app_mocks.NewMockPeerStorage(mockCtrl)

	mockWgDevice, err := generateMockWgDevice()
	require.NoError(t, err)

	mockWgPeer, err := generateMockWgPeer()
	require.NoError(t, err)

	mockWgDevice.Peers = []wgtypes.Peer{mockWgPeer}

	mockWgCtrl.EXPECT().
		Device(mockWgDevice.Name).
		Return(mockWgDevice, nil).
		AnyTimes()

	mockStorage.EXPECT().
		GetAll(context.Background(), 0, 0).
		Return(make([]*entity.Peer, 0), nil).
		AnyTimes()

	srv := deviceservice.NewDeviceService(
		mockLogger,
		mockWgCtrl,
		mockStorage,
	)

	err = srv.Setup(context.Background(), "wg0", "192.0.2.1", "10.0.0.1/24")
	require.NoError(t, err)

	peers, err := srv.GetPeers()
	require.NoError(t, err)
	require.Equal(t, len(mockWgDevice.Peers), len(peers))
}

func TestDeviceService_AddPeer(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLogger := app_mocks.NewMockLogger(mockCtrl)
	mockWgCtrl := app_mocks.NewMockWgCtrl(mockCtrl)
	mockStorage := app_mocks.NewMockPeerStorage(mockCtrl)

	mockWgDevice, err := generateMockWgDevice()
	require.NoError(t, err)

	mockWgPeerConf, err := generateMockWgPeerConfig()
	require.NoError(t, err)

	mockWgCtrl.EXPECT().
		Device(mockWgDevice.Name).
		Return(mockWgDevice, nil).
		AnyTimes()

	mockStorage.EXPECT().
		GetAll(context.Background(), 0, 0).
		Return(make([]*entity.Peer, 0), nil).
		AnyTimes()

	mockWgCtrl.EXPECT().
		ConfigureDevice(mockWgDevice.Name, gomock.Any()).
		Return(nil).
		Times(1)

	srv := deviceservice.NewDeviceService(
		mockLogger,
		mockWgCtrl,
		mockStorage,
	)

	err = srv.Setup(context.Background(), "wg0", "192.0.2.1", "10.0.0.1/24")
	require.NoError(t, err)

	err = srv.AddPeer(mockWgPeerConf)
	require.NoError(t, err)
}

func TestDeviceService_UpdatePeer(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLogger := app_mocks.NewMockLogger(mockCtrl)
	mockWgCtrl := app_mocks.NewMockWgCtrl(mockCtrl)
	mockStorage := app_mocks.NewMockPeerStorage(mockCtrl)

	mockWgDevice, err := generateMockWgDevice()
	require.NoError(t, err)

	mockWgPeerConf, err := generateMockWgPeerConfig()
	require.NoError(t, err)

	mockWgCtrl.EXPECT().
		Device(mockWgDevice.Name).
		Return(mockWgDevice, nil).
		AnyTimes()

	mockStorage.EXPECT().
		GetAll(context.Background(), 0, 0).
		Return(make([]*entity.Peer, 0), nil).
		AnyTimes()

	var wgConf wgtypes.Config

	mockWgCtrl.EXPECT().
		ConfigureDevice(mockWgDevice.Name, gomock.AssignableToTypeOf(wgConf)).
		DoAndReturn(func(_ string, cfg wgtypes.Config) error {
			wgConf = cfg
			return nil
		}).
		Times(1)

	srv := deviceservice.NewDeviceService(
		mockLogger,
		mockWgCtrl,
		mockStorage,
	)

	err = srv.Setup(context.Background(), "wg0", "192.0.2.1", "10.0.0.1/24")
	require.NoError(t, err)

	err = srv.UpdatePeer(mockWgPeerConf)
	require.NoError(t, err)
	require.True(t, wgConf.Peers[0].UpdateOnly)
}

func TestDeviceService_RemovePeer(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLogger := app_mocks.NewMockLogger(mockCtrl)
	mockWgCtrl := app_mocks.NewMockWgCtrl(mockCtrl)
	mockStorage := app_mocks.NewMockPeerStorage(mockCtrl)

	mockWgDevice, err := generateMockWgDevice()
	require.NoError(t, err)

	mockWgPeerConf, err := generateMockWgPeerConfig()
	require.NoError(t, err)

	mockWgCtrl.EXPECT().
		Device(mockWgDevice.Name).
		Return(mockWgDevice, nil).
		AnyTimes()

	mockStorage.EXPECT().
		GetAll(context.Background(), 0, 0).
		Return(make([]*entity.Peer, 0), nil).
		AnyTimes()

	var wgConf wgtypes.Config

	mockWgCtrl.EXPECT().
		ConfigureDevice(mockWgDevice.Name, gomock.AssignableToTypeOf(wgConf)).
		DoAndReturn(func(_ string, cfg wgtypes.Config) error {
			wgConf = cfg
			return nil
		}).
		Times(1)

	srv := deviceservice.NewDeviceService(
		mockLogger,
		mockWgCtrl,
		mockStorage,
	)

	err = srv.Setup(context.Background(), "wg0", "192.0.2.1", "10.0.0.1/24")
	require.NoError(t, err)

	err = srv.RemovePeer(mockWgPeerConf)
	require.NoError(t, err)
	require.True(t, wgConf.Peers[0].Remove)
}
