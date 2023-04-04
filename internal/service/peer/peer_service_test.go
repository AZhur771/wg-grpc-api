package peerservice_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	app_mocks "github.com/AZhur771/wg-grpc-api/internal/app/mocks"
	"github.com/AZhur771/wg-grpc-api/internal/dto"
	"github.com/AZhur771/wg-grpc-api/internal/entity"
	peerservice "github.com/AZhur771/wg-grpc-api/internal/service/peer"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func generateMockDevice() (*entity.Device, error) {
	privateKeyMock, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	device := &entity.Device{
		Name:         "wg0",
		Type:         0,
		PrivateKey:   privateKeyMock,
		PublicKey:    privateKeyMock.PublicKey(),
		ListenPort:   51820,
		Endpoint:     "192.0.2.1",
		Address:      "10.0.0.1/24",
		FirewallMark: 0,
	}

	return device, nil
}

func generateMockPeer() (*entity.Peer, error) {
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	presharedKey, err := wgtypes.GenerateKey()
	if err != nil {
		return nil, err
	}

	return &entity.Peer{
		ID:                          uuid.New(),
		Name:                        "Test peer",
		Email:                       "email@example.com",
		PrivateKey:                  entity.WgKey(privateKey),
		PublicKey:                   entity.WgKey(privateKey.PublicKey()),
		PresharedKey:                entity.WgKey(presharedKey),
		PersistentKeepaliveInterval: time.Second * 15,
		LastHandshakeTime:           time.Now(),
		AllowedIPs:                  []string{"10.0.0.3/32"},
		HasPresharedKey:             true,
		IsEnabled:                   true,
		Tags:                        []string{"tag1", "tag2"},
		Description:                 "Test peer description",
	}, nil
}

func TestPeerService_AddPeer(t *testing.T) {
	tests := []struct {
		name                string
		peerName            string
		email               string
		addPresharedKey     bool
		persistentKeepAlive time.Duration
		description         string
		tags                []string
	}{
		{
			name:        "Test peer creation without preshared key and keep alive",
			peerName:    "Test peer 1",
			email:       "email@example.com",
			description: "Test peer description 1",
		},
		{
			name:                "Test peer creation with preshared key and keep alive",
			addPresharedKey:     true,
			peerName:            "Test peer 2",
			persistentKeepAlive: time.Second * 15,
			description:         "Test peer description 2",
			tags:                []string{"tag1", "tag2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockLogger := app_mocks.NewMockLogger(mockCtrl)
			mockDeviceService := app_mocks.NewMockDeviceService(mockCtrl)
			mockStorage := app_mocks.NewMockPeerStorage(mockCtrl)

			mockDevice, err := generateMockDevice()
			require.NoError(t, err)

			mockStorage.EXPECT().
				Add(context.Background(), gomock.Any()).
				DoAndReturn(
					func(ctx context.Context, peer *entity.Peer) (*entity.Peer, error) {
						return peer, nil
					},
				).
				Times(1)

			mockDeviceService.EXPECT().AddPeer(gomock.Any()).Return(nil).Times(1)
			mockDeviceService.EXPECT().GetPeer(gomock.Any()).Return(wgtypes.Peer{}, nil).Times(1)
			mockDeviceService.EXPECT().GetDevice().Return(mockDevice, nil).Times(1)

			srv := peerservice.NewPeerService(
				mockLogger,
				mockDeviceService,
				mockStorage,
			)

			testAddPeerDTO := dto.AddPeerDTO{
				Name:                tt.peerName,
				Email:               tt.email,
				AddPresharedKey:     tt.addPresharedKey,
				PersistentKeepAlive: tt.persistentKeepAlive,
				Tags:                tt.tags,
				Description:         tt.description,
			}

			peer, err := srv.Add(context.Background(), testAddPeerDTO)
			require.NoError(t, err)
			require.Equal(t, tt.peerName, peer.Name)
			require.Equal(t, tt.email, peer.Email)
			require.Equal(t, tt.description, peer.Description)
			require.Equal(t, tt.tags, peer.Tags)
			require.Equal(t, tt.addPresharedKey, peer.HasPresharedKey)
			require.Equal(t, tt.persistentKeepAlive, peer.PersistentKeepaliveInterval)
			require.Equal(t, []string{"10.0.0.0/32"}, peer.AllowedIPs)
		})
	}
}

