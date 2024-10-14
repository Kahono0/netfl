package router

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
)

func Render(ctx context.Context, c templ.Component, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	c.Render(ctx, w)
}
