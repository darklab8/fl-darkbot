package api

import (
	"darkbot/api/routes"
	_ "embed"
	"net/http"
)

func newApi() {
	http.ListenAndServe(":8080", routes.Server.GetMux())
}
