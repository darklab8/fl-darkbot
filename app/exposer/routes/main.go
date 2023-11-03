package routes

import (
	"darkbot/app/exposer/web"
	"net/http"
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
