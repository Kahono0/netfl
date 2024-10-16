package p2p

import (
	"sync"

	"github.com/libp2p/go-libp2p/core/peer"
)

type PeerInfo struct {
	ID     int
	Peer   peer.AddrInfo
	Alias  string
	Avatar string
}

type PeerStore struct {
	Peers      []PeerInfo
	PeersMutex *sync.Mutex
	NextID     int
}

func NewStore() *PeerStore {
	return &PeerStore{
		PeersMutex: &sync.Mutex{},
		NextID:     0,
	}
}

func (ps *PeerStore) AddPeer(peer peer.AddrInfo, alias, avatar string) {
	ps.PeersMutex.Lock()
	defer ps.PeersMutex.Unlock()

	ps.Peers = append(ps.Peers, PeerInfo{ID: ps.NextID, Peer: peer, Alias: alias, Avatar: avatar})
	ps.NextID++
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
