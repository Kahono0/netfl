package putils

import (
	"context"
	"fmt"

	"github.com/kahono0/netfl/pkg/peers"

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

func CreateStrean(host host.Host, peer peer.AddrInfo, protocolID string) (*network.Stream, error) {
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

func CreateStreamFromUnknow(host host.Host, peerID peer.ID, protocolID string) (*network.Stream, error) {
	peer := peers.Store.GetPeerByID(peerID.String())
	if peer == nil {
		return nil, fmt.Errorf("no peer found with ID %s", peerID)
	}

	return CreateStrean(host, *peer, protocolID)

}
