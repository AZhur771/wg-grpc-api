package peer_service_test

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	app_mocks "github.com/AZhur771/wg-grpc-api/internal/app/mocks"
// 	"github.com/AZhur771/wg-grpc-api/internal/entity"
// 	"github.com/AZhur771/wg-grpc-api/internal/service"
// 	"github.com/golang/mock/gomock"
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/require"
// 	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
// )

// var (
// 	uuidMock             = "acda9b63-45ae-4352-995c-82202086cac4"
// 	devicePrivateKeyMock = "WDhbZ4+4sE8LmIu4tSA1AXINX1ly+d+ZUwzazdiRMFU="
// 	devicePublicKeyMock  = "MG+IiZS7uPsLMigRQoMch5MD7H2XCqEM+o9QJf1VGD4="
// 	deviceNameMock       = "wg0"
// 	deviceAddressMock    = "10.0.0.1/24" // address in CIDR notation
// 	deviceEndpointMock   = "50.75.40.3"
// 	peerFolderMock       = "/etc/some_folder"
// )

// func generateMockDevice(peers []wgtypes.Peer) (*wgtypes.Device, error) {
// 	privateKeyMock, err := wgtypes.ParseKey(devicePrivateKeyMock)
// 	if err != nil {
// 		return nil, err
// 	}

// 	publicKeyMock, err := wgtypes.ParseKey(devicePublicKeyMock)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &wgtypes.Device{
// 		Name:         "wg0",
// 		Type:         0,
// 		PrivateKey:   privateKeyMock,
// 		PublicKey:    publicKeyMock,
// 		ListenPort:   51820,
// 		FirewallMark: 0,
// 		Peers:        peers,
// 	}, nil
// }

// func generateMockPersistedPeer() (*entity.PersistedPeer, error) {
// 	privateKey, err := wgtypes.GeneratePrivateKey()
// 	if err != nil {
// 		return nil, err
// 	}

// 	presharedKey, err := wgtypes.GenerateKey()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &entity.PersistedPeer{
// 		ID:                          uuidMock,
// 		PrivateKey:                  privateKey.String(),
// 		PublicKey:                   privateKey.PublicKey().String(),
// 		PresharedKey:                presharedKey.String(),
// 		PersistentKeepaliveInterval: time.Second * 15,
// 		LastHandshakeTime:           time.Now(),
// 		AllowedIPs:                  []string{"10.0.0.3/32"},
// 		HasPresharedKey:             true,
// 		IsEnabled:                   true,
// 		Description:                 "some description",
// 	}, nil
// }

// func TestPeerService_AddPeer(t *testing.T) {
// 	tests := []struct {
// 		name                string
// 		addPresharedKey     bool
// 		persistentKeepAlive time.Duration
// 		description         string
// 	}{
// 		{
// 			name:        "Test peer creation without preshared key and keep alive",
// 			description: "Test peer 1",
// 		},
// 		{
// 			name:                "Test peer creation with preshared key and keep alive",
// 			addPresharedKey:     true,
// 			persistentKeepAlive: time.Second * 15,
// 			description:         "Test peer 1",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockCtrl := gomock.NewController(t)
// 			defer mockCtrl.Finish()

// 			mockLogger := app_mocks.NewMockLogger(mockCtrl)
// 			mockWgCtrl := app_mocks.NewMockWgCtrl(mockCtrl)
// 			mockStorage := app_mocks.NewMockPeerStorage(mockCtrl)

// 			idMock, err := uuid.Parse(uuidMock)
// 			require.NoError(t, err)

// 			mockDevice, err := generateMockDevice([]wgtypes.Peer{})
// 			require.NoError(t, err)

// 			mockWgCtrl.EXPECT().GetDevice().Return(mockDevice, nil).Times(1)

// 			mockWgCtrl.EXPECT().AddPeer(gomock.Any()).Return(nil).Times(1)
// 			mockWgCtrl.EXPECT().GetPeers().Return([]wgtypes.Peer{}, nil).Times(1)

