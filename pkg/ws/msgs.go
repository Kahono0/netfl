package ws

import (
	"context"

	"github.com/kahono0/netfl/pkg/p2p"
	"github.com/kahono0/netfl/views/pages"
)

type WsOps uint

const (
	RemoveElement WsOps = iota
	AddElement
	RemoveMultiple
)

type BroadcastInterface interface {
	Broadcast()
}

type HTMLContent interface {
	ToHTML() string
}

type PeerInfoHTML struct {
	p2p.PeerInfo
}

func (p *PeerInfoHTML) ToHTML() string {
	c := pages.OnlineUser(p.PeerInfo)
	return RenderToString(context.TODO(), c)
}
