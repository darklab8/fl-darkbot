/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package management

import (
	"darkbot/app/configurator"
	"darkbot/app/listener"
	"darkbot/app/scrappy"
	"darkbot/app/settings"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/utils"
	"darkbot/app/viewer"
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
		logus.Info("run called")

		// migrate db
		configurator.NewConfigurator(settings.Dbpath).AutoMigrateSchema()

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
