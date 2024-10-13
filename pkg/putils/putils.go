package putils

import (
	"bufio"
	"context"
	"fmt"

	"github.com/kahono0/netfl/pkg/msgs"
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

func WriteMessage(rw *bufio.ReadWriter, msg *msgs.Message) error {
	_, err := rw.Write(msg.Bytes())
	if err != nil {
		return err
	}

	return rw.Flush()
}

func SendMessage(host host.Host, peer peer.AddrInfo, msg *msgs.Message, protocolID string) error {
	ctx := context.Background()
	if err := Connect(host, peer); err != nil {
		return err
	}
	s, err := host.NewStream(ctx, peer.ID, protocol.ID(protocolID))
	if err != nil {
		return err
	}

	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	return WriteMessage(rw, msg)
}

func SendToUnkown(host host.Host, peerID peer.ID, msg *msgs.Message, protocolID string) error {
	peer := peers.GetPeerByID(peerID.String())
	if peer == nil {
		return fmt.Errorf("no peer found with ID %s", peerID)
	}

	return SendMessage(host, *peer, msg, protocolID)
}

func SendWithStream(msg *msgs.Message, stream network.Stream) error {
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	return WriteMessage(rw, msg)
}
