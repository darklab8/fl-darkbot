package exposer

import (
	"darkbot/app/exposer/routes"
	_ "embed"
	"net/http"
)

func NewExposer() {
	http.ListenAndServe(":8080", routes.Server.GetMux())
}
