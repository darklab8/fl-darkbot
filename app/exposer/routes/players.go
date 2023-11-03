package routes

import "net/http"

func init() {
	Server.RegisterEndpoint(
		"/players",
		"players online",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Another page"))
		},
	)
}
