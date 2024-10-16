package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kahono0/netfl/pkg/app"
	"github.com/kahono0/netfl/pkg/peers"
	"github.com/kahono0/netfl/repo"
	"github.com/kahono0/netfl/utils"
)

func SetUpRoutes(app *app.App) {
	http.HandleFunc("/movies", getMovies(app.GetMovieRepo()))
	http.HandleFunc("/peers", getPeers(app.GetPeerStore()))
}

// func SetUpRoutes(movieDir string) {
// 	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
// 	fmt.Println("Movie Dir: ", movieDir)

// 	// serve movieDir as static files
// 	http.HandleFunc("/thumb/", func(w http.ResponseWriter, r *http.Request) {
// 		// strip the /movies/ prefix
// 		fmt.Println("Request URL: ", movieDir+r.URL.Path[6:])
// 		http.ServeFile(w, r, movieDir+r.URL.Path[6:])
// 	})

// 	http.HandleFunc("/", index)
// 	http.HandleFunc("/movies", getMovies)
// 	http.HandleFunc("/peers", getPeers)

// }

func getMovies(repo *repo.MovieRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		movies := repo.ToJSON()
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(movies))
	}
}

func getPeers(peerStore *peers.PeerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		peers := peerStore.Peers
		fmt.Printf("Peers %s\n", utils.AsPrettyJson(peers))

		data, _ := json.Marshal(peers)

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}
