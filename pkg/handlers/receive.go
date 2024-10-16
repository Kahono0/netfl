package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/kahono0/netfl/pkg/msgs"
	"github.com/kahono0/netfl/pkg/p2p"
	"github.com/kahono0/netfl/pkg/ws"
	"github.com/libp2p/go-libp2p/core/network"
)

func (h *Handler) HandleInitialRequest(msg *msgs.Message, stream network.Stream) error {
	fmt.Printf("\n\nReceived an initial request from%s\n\n", stream.Conn().RemotePeer())

	initialRequestData := &msgs.InitialResponseData{
		Alias:  h.Config.Alias,
		Avatar: h.Config.Avatar,
		Movies: h.GetMovieRepo().Movies,
	}

	data, err := json.Marshal(initialRequestData)
	if err != nil {
		return err
	}

	newMsg, err := msgs.NewMessage(msgs.InitialResponse, []byte(data))
	if err != nil {
		return err
	}

	peerID := stream.Conn().RemotePeer()
	peerStore := h.GetPeerStore()
	peer := peerStore.GetPeerByID(peerID.String())
	if peer == nil {
		return fmt.Errorf("peer not found")
	}

	s, err := p2p.CreateStream(h.Host, *peer, h.Config.ProtocolID)
	if err != nil {
		return err
	}

	return newMsg.Write(*s)
}

func (h *Handler) HandleInitialResponse(msg *msgs.Message, stream network.Stream) error {
	fmt.Printf("\n\nReceived an initial response from %s\n\n", stream.Conn().RemotePeer())

	var data msgs.InitialResponseData

	err := json.Unmarshal(msg.Data, &data)
	if err != nil {
		return err
	}

	fmt.Printf("Received alias: %s, avatar: %s\n", data.Alias, data.Avatar)

	peerStore := h.GetPeerStore()
	peer := peerStore.UpdatePeer(stream.Conn().RemotePeer(), data.Alias, data.Avatar)
	if peer != nil {
		peerHTML := &ws.PeerInfoHTML{
			PeerInfo: *peer,
		}

		wsAddElement := &ws.WsAddElement{
			ParentID: "peers",
			Content:  peerHTML,
		}

		wsAddElement.Broadcast()
	}

	h.GetMovieRepo().AddMovies(peer.ID, data.Movies)

	return nil

}
