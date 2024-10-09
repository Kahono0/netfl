package msgs

import "github.com/libp2p/go-libp2p/core/network"

type MessageType uint

const (
	Ping MessageType = iota
)

// all messages must include '\n' at the end and '\n' signifies end of message
type Message struct {
	Type MessageType
	Data []byte
}

func NewMessage(t MessageType, data []byte) *Message {
	return &Message{
		Type: t,
		Data: data,
	}
}

func (m *Message) Bytes() []byte {
	msg := append([]byte{byte(m.Type)}, m.Data...)
	return append(msg, '\n')
}

func DecodeMessage(data []byte) *Message {
	return &Message{
		Type: MessageType(data[0]),
		Data: data[1 : len(data)-1],
	}
}

func (m *Message) Handle(stream network.Stream) {
	switch m.Type {
	case Ping:
		HandlePing(m, stream)
	}
}
