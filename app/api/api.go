package api

import (
	"darkbot/app/api/routes"
	_ "embed"
	"net/http"
)

func newApi() {
	http.ListenAndServe(":8080", routes.Server.GetMux())
}
