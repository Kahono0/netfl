package p2p

import (
	"crypto/rand"
	"fmt"

	"github.com/kahono0/netfl/pkg/peers"
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
	StreamHandler    func(network.Stream)
	NewPeerHandler   func(peer.AddrInfo, host.Host, string) error
}

func Init(cfg P2PConfig) host.Host {
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

	host.SetStreamHandler(protocol.ID(cfg.ProtocolID), cfg.StreamHandler)

	fmt.Printf("\n[*] Your Multiaddress Is: /ip4/%s/tcp/%v/p2p/%s\n", cfg.ListenHost, cfg.ListenPort, host.ID())

	peerChan := initMDNS(host, cfg.RendezvousString)

	go peers.ListenForPeers(peerChan, host, cfg.ProtocolID, cfg.NewPeerHandler)

	return host
}
