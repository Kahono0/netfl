package peers

import (
	"fmt"
	"sync"

	"github.com/kahono0/netfl/utils"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

type PeerInfo struct {
	Peer   peer.AddrInfo
	Alias  string
	Avatar string
}

type PeerStore struct {
	Peers      []PeerInfo
	PeersMutex sync.Mutex
}

var Store *PeerStore

func NewPeerStore() {
	Store = &PeerStore{}
}

func (ps *PeerStore) AddPeer(peer peer.AddrInfo, alias, avatar string) {
	ps.PeersMutex.Lock()
	defer ps.PeersMutex.Unlock()

	ps.Peers = append(ps.Peers, PeerInfo{Peer: peer, Alias: alias, Avatar: avatar})
}

func (ps *PeerStore) RemovePeer(peer peer.AddrInfo) {
	ps.PeersMutex.Lock()
	defer ps.PeersMutex.Unlock()

	for i, p := range ps.Peers {
		if p.Peer.ID == peer.ID {
			ps.Peers = append(ps.Peers[:i], ps.Peers[i+1:]...)
			return
		}
	}
}

func (ps *PeerStore) GetPeerByID(peerID string) *peer.AddrInfo {
	for _, p := range ps.Peers {
		if p.Peer.ID.String() == peerID {
			return &p.Peer
		}
	}

	return nil
}

func (ps *PeerStore) UpdatePeer(peerID peer.ID, alias string, avatar string) *PeerInfo {
	ps.PeersMutex.Lock()
	defer ps.PeersMutex.Unlock()

	for i, p := range ps.Peers {
		if p.Peer.ID == peerID {
			ps.Peers[i].Alias = alias
			ps.Peers[i].Avatar = avatar
			return &ps.Peers[i]
		}
	}

	return nil
}

func ListenForPeers(peerChan chan peer.AddrInfo, host host.Host, protocalID string, handleNewPeer func(peer.AddrInfo, host.Host, string) error) {
	for {
		peer := <-peerChan
		fmt.Printf("Found peer: %s\n", utils.AsPrettyJson(peer))
		if Store == nil {
			panic("Peer store not initialized")
		}

		Store.AddPeer(peer, "", "")

		go handleNewPeer(peer, host, protocalID)
	}
}
