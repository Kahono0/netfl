package handlers

import (
	"bufio"
	"fmt"

	"github.com/kahono0/netfl/pkg/app"
	"github.com/kahono0/netfl/pkg/msgs"
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

func (h *Handler) HandleNewPeer(peer peer.AddrInfo) error {
	return h.InitialRequest(peer)
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
