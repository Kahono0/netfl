package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/kahono0/netfl/handlers"
	"github.com/kahono0/netfl/p2p"
	"github.com/kahono0/netfl/repo"
)

type config struct {
	RendezvousString string
	ProtocolID       string
	ListenHost       string
	ListenPort       int
	path             string
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

func setUpRoutes(movieRepo *repo.MovieRepo) {
	http.HandleFunc("/peers", handlers.ShowPeers)
	http.HandleFunc("/send", handlers.SendSampleMsgHandler())
	http.HandleFunc("/movies", movieRepo.GetMovies)
}

func p2pConfig(c *config) *p2p.P2Pconfig {
	return &p2p.P2Pconfig{
		RendezvousString: c.RendezvousString,
		ProtocolID:       c.ProtocolID,
		ListenHost:       c.ListenHost,
		ListenPort:       c.ListenPort,
	}
}

func main() {
	config := parseFlags()

	p2p.InitConfig(config.RendezvousString, config.ProtocolID, config.ListenHost, config.ListenPort)

	p2p.Init()
	// go p2p.PingPeers(ctx, host, p2pConfig.ProtocolID)

	listener := createListener()

	defer listener.Close()

	serverPort := listener.Addr().(*net.TCPAddr).Port

	movieRepo := repo.New(serverPort, config.path, false)
	fmt.Printf("Movies:\n%s\n", movieRepo.ToJSON())
	setUpRoutes(movieRepo)

	fmt.Printf("Listening on http://localhost:%d\n", serverPort)

	log.Fatal(http.Serve(listener, nil))

	fmt.Println("Exiting...")

}
