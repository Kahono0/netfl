package main

import (
	"fmt"
	"net/http"

	"github.com/kahono0/netfl/p2p"
	"github.com/kahono0/netfl/repo"
)

func main() {

	movieRepo := repo.NewMovieRepo(".", "http://localhost:8080", false)
	movieRepo.Load()
	fmt.Println(movieRepo)

	go p2p.Init()

	http.HandleFunc("/movies", movieRepo.GetMovies)

	http.ListenAndServe(":8080", nil)
}
