package ws

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/a-h/templ"
	"github.com/gorilla/websocket"
)

var (
	Clients        = make(map[*websocket.Conn]bool)
	mutex          = &sync.Mutex{}
	Broadcast      = make(chan []byte)
	ImplementedIDs = []string{"p", "m"}
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

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

func Handle(w http.ResponseWriter, r *http.Request) {
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

func RenderToString(ctx context.Context, c templ.Component) string {
	var b strings.Builder
	c.Render(ctx, &b)
	return b.String()
}
