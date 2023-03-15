package peerstorage_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	app_mocks "github.com/AZhur771/wg-grpc-api/internal/app/mocks"
	"github.com/AZhur771/wg-grpc-api/internal/entity"
	peerstorage "github.com/AZhur771/wg-grpc-api/internal/storage/peer"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
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
		PrivateKey:                  privateKey,
		PublicKey:                   privateKey.PublicKey(),
		PresharedKey:                presharedKey,
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
}
