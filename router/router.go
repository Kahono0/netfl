package router

import (
	"net/http"

	"github.com/kahono0/netfl/repo"
	"github.com/kahono0/netfl/views/pages"
)

func SetUpRoutes() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", index)

}

func index(w http.ResponseWriter, r *http.Request) {
	m := repo.Movie{
		Name:     "Stranger Things",
		MovieUrl: "When a young boy disappears, his mother, a police chief, and his friends must confront terrifying forces in order to get him back.",
	}

	c := pages.Index(m)
	c.Render(r.Context(), w)
}
