package ws

import "fmt"

type WsAddElement struct {
	ParentID string
	Content  HTMLContent
}

func (w *WsAddElement) Broadcast() {
	msg := fmt.Sprintf("%d%s,%s", AddElement, w.ParentID, w.Content.ToHTML())
	Broadcast <- []byte(msg)
}
