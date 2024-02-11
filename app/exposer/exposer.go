package exposer

import (
	_ "embed"
	"net/http"

	"github.com/darklab8/fl-darkbot/app/exposer/routes"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
)

func NewExposer() {
	err := http.ListenAndServe(":8080", routes.Server.GetMux())
	logus.Log.CheckError(err, "failed to listen to this port")
}
