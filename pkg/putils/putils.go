package putils

import (
	"context"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

func Connect(host host.Host, peer peer.AddrInfo) error {
	if host.Network().Connectedness(peer.ID) != network.Connected {
		return host.Connect(context.Background(), peer)
	}

	return nil
}

func CreateStream(host host.Host, peer peer.AddrInfo, protocolID string) (*network.Stream, error) {
	ctx := context.Background()
	if err := Connect(host, peer); err != nil {
		return nil, err
	}
	s, err := host.NewStream(ctx, peer.ID, protocol.ID(protocolID))
	if err != nil {
		return nil, err
	}

	return &s, nil
}
