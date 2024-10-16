package router

import (
	"encoding/json"
	"net/http"

	"github.com/kahono0/netfl/pkg/app"
	"github.com/kahono0/netfl/pkg/p2p"
	"github.com/kahono0/netfl/pkg/repo/movies"
	"github.com/kahono0/netfl/pkg/ws"
	"github.com/kahono0/netfl/views/pages"
)

func SetUpRoutes(app *app.App) {
	http.HandleFunc("/ws", ws.Handle)
	http.HandleFunc("/avi", serveAvi(app.Config.Alias))
	http.HandleFunc("/", index(app))
	http.HandleFunc("/movies", getMovies(app.GetMovieRepo()))
	http.HandleFunc("/peers", getPeers(app.GetPeerStore()))
	http.HandleFunc("/thumb/", serveThumbs(app.Config.Path))
	http.HandleFunc("/movies/", serveMovies(app.Config.Path))

}

func index(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		peers := app.GetPeerStore().Peers
		movies := app.GetMovieRepo().Movies

		c := pages.Index(peers, movies)
		w.Header().Set("Content-Type", "text/html")
		Render(r.Context(), c, w)
	}
}

func getMovies(repo *movies.MovieRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		movies := repo.ToJSON()
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(movies))
	}
}

func getPeers(peerStore *p2p.PeerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		peers := peerStore.Peers

		data, _ := json.Marshal(peers)

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func serveThumbs(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path+r.URL.Path[6:])
	}
}

func serveMovies(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path+r.URL.Path[7:])
	}
}

func serveAvi(alias string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file := "~/.netfl/assets/" + alias + ".png"
		http.ServeFile(w, r, file)
	}
}
