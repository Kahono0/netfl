package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/kahono0/netfl/pkg/p2p"
	"github.com/kahono0/netfl/views/pages"
)

type WsOps uint

const (
	RemoveElement WsOps = iota
	AddElement
)

type BroadcastInterface interface {
	Broadcast()
}

type HTMLContent interface {
	ToHTML() string
}

type WsRemoveElement struct {
	ID      int
	Element string
}

func (w *WsRemoveElement) Broadcast() {
	handledTypes := []string{"p", "m"}
	for _, t := range handledTypes {
		if t == w.Element {
			msg := fmt.Sprintf("%d%s-%d", RemoveElement, w.Element, w.ID)
			Broadcast <- []byte(msg)
		}
	}
}

type PeerInfoHTML struct {
	p2p.PeerInfo
}

func (p *PeerInfoHTML) ToHTML() string {
	c := pages.OnlineUser(p.PeerInfo)
	return RenderToString(context.TODO(), c)
}

type WsAddElement struct {
	ParentID string
	Content  HTMLContent
}

func (w *WsAddElement) Broadcast() {
	msg := fmt.Sprintf("%d%s,%s", AddElement, w.ParentID, w.Content.ToHTML())
	Broadcast <- []byte(msg)

}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

var Clients = make(map[*websocket.Conn]bool)
var mutex = &sync.Mutex{}

var Broadcast = make(chan []byte)

func HandleBroadCasts() {
	for {
		msg := <-Broadcast

		mutex.Lock()
		for client := range Clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				client.Close()
				delete(Clients, client)
			}
		}

		mutex.Unlock()
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	mutex.Lock()
	Clients[conn] = true
	mutex.Unlock()

	// Continuously read messages from the client
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("Received: %s", msg)

		// Send the message back to the client
		response := fmt.Sprintf("Echo: %s", msg)
		if err = conn.WriteMessage(msgType, []byte(response)); err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
