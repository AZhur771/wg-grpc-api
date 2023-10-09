package handlers

import (
	wgpb "github.com/AZhur771/wg-grpc-api/gen"
	"github.com/AZhur771/wg-grpc-api/internal/entity"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapEntityDeviceToPbDeivce(dev *entity.Device) *wgpb.Device {
	return &wgpb.Device{
		Id:                  dev.ID.String(),
		Name:                dev.Name,
		Description:         dev.Description,
		Type:                dev.Type.String(),
		PublicKey:           dev.PublicKey.String(),
		FirewallMark:        int32(dev.FirewallMark),
		MaxPeersCount:       int32(dev.MaxPeersCount),
		CurrentPeersCount:   int32(dev.CurrentPeersCount),
		Endpoint:            dev.Endpoint,
		Address:             dev.Address,
		Mtu:                 int32(dev.MTU),
		Dns:                 dev.DNS,
		Table:               dev.Table,
		PersistentKeepAlive: int32(dev.PersistentKeepAlive),
		PreUp:               dev.PreUp,
		PreDown:             dev.PreDown,
		PostUp:              dev.PostDown,
		PostDown:            dev.PostDown,
	}
}

func mapEntityPeerToPbPeer(peer *entity.Peer) *wgpb.Peer {
	return &wgpb.Peer{
		Id:                  peer.ID.String(),
		DeviceId:            peer.DeviceID.String(),
		Name:                peer.Name,
		Email:               peer.Email,
		PublicKey:           peer.PublicKey.String(),
		Endpoint:            peer.Endpoint.String(),
		PersistentKeepAlive: int32(peer.PersistentKeepaliveInterval.Seconds()),
		AllowedIps:          peer.AllowedIPs,
		ProtocolVersion:     uint32(peer.ProtocolVersion),
		ReceiveBytes:        peer.ReceiveBytes,
		TransmitBytes:       peer.TransmitBytes,
		LastHandshake:       timestamppb.New(peer.LastHandshakeTime),
		HasPresharedKey:     peer.HasPresharedKey,
		IsEnabled:           peer.IsEnabled,
		IsActive:            peer.IsActive,
		Description:         peer.Description,
		Dns:                 peer.DNS,
		Mtu:                 int32(peer.MTU),
	}
}

func mapEntityPeerToPbPeerAbridged(peer *entity.Peer) *wgpb.PeerAbridged {
	return &wgpb.PeerAbridged{
		Id:                  peer.ID.String(),
		DeviceId:            peer.DeviceID.String(),
		Name:                peer.Name,
		Email:               peer.Email,
		PublicKey:           peer.PublicKey.String(),
		PersistentKeepAlive: int32(peer.PersistentKeepaliveInterval.Seconds()),
		AllowedIps:          peer.AllowedIPs,
		HasPresharedKey:     peer.HasPresharedKey,
		IsEnabled:           peer.IsEnabled,
		Description:         peer.Description,
		Dns:                 peer.DNS,
		Mtu:                 int32(peer.MTU),
	}
}

func mapNames(s string) string {
	switch s {
	case "id":
		return "ID"
	case "name":
		return "Name"
	case "email":
		return "Email"
	case "description":
		return "Description"
	case "endpoint":
		return "Endpoint"
	case "address":
		return "Address"
	case "table":
		return "Table"
	case "firewall_mark":
		return "FirewallMark"
	case "dns":
		return "DNS"
	case "mtu":
		return "MTU"
	case "persistent_keep_alive":
		return "PersistentKeepAlive"
	case "add_preshared_key":
		return "AddPresharedKey"
	case "remove_preshared_key":
		return "RemovePresharedKey"
	case "pre_up":
		return "PreUp"
	case "pre_down":
		return "PreDown"
	case "post_up":
		return "PostUp"
	case "post_down":
		return "PostDown"
	default:
		return ""
	}
}
