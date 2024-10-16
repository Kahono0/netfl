package ws

import "fmt"

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

type WsRemoveMultipleElements struct {
	IDs     []int
	Element string
}

func (w *WsRemoveMultipleElements) Broadcast() {
	for _, t := range ImplementedIDs {
		if t == w.Element {
			msg := fmt.Sprintf("%d", RemoveMultiple)
			for _, id := range w.IDs {
				msg += fmt.Sprintf("%s-%d,", w.Element, id)
			}

			msg = msg[:len(msg)-1]

			Broadcast <- []byte(msg)
		}
	}
}
