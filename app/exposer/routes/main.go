package routes

import (
	"net/http"

	"github.com/darklab8/fl-darkbot/app/exposer/web"
)

var Server *web.Server = web.NewServer()

func init() {
	Server.RegisterEndpoint(
		"/ping",
		"api healthcheck",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong!"))
		},
	)
}