func TestPeerService_RemovePeer(t *testing.T) {
	tests := []struct {
		name    string
		enabled bool
	}{
		{
			name:    "Test enabled peer deletion",
			enabled: true,
		},
		{
			name:    "Test disabled peer deletion",
			enabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockLogger := app_mocks.NewMockLogger(mockCtrl)
			mockDeviceService := app_mocks.NewMockDeviceService(mockCtrl)
			mockStorage := app_mocks.NewMockPeerStorage(mockCtrl)

			mockPeer, err := generateMockPeer()
			require.NoError(t, err)
			mockPeer.IsEnabled = tt.enabled

			mockStorage.EXPECT().
				Remove(context.Background(), mockPeer.ID).
				Return(mockPeer, nil).
				Times(1)

			expectedRemovePeerCalls := 0
			if tt.enabled {
				expectedRemovePeerCalls = 1
			}

			mockDeviceService.EXPECT().
				RemovePeer(gomock.Any()).
				Return(nil).
				Times(expectedRemovePeerCalls)

			srv := peerservice.NewPeerService(
				mockLogger,
				mockDeviceService,
				mockStorage,
			)

			err = srv.Remove(context.Background(), mockPeer.ID)
			require.NoError(t, err)
		})
	}
}

func TestPeerService_UpdatePeer(t *testing.T) {
	tests := []struct {
		name                string
		peerName            string
		email               string
		addPresharedKey     bool
		removePresharedKey  bool
		persistentKeepAlive time.Duration
		description         string
		tags                []string
		hasPresharedKey     bool
	}{
		{
			name:            "Test updating peer",
			peerName:        "Test peer 1",
			email:           "email@example.com",
			description:     "Test peer description 1",
			hasPresharedKey: true,
		},
		{
			name:               "Test removing preshared key",
			peerName:           "Test peer 2",
			email:              "email@example.com",
			description:        "Test peer description 2",
			removePresharedKey: true,
			hasPresharedKey:    false,
		},
		{
			name:                "Test adding preshared key and keep alive to peer",
			addPresharedKey:     true,
			peerName:            "Test peer 3",
			persistentKeepAlive: time.Second * 15,
			description:         "Test peer description 3",
			tags:                []string{"newtag1", "newtag2"},
			hasPresharedKey:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockLogger := app_mocks.NewMockLogger(mockCtrl)
			mockDeviceService := app_mocks.NewMockDeviceService(mockCtrl)
			mockStorage := app_mocks.NewMockPeerStorage(mockCtrl)

			mockPeer, err := generateMockPeer()
			require.NoError(t, err)

			if tt.addPresharedKey {
				mockPeer.PresharedKey = entity.WgKey(wgtypes.Key{})
				mockPeer.HasPresharedKey = false
			}

			mockStorage.EXPECT().
				Get(context.Background(), mockPeer.ID).
				Return(mockPeer, nil).
				Times(1)
			mockStorage.EXPECT().
				Update(context.Background(), gomock.Any()).
				DoAndReturn(
					func(ctx context.Context, peer *entity.Peer) (*entity.Peer, error) {
						return peer, nil
					},
				).
				Times(1)

			mockDeviceService.EXPECT().UpdatePeer(gomock.Any()).Return(nil).Times(1)
			mockDeviceService.EXPECT().GetPeer(gomock.Any()).Return(wgtypes.Peer{}, nil).Times(1)

			srv := peerservice.NewPeerService(
				mockLogger,
				mockDeviceService,
				mockStorage,
			)

			testUpdatePeerDTO := dto.UpdatePeerDTO{
				ID:                  mockPeer.ID,
				Name:                tt.peerName,
				Email:               tt.email,
				AddPresharedKey:     tt.addPresharedKey,
				RemovePresharedKey:  tt.removePresharedKey,
				PersistentKeepAlive: tt.persistentKeepAlive,
				Tags:                tt.tags,
				Description:         tt.description,
			}

			peer, err := srv.Update(context.Background(), testUpdatePeerDTO)
			require.NoError(t, err)
			require.Equal(t, tt.peerName, peer.Name)
			require.Equal(t, tt.email, peer.Email)
			require.Equal(t, tt.description, peer.Description)
			require.Equal(t, tt.tags, peer.Tags)
			require.Equal(t, tt.persistentKeepAlive, peer.PersistentKeepaliveInterval)
			require.Equal(t, tt.hasPresharedKey, peer.HasPresharedKey)
		})
	}
}

