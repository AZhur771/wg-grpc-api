package grpcserver

import (
	"net"

	peerpb "github.com/AZhur771/wg-grpc-api/gen"
	"github.com/AZhur771/wg-grpc-api/internal/entity"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertAllowedIpsToStringSlice(data []net.IPNet) []string {
	res := make([]string, 0, len(data))

	for _, e := range data {
		res = append(res, e.String())
	}

	return res
}

func mapEntityPeerToPbPeer(peer *entity.Peer) *peerpb.Peer {
	return &peerpb.Peer{
		Id:                  peer.ID.String(),
		PublicKey:           peer.PublicKey.String(),
		Endpoint:            peer.Endpoint.String(),
		PersistentKeepAlive: durationpb.New(peer.PersistentKeepaliveInterval),
		AllowedIps:          convertAllowedIpsToStringSlice(peer.AllowedIPs),
		ProtocolVersion:     uint32(peer.ProtocolVersion),
		ReceiveBytes:        peer.ReceiveBytes,
		TransmitBytes:       peer.TransmitBytes,
		LastHandshake:       timestamppb.New(peer.LastHandshakeTime),
		HasPresharedKey:     peer.HasPresharedKey,
		IsEnabled:           peer.IsEnabled,
		IsActive:            peer.IsActive,
		Description:         peer.Description,
	}
}
