package msgs

import (
	"bufio"
	"errors"
)

var ErrUnknownMessageType = errors.New("unknown message type")

type MessageTypeID uint

type MessageType struct {
	Code MessageTypeID
}

const (
	Ping MessageTypeID = iota
	Sample
	RequestMovies
	ResponseMovies
)

// utility to chack if method type is valid
func NewMessageType(t MessageTypeID) (*MessageType, error) {
	return &MessageType{
		Code: t,
	}, nil
}

// Message structure
// -------------------------------------------------
// | Type (1 byte) | Data (variable) | \n (1 byte) |
// -------------------------------------------------
type Message struct {
	Type MessageType
	Data []byte
}

func NewMessage(c MessageTypeID, data []byte) (*Message, error) {
	msgType, error := NewMessageType(c)
	if error != nil {
		return nil, error
	}

	return &Message{
		Type: *msgType,
		Data: data,
	}, nil
}

func NewFromReader(r *bufio.Reader) (*Message, error) {
	msg, err := r.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	return DecodeMessage(msg)
}

func (m *Message) Bytes() []byte {
	msg := append([]byte{byte(m.Type.Code)}, m.Data...)
	return append(msg, '\n')
}

func DecodeMessage(data []byte) (*Message, error) {
	code := MessageTypeID(data[0])
	msgType, err := NewMessageType(code)
	if err != nil {
		return nil, err
	}

	return &Message{
		Type: *msgType,
		Data: data[1 : len(data)-1],
	}, nil
}
