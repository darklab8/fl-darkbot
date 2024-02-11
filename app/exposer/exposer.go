package exposer

import (
	_ "embed"
	"net/http"

	"github.com/darklab/fl-darkbot/app/exposer/routes"
)

func NewExposer() {
	http.ListenAndServe(":8080", routes.Server.GetMux())
}
