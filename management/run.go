/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package management

import (
	"darkbot/configurator"
	"darkbot/listener"
	"darkbot/scrappy"
	"darkbot/settings"
	"darkbot/utils"
	"darkbot/utils/logger"
	"darkbot/viewer"

	"net/http"
	_ "net/http/pprof"

	"github.com/pkg/profile"
	"github.com/spf13/cobra"
)

// runCmd represents the check command
var runCmd = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("run called")

		// migrate db
		configurator.NewConfigurator(settings.Dbpath).Migrate()

		scrappy.Storage.Update()
		go scrappy.Run()
		go listener.Run()
		go viewer.Run()

		// profiler
		if settings.Config.ProfilingEnabled == settings.EnvTrue {
			p := profile.Start(profile.MemProfile, profile.ProfilePath("."), profile.NoShutdownHook)
			defer p.Stop()

			go func() {
				http.ListenAndServe(":8080", nil)
			}()
		}

		utils.SleepAwaitCtrlC()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