// 			mockStorage.EXPECT().Add(context.Background(), gomock.Any()).Return(idMock, nil).Times(1)
// 			mockStorage.EXPECT().GetAll(context.Background()).Return([]*entity.PersistedPeer{}, nil).Times(1)

// 			srv := service.NewPeerService(
// 				mockLogger,
// 				mockWgCtrl,
// 				mockStorage,
// 			)
// 			srv.Setup(
// 				context.Background(),
// 				deviceNameMock,
// 				deviceAddressMock,
// 				deviceEndpointMock,
// 				peerFolderMock,
// 			)

// 			peer, err := srv.Add(context.Background(), tt.addPresharedKey, tt.persistentKeepAlive, tt.description)
// 			require.NoError(t, err)
// 			require.Equal(t, uuidMock, peer.ID.String())
// 			require.Equal(t, tt.addPresharedKey, peer.HasPresharedKey)
// 			require.Equal(t, tt.persistentKeepAlive, peer.PersistentKeepaliveInterval)
// 			require.Equal(t, tt.description, peer.Description)
// 		})
// 	}
// }

// func TestPeerService_RemovePeer(t *testing.T) {
// 	tests := []struct {
// 		name                string
// 		addPresharedKey     bool
// 		persistentKeepAlive time.Duration
// 		description         string
// 	}{
// 		{
// 			name:        "Test peer creation without preshared key and keep alive",
// 			description: "Test peer 1",
// 		},
// 		{
// 			name:                "Test peer creation with preshared key and keep alive",
// 			addPresharedKey:     true,
// 			persistentKeepAlive: time.Second * 15,
// 			description:         "Test peer 1",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockCtrl := gomock.NewController(t)
// 			defer mockCtrl.Finish()

// 			mockLogger := app_mocks.NewMockLogger(mockCtrl)
// 			mockWgCtrl := app_mocks.NewMockWgCtrl(mockCtrl)
// 			mockStorage := app_mocks.NewMockPeerStorage(mockCtrl)

// 			idMock, err := uuid.Parse(uuidMock)
// 			require.NoError(t, err)

// 			mockDevice, err := generateMockDevice([]wgtypes.Peer{})
// 			require.NoError(t, err)

// 			mockWgCtrl.EXPECT().GetDevice().Return(mockDevice, nil).Times(1)

// 			mockWgCtrl.EXPECT().RemovePeer(gomock.Any()).Return(nil).Times(1)
// 			mockWgCtrl.EXPECT().GetPeers().Return([]wgtypes.Peer{}, nil).Times(1)

// 			persistedPeerMock, err := generateMockPersistedPeer()
// 			require.NoError(t, err)

// 			mockStorage.EXPECT().Get(context.Background(), idMock).Return(persistedPeerMock, nil).Times(1)
// 			mockStorage.EXPECT().GetAll(context.Background()).Return([]*entity.PersistedPeer{}, nil).Times(1)

// 			srv := service.NewPeerService(
// 				mockLogger,
// 				mockWgCtrl,
// 				mockStorage,
// 			)
// 			srv.Setup(
// 				context.Background(),
// 				deviceNameMock,
// 				deviceAddressMock,
// 				deviceEndpointMock,
// 				peerFolderMock,
// 			)

// 			peer, err := srv.Remove(context.Background(), idMock)
// 			require.NoError(t, err)
// 			require.Equal(t, uuidMock, peer.ID.String())
// 			require.Equal(t, persistedPeerMock.HasPresharedKey, peer.PresharedKey)
// 			require.Equal(t, persistedPeerMock.PersistentKeepaliveInterval, persistedPeerMock.PersistentKeepaliveInterval)
// 			require.Equal(t, persistedPeerMock.Description, persistedPeerMock.Description)
// 		})
// 	}
// }

// func TestPeerService_UpdatePeer(t *testing.T) {

// }

// func TestPeerService_GetPeer(t *testing.T) {

// }

// func TestPeerService_GetPeers(t *testing.T) {

// }
