package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	mhandlers "github.com/kahono0/netfl/pkg/handlers"
	"github.com/kahono0/netfl/pkg/p2p"
	"github.com/kahono0/netfl/pkg/peers"
	"github.com/kahono0/netfl/repo"
	"github.com/kahono0/netfl/router"
	"github.com/kahono0/netfl/utils"
	"github.com/libp2p/go-libp2p/core/host"
)

type config struct {
	p2p.P2PConfig
	// dir to get movies from
	Path string

	// port the server will be running on
	SPort int

	// name referring to this peer
	Alias string

	// avatar url
	Avatar string
}

func parseFlags() *config {
	f := &config{}

	flag.StringVar(&f.RendezvousString, "rendezvous", "meetme", "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	flag.StringVar(&f.ListenHost, "host", "0.0.0.0", "The bootstrap node host listen address\n")
	flag.StringVar(&f.ProtocolID, "pid", "/chat/1.1.0", "Sets a protocol id for stream headers")
	flag.IntVar(&f.ListenPort, "port", 0, "node listen port (0 pick a random unused port)")
	flag.IntVar(&f.SPort, "sport", 0, "server port")

	flag.StringVar(&f.Path, "path", "", "Path to store movie data")

	alias := utils.Whoami()
	flag.StringVar(&f.Alias, "alias", alias, "Alias for this peer")

	flag.Parse()

	return f
}

func createListener(port int) net.Listener {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	return listener
}

func setUpHandler(host host.Host, protocolID string, alias, avatar string) {
	mhandlers.NewHandler(host, protocolID, alias, avatar)

	mhandlers.MsgHandler.RegisterHandlers()

}

func main() {
	config := parseFlags()
	avatar, err := utils.GenerateIdenticon(config.Alias, 400)
	if err != nil {
		panic(err)
	}

	listener := createListener(config.SPort)

	defer listener.Close()

	serverPort := listener.Addr().(*net.TCPAddr).Port

	config.Avatar = fmt.Sprintf("http://%s:%d/%s", utils.GetPrivateIP(), serverPort, avatar)
	config.NewPeerHandler = mhandlers.HandleNewPeer
	config.StreamHandler = mhandlers.HandleStream

	host := p2p.Init(config.P2PConfig)
	setUpHandler(host, config.ProtocolID, config.Alias, config.Avatar)
	peers.NewPeerStore()

	go mhandlers.MsgHandler.PingPeers()

	repo.NewThumbNailGenWorker()
	repo.ThumbNailGenWorkerInstance.Start()
	defer repo.ThumbNailGenWorkerInstance.Close()

	repo.Init(serverPort, config.Path, false)
	// fmt.Printf("Movies:\n%s\n", repo.Repo.ToJSON())
	// router.SetUpRoutes(host, config.ProtocolID)
	router.SetUpRoutes(config.Path)

	fmt.Printf("Listening on http://localhost:%d\n", serverPort)

	log.Fatal(http.Serve(listener, nil))

	fmt.Println("Exiting...")

}
