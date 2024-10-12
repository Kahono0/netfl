package main

import (
	"flag"

	"github.com/kahono0/netfl/p2p"
)

type Flags struct {
	p2p.P2Pconfig
	path string
}

func parseFlags() *Flags {
	f := &Flags{}

	flag.StringVar(&f.P2Pconfig.RendezvousString, "rendezvous", "meetme", "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	flag.StringVar(&f.P2Pconfig.ListenHost, "host", "0.0.0.0", "The bootstrap node host listen address\n")
	flag.StringVar(&f.P2Pconfig.ProtocolID, "pid", "/chat/1.1.0", "Sets a protocol id for stream headers")
	flag.IntVar(&f.P2Pconfig.ListenPort, "port", 0, "node listen port (0 pick a random unused port)")

	flag.StringVar(&f.path, "path", "", "Path to store movie data")

	flag.Parse()

	return f
}
