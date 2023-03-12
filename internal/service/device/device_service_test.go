package deviceservice_test

import (
	"context"
	"testing"

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
