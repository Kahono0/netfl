package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"

	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type discoveryNotifee struct {
	PeerChan chan peer.AddrInfo
}

// interface to be called when new  peer is found
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	n.PeerChan <- pi
}

// Initialize the MDNS service
func initMDNS(peerhost host.Host, rendezvous string) chan peer.AddrInfo {
	// register with service so that we get notified about peer discovery
	n := &discoveryNotifee{}
	n.PeerChan = make(chan peer.AddrInfo)

	// An hour might be a long long period in practical applications. But this is fine for us
	ser := mdns.NewMdnsService(peerhost, rendezvous, n)
	if err := ser.Start(); err != nil {
		panic(err)
	}
	return n.PeerChan
}

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
