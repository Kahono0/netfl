package p2p

import (
	"bufio"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kahono0/netfl/msgs"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

var Peers []peer.AddrInfo
var peersMutex sync.Mutex

func GetPeerByID(peerID string) *peer.AddrInfo {
	peersMutex.Lock()
	defer peersMutex.Unlock()

	for _, p := range Peers {
		if p.ID.String() == peerID {
			return &p
		}
	}

	return nil
}

func listenForPeers(peerChan chan peer.AddrInfo) {
	for {
		peer := <-peerChan
		fmt.Printf("Found peer: %v\n", peer)
		Peers = append(Peers, peer)
	}
}

func removePeer(peer peer.AddrInfo) {
	peersMutex.Lock()
	defer peersMutex.Unlock()

	for i, p := range Peers {
		if p.ID == peer.ID {
			Peers = append(Peers[:i], Peers[i+1:]...)
			return
		}
	}
}

func PingPeers(ctx context.Context, host host.Host, cfg *P2Pconfig) {
	for {
		for _, peer := range Peers {
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
	msg, err := msgs.NewMessage(msgs.Ping, []byte("ping"))
	if err != nil {
		fmt.Println("Error creating message")
		return err
	}

	_, err = rw.Write(msg.Bytes())
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
