package handlers

import (
	"fmt"

	"github.com/kahono0/netfl/pkg/msgs"
	"github.com/kahono0/netfl/pkg/peers"
	"github.com/kahono0/netfl/pkg/putils"

	"github.com/kahono0/netfl/repo"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

type Handler struct {
	Host       host.Host
	ProtocolID string
	Handlers   map[msgs.MessageTypeID]func(*msgs.Message, network.Stream) error
}

var MsgHandler *Handler

func NewHandler(host host.Host, protocolID string) {
	MsgHandler = &Handler{
		Host:       host,
		ProtocolID: protocolID,
	}
}

func (h *Handler) RegisterHandlers() {
	h.Handlers = map[msgs.MessageTypeID]func(*msgs.Message, network.Stream) error{
		msgs.Ping:           h.HandlePing,
		msgs.Sample:         h.HandleSample,
		msgs.RequestMovies:  h.HandleRequestMovies,
		msgs.ResponseMovies: h.HandleResponseMovies,
	}
}

func (h *Handler) HandleMessage(msg *msgs.Message, stream network.Stream) error {
	handler, ok := h.Handlers[msg.Type.Code]
	if !ok {
		return fmt.Errorf("no handler found for message type %d", msg.Type.Code)
	}

	return handler(msg, stream)
}

func (h *Handler) HandlePing(msg *msgs.Message, stream network.Stream) error {
	if string(msg.Data) != "\n" {
		fmt.Print("-")
	}

	return nil
}

func (h *Handler) HandleSample(msg *msgs.Message, stream network.Stream) error {
	fmt.Printf("\n\nReceived sample message %s\n\n", string(msg.Data))
	return nil
}

func (h *Handler) RequestMovies(peer peer.AddrInfo) error {
	msg, err := msgs.NewMessage(msgs.RequestMovies, []byte{})
	if err != nil {
		return err
	}

	return putils.SendMessage(h.Host, peer, msg, h.ProtocolID)
}

func (h *Handler) HandleRequestMovies(msg *msgs.Message, stream network.Stream) error {
	fmt.Printf("\n\nReceived request for movies %s\n\n", string(msg.Data))
	movies := repo.Repo.ToJSON()

	newMsg, err := msgs.NewMessage(msgs.ResponseMovies, []byte(movies))
	if err != nil {
		return err
	}

	peerID := stream.Conn().RemotePeer()
	return putils.SendToUnkown(h.Host, peerID, newMsg, h.ProtocolID)
}

func (h *Handler) HandleResponseMovies(msg *msgs.Message, stream network.Stream) error {
	movies := string(msg.Data)

	err := repo.Repo.AddFromJSON(movies)
	if err != nil {
		return err
	}

	fmt.Printf("\n\nReceived movies from peer %s\n\n", stream.Conn().RemotePeer())
	return nil
}

func HandleNewPeer(peer peer.AddrInfo, host host.Host, protocolID string) error {
	h := &Handler{
		Host:       host,
		ProtocolID: protocolID,
	}

	return h.RequestMovies(peer)
}

func (h *Handler) Ping(peer peer.AddrInfo) error {
	msg, err := msgs.NewMessage(msgs.Ping, []byte("\n"))
	if err != nil {
		return err
	}

	return putils.SendMessage(h.Host, peer, msg, h.ProtocolID)
}

func (h *Handler) PingPeers(ps []peer.AddrInfo) {
	for _, peer := range ps {
		err := h.Ping(peer)
		if err != nil {
			peers.RemovePeer(peer)
		}
	}
}
