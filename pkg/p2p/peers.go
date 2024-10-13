package p2p

import (
	"fmt"
	"sync"
	"time"

	"github.com/kahono0/netfl/pkg/handlers"
	"github.com/kahono0/netfl/utils"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
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

func PingPeers(host host.Host) {
	for {
		for _, peer := range Peers {
			if peer.ID == host.ID() {
				continue
			}

			err := handlers.MsgHandler.Ping(peer)
			if err != nil {
				fmt.Printf("Error pinging peer %s: %s\n", peer.ID, err)
				removePeer(peer)
			}

		}

		time.Sleep(1 * time.Second)
	}
}

func listenForPeers(peerChan chan peer.AddrInfo, host host.Host, protocalID string, handleNewPeer func(peer.AddrInfo, host.Host, string) error) {
	for {
		peer := <-peerChan
		fmt.Printf("Found peer: %s\n", utils.AsPrettyJson(peer))
		Peers = append(Peers, peer)

		go handleNewPeer(peer, host, protocalID)
	}
}
