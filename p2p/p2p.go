package p2p

import (
	"bufio"
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/kahono0/netfl/p2p/msgs"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"

	"github.com/multiformats/go-multiaddr"
)

var peers []peer.AddrInfo
var peersMutex sync.Mutex

func listenForPeers(peerChan chan peer.AddrInfo) {
	for {
		peer := <-peerChan
		fmt.Printf("Found peer: %v\n", peer)
		peers = append(peers, peer)
	}
}

func removePeer(peer peer.AddrInfo) {
	peersMutex.Lock()
	defer peersMutex.Unlock()

	for i, p := range peers {
		if p.ID == peer.ID {
			peers = append(peers[:i], peers[i+1:]...)
			return
		}
	}
}

func pingPeers(ctx context.Context, host host.Host, cfg *config) {
	for {
		for _, peer := range peers {
			if peer.ID == host.ID() {
				continue
			}

			if err := host.Connect(ctx, peer); err != nil {
				fmt.Println("Connection failed:", err)
				removePeer(peer)
				continue
			}

			s, err := host.NewStream(ctx, peer.ID, protocol.ID(cfg.ProtocolID))
			if err != nil {
				fmt.Printf("Error opening stream to %s: %s\n", peer.ID, err)
				removePeer(peer)
				continue
			}

			rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
			err = ping(rw)
			if err != nil {
				fmt.Printf("Error pinging %s: %s\n", peer.ID, err)
				continue
			}

			s.Close()
		}

		time.Sleep(1 * time.Second)
	}
}

func ping(rw *bufio.ReadWriter) error {
	msg := msgs.NewMessage(msgs.Ping, []byte("ping"))

	_, err := rw.Write(msg.Bytes())
	if err != nil {
		fmt.Println("Error writing to buffer")
		return err
	}

	err = rw.Flush()
	if err != nil {
		fmt.Println("Error flushing buffer")
		return err
	}

	return nil
}

func handleStream(stream network.Stream) {
	fmt.Println("Got a new stream!")

	// Create a buffer stream for non-blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go readData(rw, stream)
	// go writeData(rw)

	// 'stream' will stay open until you close it (or the other side closes it).
}

func readData(rw *bufio.ReadWriter, stream network.Stream) {
	bytes, err := rw.ReadBytes('\n')
	if err != nil {
		fmt.Println("Error reading from buffer")
		panic(err)
	}

	msg := msgs.DecodeMessage(bytes)
	msg.Handle(stream)
}

func Init() {
	help := flag.Bool("help", false, "Display Help")
	cfg := parseFlags()

	if *help {
		fmt.Printf("Simple example for peer discovery using mDNS. mDNS is great when you have multiple peers in local LAN.")
		fmt.Printf("Usage: \n   Run './chat-with-mdns'\nor Run './chat-with-mdns -host [host] -port [port] -rendezvous [string] -pid [proto ID]'\n")

		os.Exit(0)
	}

	fmt.Printf("[*] Listening on: %s with port: %d\n", cfg.listenHost, cfg.listenPort)

	// ctx := context.Background()
	r := rand.Reader

	// Creates a new RSA key pair for this host.
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}

	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", cfg.listenHost, cfg.listenPort))

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	host, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)
	if err != nil {
		panic(err)
	}

	// Set a function as stream handler.
	// This function is called when a peer initiates a connection and starts a stream with this peer.
	host.SetStreamHandler(protocol.ID(cfg.ProtocolID), handleStream)

	fmt.Printf("\n[*] Your Multiaddress Is: /ip4/%s/tcp/%v/p2p/%s\n", cfg.listenHost, cfg.listenPort, host.ID())

	peerChan := initMDNS(host, cfg.RendezvousString)

	go listenForPeers(peerChan)

	ctx := context.Background()
	go pingPeers(ctx, host, cfg)

	select {}
}