func TestPeerService_GetPeer(t *testing.T) {
	tests := []struct {
		name    string
		enabled bool
	}{
		{
			name:    "Test getting enabled peer",
			enabled: true,
		},
		{
			name:    "Test getting disabled peer",
			enabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockLogger := app_mocks.NewMockLogger(mockCtrl)
			mockDeviceService := app_mocks.NewMockDeviceService(mockCtrl)
			mockStorage := app_mocks.NewMockPeerStorage(mockCtrl)

			mockPeer, err := generateMockPeer()
			require.NoError(t, err)
			mockPeer.IsEnabled = tt.enabled

			mockStorage.EXPECT().
				Get(context.Background(), mockPeer.ID).
				Return(mockPeer, nil).
				Times(1)

			expectedGetPeerCalls := 0
			if tt.enabled {
				expectedGetPeerCalls = 1
			}

			mockDeviceService.EXPECT().
				GetPeer(wgtypes.Key(mockPeer.PublicKey)).
				Return(wgtypes.Peer{}, nil).
				Times(expectedGetPeerCalls)

			srv := peerservice.NewPeerService(
				mockLogger,
				mockDeviceService,
				mockStorage,
			)

			peer, err := srv.Get(context.Background(), mockPeer.ID)
			require.NoError(t, err)
			require.Equal(t, mockPeer.ID, peer.ID)
		})
	}
}

func TestPeerService_GetPeers(t *testing.T) {
	tests := []struct {
		name  string
		skip  int
		limit int

		expectedPeersCount int
		expectedTotal      int
		expectedHasNext    bool
	}{
		{
			name:               "Test getting all peers with default skip and limit",
			expectedPeersCount: 20,
			expectedTotal:      30,
			expectedHasNext:    true,
		},
		{
			name:               "Test getting all peers with skip=5",
			skip:               5,
			expectedPeersCount: 20,
			expectedTotal:      30,
			expectedHasNext:    true,
		},
		{
			name:               "Test getting all peers with skip=10",
			skip:               10,
			expectedPeersCount: 20,
			expectedTotal:      30,
			expectedHasNext:    false,
		},
		{
			name:               "Test getting all peers with skip=5 and limit=10",
			skip:               5,
			limit:              10,
			expectedPeersCount: 10,
			expectedTotal:      30,
			expectedHasNext:    true,
		},
		{
			name:               "Test getting all peers with limit=1",
			limit:              1,
			expectedPeersCount: 1,
			expectedTotal:      30,
			expectedHasNext:    true,
		},
		{
			name:               "Test getting all peers with limit=30",
			limit:              30,
			expectedPeersCount: 30,
			expectedTotal:      30,
			expectedHasNext:    false,
		},
		{
			name:               "Test getting all peers with limit=40",
			limit:              40,
			expectedPeersCount: 30,
			expectedTotal:      30,
			expectedHasNext:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockLogger := app_mocks.NewMockLogger(mockCtrl)
			mockDeviceService := app_mocks.NewMockDeviceService(mockCtrl)
			mockStorage := app_mocks.NewMockPeerStorage(mockCtrl)

			mockPeers := make([]*entity.Peer, 0, 30)
			mockWgPeers := make([]wgtypes.Peer, 0, 30)

			n := 0
			for n < tt.expectedTotal {
				mockPeer, err := generateMockPeer()
				require.NoError(t, err)

				mockPeer.Name = fmt.Sprintf("Test peer %d", n)
				mockPeers = append(mockPeers, mockPeer)

				mockWgPeer := wgtypes.Peer{
					PublicKey: wgtypes.Key(mockPeer.PublicKey),
				}
				mockWgPeers = append(mockWgPeers, mockWgPeer)

				n++
			}

			mockStorage.EXPECT().
				CountAll(context.Background()).
				Return(len(mockPeers), nil).
				Times(1)

			mockStorage.EXPECT().
				GetAll(
					context.Background(),
					gomock.AssignableToTypeOf(tt.skip),
					gomock.AssignableToTypeOf(tt.limit),
				).
				DoAndReturn(
					func(ctx context.Context, skip, limit int) ([]*entity.Peer, error) {
						if skip >= len(mockPeers) {
							return make([]*entity.Peer, 0), nil
						}

						if limit == 0 || skip+limit >= len(mockPeers) {
							return mockPeers[skip:], nil
						}

						return mockPeers[skip : skip+limit], nil
					},
				).
				Times(1)

			mockDeviceService.EXPECT().
				GetPeers().
				Return(mockWgPeers, nil).
				Times(1)

			srv := peerservice.NewPeerService(
				mockLogger,
				mockDeviceService,
				mockStorage,
			)

			getPeersRequestDTO := dto.GetPeersRequestDTO{
				Skip:  tt.skip,
				Limit: tt.limit,
			}
			getPeersResponseDTO, err := srv.GetAll(context.Background(), getPeersRequestDTO)
			require.NoError(t, err)
			require.Equal(t, tt.expectedPeersCount, len(getPeersResponseDTO.Peers))
			require.Equal(t, tt.expectedTotal, getPeersResponseDTO.Total)
			require.Equal(t, tt.expectedHasNext, getPeersResponseDTO.HasNext)
		})
	}
}
