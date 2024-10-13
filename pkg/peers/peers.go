package peers

import (
	"fmt"
	"sync"

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

func RemovePeer(peer peer.AddrInfo) {
	peersMutex.Lock()
	defer peersMutex.Unlock()

	for i, p := range Peers {
		if p.ID == peer.ID {
			Peers = append(Peers[:i], Peers[i+1:]...)
			return
		}
	}
}

func ListenForPeers(peerChan chan peer.AddrInfo, host host.Host, protocalID string, handleNewPeer func(peer.AddrInfo, host.Host, string) error) {
	for {
		peer := <-peerChan
		fmt.Printf("Found peer: %s\n", utils.AsPrettyJson(peer))
		Peers = append(Peers, peer)

		go handleNewPeer(peer, host, protocalID)
	}
}
