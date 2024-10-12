package p2p

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"

	"github.com/kahono0/netfl/msgs"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"

	"github.com/multiformats/go-multiaddr"
)

type P2Pconfig struct {
	RendezvousString string
	ProtocolID       string
	ListenHost       string
	ListenPort       int
}

func WriteMessage(rw *bufio.ReadWriter, msg *msgs.Message) error {
	_, err := rw.Write(msg.Bytes())
	if err != nil {
		return err
	}

	return rw.Flush()
}

func SendMessage(ctx context.Context, host host.Host, peerID peer.ID, msg *msgs.Message, protocalID string) error {
	s, err := host.NewStream(ctx, peerID, protocol.ID(protocalID))
	if err != nil {
		return err
	}

	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	return WriteMessage(rw, msg)
}

func handleStream(stream network.Stream) {
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go readData(rw, stream)
}

func readData(rw *bufio.ReadWriter, stream network.Stream) {
	msg, err := msgs.NewFromReader(rw.Reader)
	if err != nil {
		fmt.Println("Error reading message")
		return
	}

	err = msg.Handle(stream)
	if err != nil {
		fmt.Println("Error handling message")
		return
	}
}

func Init(ctx context.Context, cfg *P2Pconfig) *host.Host {
	fmt.Printf("[*] Listening on: %s with port: %d\n", cfg.ListenHost, cfg.ListenPort)

	r := rand.Reader

	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}

	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", cfg.ListenHost, cfg.ListenPort))

	host, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)
	if err != nil {
		panic(err)
	}

	host.SetStreamHandler(protocol.ID(cfg.ProtocolID), handleStream)

	fmt.Printf("\n[*] Your Multiaddress Is: /ip4/%s/tcp/%v/p2p/%s\n", cfg.ListenHost, cfg.ListenPort, host.ID())

	peerChan := initMDNS(host, cfg.RendezvousString)

	go listenForPeers(peerChan)

	return &host
}
