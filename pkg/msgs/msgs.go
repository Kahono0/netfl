package msgs

import (
	"bufio"
	"errors"

	"github.com/kahono0/netfl/repo"
	"github.com/libp2p/go-libp2p/core/network"
)

var ErrUnknownMessageType = errors.New("unknown message type")

var EOFDelim = byte(0)

type MessageTypeID uint

type MessageType struct {
	Code MessageTypeID
}

const (
	Ping MessageTypeID = iota
	Sample
	InitialRequest
	InitialResponse
)

type InitialResponseData struct {
	Alias  string
	Avatar string
	Movies []repo.Movie
}

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
	msg, err := r.ReadBytes(EOFDelim)
	if err != nil {
		return nil, err
	}

	return DecodeMessage(msg)
}

func (m *Message) bytes() []byte {
	msg := append([]byte{byte(m.Type.Code)}, m.Data...)
	return append(msg, EOFDelim)
}

func (m *Message) Write(stream network.Stream) error {
	defer stream.Close()

	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	_, err := rw.Write(m.bytes())
	if err != nil {
		return err
	}

	return rw.Flush()
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
