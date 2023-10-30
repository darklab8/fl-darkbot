/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package management

import (
	"darkbot/configurator"
	"darkbot/listener"
	"darkbot/scrappy"
	"darkbot/settings"
	"darkbot/settings/utils"
	"darkbot/settings/utils/logger"
	"darkbot/viewer"
	"runtime"
	"time"

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

		if settings.Config.DevEnvMockApi == "true" {
			scrappy.Storage = scrappy.FixtureMockedStorage()
			scrappy.Storage.Update()
		}

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

		// May be garbage collector will help
		go func() {
			for {
				time.Sleep(time.Second * 10)
				runtime.GC()
			}
		}()

		utils.SleepAwaitCtrlC()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
