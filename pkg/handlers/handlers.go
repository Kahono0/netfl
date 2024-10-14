package handlers

import (
	"bufio"
	"encoding/json"
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
	Alias      string
	Avatar     string
	Handlers   map[msgs.MessageTypeID]func(*msgs.Message, network.Stream) error
}

var MsgHandler *Handler

func NewHandler(host host.Host, protocolID string, alias, avatar string) {
	MsgHandler = &Handler{
		Host:       host,
		ProtocolID: protocolID,
		Alias:      alias,
		Avatar:     avatar,
	}
}

func (h *Handler) RegisterHandlers() {
	h.Handlers = map[msgs.MessageTypeID]func(*msgs.Message, network.Stream) error{
		msgs.Ping:            h.HandlePing,
		msgs.Sample:          h.HandleSample,
		msgs.InitialRequest:  h.HandleInitialRequest,
		msgs.InitialResponse: h.HandleInitialResponse,
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

func (h *Handler) InitialRequest(peer peer.AddrInfo) error {
	msg, err := msgs.NewMessage(msgs.InitialRequest, []byte{})
	if err != nil {
		return err
	}

	s, err := putils.CreateStrean(h.Host, peer, h.ProtocolID)
	if err != nil {
		return err
	}

	return msg.Write(*s)
}

func (h *Handler) HandleInitialRequest(msg *msgs.Message, stream network.Stream) error {
	fmt.Printf("\n\nReceived an initial request from%s\n\n", stream.Conn().RemotePeer())

	initialRequestData := &msgs.InitialResponseData{
		Alias:  h.Alias,
		Avatar: h.Avatar,
		Movies: repo.Repo.GetMovies(),
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
	s, err := putils.CreateStreamFromUnknow(h.Host, peerID, h.ProtocolID)
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

	_ = peers.Store.UpdatePeer(stream.Conn().RemotePeer(), data.Alias, data.Avatar)

	repo.Repo.AddMovies(data.Movies)

	return nil

}

func HandleNewPeer(peer peer.AddrInfo, host host.Host, protocolID string) error {
	h := &Handler{
		Host:       host,
		ProtocolID: protocolID,
	}

	return h.InitialRequest(peer)

}

func (h *Handler) Ping(peer peer.AddrInfo) error {
	msg, err := msgs.NewMessage(msgs.Ping, []byte("\n"))
	if err != nil {
		return err
	}

	s, err := putils.CreateStrean(h.Host, peer, h.ProtocolID)
	if err != nil {
		return err
	}

	return msg.Write(*s)
}

func (h *Handler) PingPeers() {
	ps := peers.Store.Peers
	for _, peer := range ps {
		err := h.Ping(peer.Peer)
		if err != nil {
			peers.Store.RemovePeer(peer.Peer)
		}
	}
}

func handleMsg(rw *bufio.ReadWriter, stream network.Stream) {
	msg, err := msgs.NewFromReader(rw.Reader)
	if err != nil {
		fmt.Println("Error reading message")
		return
	}
	err = MsgHandler.HandleMessage(msg, stream)
	if err != nil {
		fmt.Println("Error handling message")
		return
	}
}

func HandleStream(stream network.Stream) {
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	defer stream.Close()

	handleMsg(rw, stream)
}
