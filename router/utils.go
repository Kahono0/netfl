package router

import (
	"context"
	"net/http"
	"strings"

	"github.com/a-h/templ"
)

func Render(ctx context.Context, c templ.Component, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	c.Render(ctx, w)
}

func RenderToString(ctx context.Context, c templ.Component) string {
	var b strings.Builder
	c.Render(ctx, &b)
	return b.String()
}
