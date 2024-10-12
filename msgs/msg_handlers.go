package msgs

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/network"
)

func HandlePing(msg *Message, stream network.Stream) error {
	// if string(msg.Data) == "" {
	// 	return nil
	// }

	if string(msg.Data) != "\n" {
		fmt.Print("-")
		// Green console colour: 	\x1b[32m
		// Reset console colour: 	\x1b[0m
		// fmt.Printf("\x1b[32m%s from %s\x1b[0m> ", string(msg.Data), stream.Conn().RemotePeer())
	}

	return nil
}

func HandleSample(msg *Message, stream network.Stream) error {
	fmt.Printf("\n\nReceived sample message %s\n\n", string(msg.Data))
	return nil
}
