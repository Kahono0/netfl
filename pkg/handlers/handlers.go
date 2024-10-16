package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"

	"github.com/kahono0/netfl/pkg/app"
	"github.com/kahono0/netfl/pkg/msgs"
	"github.com/kahono0/netfl/pkg/putils"
	"github.com/kahono0/netfl/utils"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

type Handler struct {
	*app.App
	Handlers map[msgs.MessageTypeID]func(*msgs.Message, network.Stream) error
}

var MsgHandler *Handler

func Setup(app *app.App) {
	NewHandler(app)

	MsgHandler.RegisterHandlers()

	go MsgHandler.ListenForPeers()
}

func NewHandler(app *app.App) {
	MsgHandler = &Handler{
		App: app,
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

	s, err := putils.CreateStream(h.Host, peer, h.Config.ProtocolID)
	if err != nil {
		return err
	}

	return msg.Write(*s)
}

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

	s, err := putils.CreateStream(h.Host, *peer, h.Config.ProtocolID)
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
	_ = peerStore.UpdatePeer(stream.Conn().RemotePeer(), data.Alias, data.Avatar)

	h.GetMovieRepo().AddMovies(data.Movies)

	return nil

}

func (h *Handler) HandleNewPeer(peer peer.AddrInfo) error {
	peerStore := h.GetPeerStore()
	peerStore.AddPeer(peer, "", "")

	err := h.InitialRequest(peer)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) Ping(peer peer.AddrInfo) error {
	msg, err := msgs.NewMessage(msgs.Ping, []byte("\n"))
	if err != nil {
		return err
	}

	s, err := putils.CreateStream(h.Host, peer, h.Config.ProtocolID)
	if err != nil {
		return err
	}

	return msg.Write(*s)
}

func (h *Handler) PingPeers() {
	peerStore := h.GetPeerStore()
	ps := peerStore.Peers
	for _, peer := range ps {
		err := h.Ping(peer.Peer)
		if err != nil {
			peerStore.RemovePeer(peer.Peer)
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

func (h *Handler) ListenForPeers() {
	for {
		peer := <-h.PeerChan
		fmt.Printf("Found peer: %s\n", utils.AsPrettyJson(peer))
		ps := h.GetPeerStore()
		ps.AddPeer(peer, "", "")

		go h.HandleNewPeer(peer)
	}
}
