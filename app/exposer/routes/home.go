package routes

import (
	_ "embed"
	"net/http"

	"github.com/darklab/fl-darkbot/app/exposer/web"
)

//go:embed home.md
var homePage string

func init() {
	Server.RegisterEndpoint(
		"/",
		"404 redirect",
		func(w http.ResponseWriter, r *http.Request) {
			web.NotFoundHandler(w, r, Server.GetRouter())
		},
	)
}
