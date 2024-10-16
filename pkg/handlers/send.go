package handlers

import (
	"time"

	"github.com/kahono0/netfl/pkg/msgs"
	"github.com/kahono0/netfl/pkg/p2p"
	"github.com/kahono0/netfl/pkg/ws"
	"github.com/libp2p/go-libp2p/core/peer"
)

func (h *Handler) InitialRequest(peer peer.AddrInfo) error {
	msg, err := msgs.NewMessage(msgs.InitialRequest, []byte{})
	if err != nil {
		return err
	}

	s, err := p2p.CreateStream(h.Host, peer, h.Config.ProtocolID)
	if err != nil {
		return err
	}

	return msg.Write(*s)
}

func (h *Handler) Ping(peer peer.AddrInfo) error {
	msg, err := msgs.NewMessage(msgs.Ping, []byte("\n"))
	if err != nil {
		return err
	}

	s, err := p2p.CreateStream(h.Host, peer, h.Config.ProtocolID)
	if err != nil {
		return err
	}

	return msg.Write(*s)
}

func (h *Handler) PingPeers() {
	for {
		peerStore := h.GetPeerStore()
		ps := peerStore.Peers
		for _, peer := range ps {
			err := h.Ping(peer.Peer)
			if err != nil {
				h.updatePeerLeft(peer)

				peerStore.RemovePeer(peer.Peer)
			}
		}

		time.Sleep(1 * time.Second)
	}
}

func (h *Handler) updatePeerLeft(peer p2p.PeerInfo) {
	removePeer := &ws.WsRemoveElement{
		ID:      peer.ID,
		Element: "p",
	}
	removePeer.Broadcast()

	// remove movies
	movieIDs := h.GetMovieRepo().RemoveForOwner(peer.ID)
	removeMovieElements := ws.WsRemoveMultipleElements{
		IDs:     movieIDs,
		Element: "m",
	}

	removeMovieElements.Broadcast()
}
