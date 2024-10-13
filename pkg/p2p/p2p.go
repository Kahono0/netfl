package p2p

import (
	"bufio"
	"crypto/rand"
	"fmt"

	"github.com/kahono0/netfl/pkg/handlers"
	"github.com/kahono0/netfl/pkg/msgs"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"

	"github.com/multiformats/go-multiaddr"
)

type P2PConfig struct {
	RendezvousString string
	ProtocolID       string
	ListenHost       string
	ListenPort       int
}

func handleStream(stream network.Stream) {
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	defer stream.Close()

	readData(rw, stream)
}

func readData(rw *bufio.ReadWriter, stream network.Stream) {
	msg, err := msgs.NewFromReader(rw.Reader)
	if err != nil {
		fmt.Println("Error reading message")
		return
	}
	err = handlers.MsgHandler.HandleMessage(msg, stream)
	if err != nil {
		fmt.Println("Error handling message")
		return
	}
}

func Init(cfg P2PConfig, handleNewPeer func(peer.AddrInfo, host.Host, string) error) host.Host {
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

	go listenForPeers(peerChan, host, cfg.ProtocolID, handleNewPeer)

	return host
}
