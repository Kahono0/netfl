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

var Host host.Host

type P2Pconfig struct {
	RendezvousString string
	ProtocolID       string
	ListenHost       string
	ListenPort       int
}

var config *P2Pconfig

func InitConfig(rendezvousString, protocolID, listenHost string, listenPort int) {
	config = &P2Pconfig{
		RendezvousString: rendezvousString,
		ProtocolID:       protocolID,
		ListenHost:       listenHost,
		ListenPort:       listenPort,
	}
}

func WriteMessage(rw *bufio.ReadWriter, msg *msgs.Message) error {
	_, err := rw.Write(msg.Bytes())
	if err != nil {
		return err
	}

	return rw.Flush()
}

func SendMessage(peer peer.AddrInfo, msg *msgs.Message) error {
	ctx := context.Background()
	if err := Connect(Host, peer); err != nil {
		return err
	}
	s, err := Host.NewStream(ctx, peer.ID, protocol.ID(config.ProtocolID))
	if err != nil {
		return err
	}

	return SendWithStream(s, msg)
}

func SendWithStream(stream network.Stream, msg *msgs.Message) error {
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	return WriteMessage(rw, msg)
}

func handleStream(stream network.Stream) {
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	readData(rw, stream)
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

func Init() {
	fmt.Printf("[*] Listening on: %s with port: %d\n", config.ListenHost, config.ListenPort)

	r := rand.Reader

	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}

	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", config.ListenHost, config.ListenPort))

	host, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)
	if err != nil {
		panic(err)
	}

	host.SetStreamHandler(protocol.ID(config.ProtocolID), handleStream)

	fmt.Printf("\n[*] Your Multiaddress Is: /ip4/%s/tcp/%v/p2p/%s\n", config.ListenHost, config.ListenPort, host.ID())

	peerChan := initMDNS(host, config.RendezvousString)

	go listenForPeers(peerChan)

	Host = host
}
