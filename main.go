package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/kahono0/netfl/handlers"
	"github.com/kahono0/netfl/p2p"
	"github.com/kahono0/netfl/repo"
)

func createListener() net.Listener {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	return listener
}

func setUpRoutes(movieRepo *repo.MovieRepo) {
	http.HandleFunc("/peers", handlers.ShowPeers)
	http.HandleFunc("/send", handlers.SendSampleMsg)
	http.HandleFunc("/movies", movieRepo.GetMovies)
}

func main() {
	ctx := context.Background()

	config := parseFlags()
	host := p2p.Init(ctx, &config.P2Pconfig)
	go p2p.PingPeers(ctx, *host, &config.P2Pconfig)

	listener := createListener()

	defer listener.Close()

	serverPort := listener.Addr().(*net.TCPAddr).Port

	movieRepo := repo.New(serverPort, ".", false)
	setUpRoutes(movieRepo)

	fmt.Printf("Listening on http://localhost:%d\n", serverPort)

	log.Fatal(http.Serve(listener, nil))

	fmt.Println("Exiting...")

}
