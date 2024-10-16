package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kahono0/netfl/pkg/app"
	"github.com/kahono0/netfl/pkg/peers"
	"github.com/kahono0/netfl/repo"
	"github.com/kahono0/netfl/utils"
	"github.com/kahono0/netfl/views/pages"
)

func SetUpRoutes(app *app.App) {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", index(app))
	http.HandleFunc("/movies", getMovies(app.GetMovieRepo()))
	http.HandleFunc("/peers", getPeers(app.GetPeerStore()))
	http.HandleFunc("/thumb/", serveThumbs(app.Config.Path))

}

func index(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		peers := app.GetPeerStore().Peers
		fmt.Println(utils.AsPrettyJson(peers))
		movies := app.GetMovieRepo().Movies

		c := pages.Index(peers, movies)
		w.Header().Set("Content-Type", "text/html")
		Render(r.Context(), c, w)
	}
}

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

func serveThumbs(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request URL: %s\n", path+r.URL.Path[6:])
		http.ServeFile(w, r, path+r.URL.Path[6:])
	}
}
