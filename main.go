package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/kahono0/netfl/handlers"
	mhandlers "github.com/kahono0/netfl/pkg/handlers"
	"github.com/kahono0/netfl/pkg/p2p"
	"github.com/libp2p/go-libp2p/core/host"

	"github.com/kahono0/netfl/repo"
)

type config struct {
	p2p.P2PConfig
	path string
}

func parseFlags() *config {
	f := &config{}

	flag.StringVar(&f.RendezvousString, "rendezvous", "meetme", "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	flag.StringVar(&f.ListenHost, "host", "0.0.0.0", "The bootstrap node host listen address\n")
	flag.StringVar(&f.ProtocolID, "pid", "/chat/1.1.0", "Sets a protocol id for stream headers")
	flag.IntVar(&f.ListenPort, "port", 0, "node listen port (0 pick a random unused port)")

	flag.StringVar(&f.path, "path", "", "Path to store movie data")

	flag.Parse()

	return f
}

func createListener() net.Listener {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	return listener
}

func setUpRoutes(host host.Host, protocolID string) {
	http.HandleFunc("/peers", handlers.ShowPeers)
	http.HandleFunc("/send", handlers.SendSampleMsgHandler(host, protocolID))
	http.HandleFunc("/movies", repo.Repo.GetMovies)
}

func setUpHandler(host host.Host, protocolID string) {
	mhandlers.NewHandler(host, protocolID)

	mhandlers.MsgHandler.RegisterHandlers()

}

func main() {
	config := parseFlags()

	host := p2p.Init(config.P2PConfig, mhandlers.HandleNewPeer)
	setUpHandler(host, config.ProtocolID)

	go p2p.PingPeers(host)

	listener := createListener()

	defer listener.Close()

	serverPort := listener.Addr().(*net.TCPAddr).Port

	repo.Init(serverPort, config.path, false)
	fmt.Printf("Movies:\n%s\n", repo.Repo.ToJSON())
	setUpRoutes(host, config.ProtocolID)

	fmt.Printf("Listening on http://localhost:%d\n", serverPort)

	log.Fatal(http.Serve(listener, nil))

	fmt.Println("Exiting...")

}
