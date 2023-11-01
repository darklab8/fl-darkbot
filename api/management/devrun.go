package main

import (
	"darkbot/api/routes"
	_ "embed"
	"net/http"
)

func main() {
	http.ListenAndServe(":8080", newApi())
}

func newApi() *http.ServeMux {
	return routes.Server.GetMux()
}
