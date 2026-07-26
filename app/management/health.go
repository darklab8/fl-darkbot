package management

import (
	"net/http"
	"time"

	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/go-utils/typelog"
	"github.com/spf13/cobra"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Health command",
	Run: func(cmd *cobra.Command, args []string) {
		tr := &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    10 * time.Second,
			DisableCompression: true,
		}

		client := &http.Client{Transport: tr}
		resp, err := client.Get("http://localhost:8000/metrics")
		logus.Log.CheckPanic(err, "failed to health check")
		if resp.StatusCode != 200 {
			logus.Log.Panic("status code is not 200", typelog.Any("code", resp.StatusCode))
		}
		logus.Log.Debug("service is healthy")
	},
}

func init() {
	rootCmd.AddCommand(healthCmd)
}
