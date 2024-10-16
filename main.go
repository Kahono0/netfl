package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/kahono0/netfl/pkg/app"
	"github.com/kahono0/netfl/pkg/handlers"
	"github.com/kahono0/netfl/pkg/ws"
	"github.com/kahono0/netfl/router"
	"github.com/kahono0/netfl/utils"
)

func parseFlags() app.Config {
	f := app.Config{}

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

func createListener(port int) (net.Listener, int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	return listener, listener.Addr().(*net.TCPAddr).Port
}

func main() {
	config := parseFlags()
	avatar, err := utils.GenerateIdenticon(config.Alias, 400)
	if err != nil {
		panic(err)
	}

	listener, port := createListener(config.SPort)

	defer listener.Close()

	app, _ := app.New(config, avatar, config.Alias, port, handlers.HandleStream)

	handlers.Setup(app)

	go handlers.MsgHandler.PingPeers()

	router.SetUpRoutes(app)

	go ws.HandleBroadCasts()

	fmt.Printf("Listening on http://localhost:%d\n", port)

	log.Fatal(http.Serve(listener, nil))

	fmt.Println("Exiting...")

}
