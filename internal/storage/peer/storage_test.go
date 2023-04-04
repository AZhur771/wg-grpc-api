package peerstorage_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	app_mocks "github.com/AZhur771/wg-grpc-api/internal/app/mocks"
	"github.com/AZhur771/wg-grpc-api/internal/entity"
	peerstorage "github.com/AZhur771/wg-grpc-api/internal/storage/peer"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/require"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

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

func TestPeerStorage_Add(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLogger := app_mocks.NewMockLogger(mockCtrl)

	testDir := t.TempDir()
	storage := peerstorage.NewPeerStorage(mockLogger, testDir)

	mockPeer, err := generateMockPeer()
	require.NoError(t, err)

	peer, err := storage.Add(context.Background(), mockPeer)
	require.NoError(t, err)
	require.Equal(t, mockPeer.Name, peer.Name)

	files, err := ioutil.ReadDir(testDir)
	require.NoError(t, err)
	require.Equal(t, 1, len(files))
	require.Equal(t, fmt.Sprintf("%s.json", peer.ID), files[0].Name())

	peerEntity := &entity.Peer{}

	bytes, err := ioutil.ReadFile(path.Join(testDir, files[0].Name()))
	require.NoError(t, err)

	err = easyjson.Unmarshal(bytes, peerEntity)
	require.NoError(t, err)
	require.Equal(t, peerEntity.ID, peer.ID)
	require.Equal(t, peerEntity.Name, peer.Name)
}

func TestPeerStorage_Update(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLogger := app_mocks.NewMockLogger(mockCtrl)

	testDir := t.TempDir()
	storage := peerstorage.NewPeerStorage(mockLogger, testDir)

	mockPeer, err := generateMockPeer()
	require.NoError(t, err)

	bytes, err := easyjson.Marshal(mockPeer)
	require.NoError(t, err)

	file, err := os.Create(path.Join(testDir, fmt.Sprintf("%s.json", mockPeer.ID)))
	require.NoError(t, err)

	_, err = file.Write(bytes)
	require.NoError(t, err)

	mockPeer.Name = "New name"

	peer, err := storage.Update(context.Background(), mockPeer)
	require.NoError(t, err)
	require.Equal(t, mockPeer.Name, peer.Name)
}

func TestPeerStorage_Remove(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLogger := app_mocks.NewMockLogger(mockCtrl)

	testDir := t.TempDir()
	storage := peerstorage.NewPeerStorage(mockLogger, testDir)

	mockPeer, err := generateMockPeer()
	require.NoError(t, err)

	bytes, err := easyjson.Marshal(mockPeer)
	require.NoError(t, err)

	file, err := os.Create(path.Join(testDir, fmt.Sprintf("%s.json", mockPeer.ID)))
	require.NoError(t, err)

	_, err = file.Write(bytes)
	require.NoError(t, err)

	peer, err := storage.Remove(context.Background(), mockPeer.ID)
	require.NoError(t, err)
	require.Equal(t, mockPeer.Name, peer.Name)

	_, err = os.Stat(path.Join(testDir, fmt.Sprintf("%s.json", mockPeer.ID)))
	require.ErrorIs(t, err, os.ErrNotExist)
}

func TestPeerStorage_Get(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLogger := app_mocks.NewMockLogger(mockCtrl)

	testDir := t.TempDir()
	storage := peerstorage.NewPeerStorage(mockLogger, testDir)

	mockPeer, err := generateMockPeer()
	require.NoError(t, err)

	bytes, err := easyjson.Marshal(mockPeer)
	require.NoError(t, err)

	file, err := os.Create(path.Join(testDir, fmt.Sprintf("%s.json", mockPeer.ID)))
	require.NoError(t, err)

	_, err = file.Write(bytes)
	require.NoError(t, err)

	peer, err := storage.Get(context.Background(), mockPeer.ID)
	require.NoError(t, err)
	require.Equal(t, mockPeer.Name, peer.Name)
}

func TestPeerStorage_GetAll(t *testing.T) {
	tests := []struct {
		name  string
		skip  int
		limit int

		expectedPeersCount int
	}{
		{
			name:               "Test getting all peers with default skip and limit",
			expectedPeersCount: 30,
		},
		{
			name:               "Test getting all peers with skip=5",
			skip:               5,
			expectedPeersCount: 25,
		},
		{
			name:               "Test getting all peers with skip=10",
			skip:               10,
			expectedPeersCount: 20,
		},
		{
			name:               "Test getting all peers with skip=5 and limit=10",
			skip:               5,
			limit:              10,
			expectedPeersCount: 10,
		},
		{
			name:               "Test getting all peers with limit=1",
			limit:              1,
			expectedPeersCount: 1,
		},
		{
			name:               "Test getting all peers with limit=30",
			limit:              30,
			expectedPeersCount: 30,
		},
		{
			name:               "Test getting all peers with limit=40",
			limit:              40,
			expectedPeersCount: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockLogger := app_mocks.NewMockLogger(mockCtrl)

			testDir := t.TempDir()
			storage := peerstorage.NewPeerStorage(mockLogger, testDir)

			n := 0
			for n < 30 {
				mockPeer, err := generateMockPeer()
				require.NoError(t, err)
				mockPeer.Name = fmt.Sprintf("Test peer %d", n)

				bytes, err := easyjson.Marshal(mockPeer)
				require.NoError(t, err)

				file, err := os.Create(path.Join(testDir, fmt.Sprintf("%s.json", mockPeer.ID)))
				require.NoError(t, err)

				_, err = file.Write(bytes)
				require.NoError(t, err)

				n++
			}

			peers, err := storage.GetAll(context.Background(), tt.skip, tt.limit)
			require.NoError(t, err)
			require.Equal(t, tt.expectedPeersCount, len(peers))
		})
	}
}
