package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kahono0/netfl/pkg/peers"
	"github.com/kahono0/netfl/repo"
	"github.com/kahono0/netfl/utils"
	"github.com/kahono0/netfl/views/pages"
)

func SetUpRoutes(movieDir string) {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Movie Dir: ", movieDir)

	// serve movieDir as static files
	http.HandleFunc("/thumb/", func(w http.ResponseWriter, r *http.Request) {
		// strip the /movies/ prefix
		fmt.Println("Request URL: ", movieDir+r.URL.Path[6:])
		http.ServeFile(w, r, movieDir+r.URL.Path[6:])
	})

	http.HandleFunc("/", index)
	http.HandleFunc("/movies", getMovies)
	http.HandleFunc("/peers", getPeers)

}

func index(w http.ResponseWriter, r *http.Request) {
	// m := repo.Movie{
	// 	Name:     "Stranger Things",
	// 	MovieUrl: "When a young boy disappears, his mother, a police chief, and his friends must confront terrifying forces in order to get him back.",
	// }
	c := pages.Index(peers.Store.Peers, repo.Repo.Movies)
	c.Render(r.Context(), w)
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	movies := repo.Repo.ToJSON()
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(movies))
}

func getPeers(w http.ResponseWriter, r *http.Request) {
	peers := peers.Store.Peers
	fmt.Printf("Peers %s\n", utils.AsPrettyJson(peers))

	data, _ := json.Marshal(peers)

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
