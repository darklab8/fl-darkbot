package routes

import (
	"darkbot/app/exposer/web"
	_ "embed"
	"net/http"
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
