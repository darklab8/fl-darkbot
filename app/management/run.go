/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package management

import (
	"fmt"
	"runtime"
	"time"

	"github.com/darklab8/fl-darkbot/app/configurator"
	"github.com/darklab8/fl-darkbot/app/discorder"
	"github.com/darklab8/fl-darkbot/app/forumer"
	"github.com/darklab8/fl-darkbot/app/listener"
	"github.com/darklab8/fl-darkbot/app/prometheuser"
	"github.com/darklab8/fl-darkbot/app/scrappy"
	"github.com/darklab8/fl-darkbot/app/settings"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkbot/app/viewer"

	"github.com/darklab8/go-utils/utils"

	"net/http"
	_ "net/http/pprof"

	"github.com/pkg/profile"
	"github.com/spf13/cobra"
)

// runCmd represents the check command
var runCmd = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
		logus.Log.Info("run called")

		settings.LoadEnv(settings.Environ.GetValidating())

		// migrate db
		configurator.NewConfigurator(settings.Dbpath).AutoMigrateSchema()
		forumenacer := forumer.NewForumer(settings.Dbpath)

		var scrappy_storage *scrappy.ScrappyStorage
		if settings.Env.DevEnvMockApi {
			scrappy_storage = scrappy.FixtureMockedStorage()
		} else {
			scrappy_storage = scrappy.NewScrappyWithApis()
		}
		dg := discorder.NewClient(discorder.WithWebSocket())
		scrappy_storage.GetPlayerStorage().RegisterObserve(dg) // updates number of players in Bot description

		scrappy_storage.Update()
		go scrappy_storage.Run()
		go listener.Run()
		go viewer.NewViewer(settings.Dbpath, scrappy_storage).Run()
		go forumenacer.Run()
		// probably bugged

		if settings.Env.PrometheuserOn {
			go prometheuser.Prometheuser(dg)
		}

		// profiler
		if settings.Env.ProfilingEnabled {
			p := profile.Start(profile.MemProfile, profile.ProfilePath("."), profile.NoShutdownHook)
			defer p.Stop()

			go func() {
				err := http.ListenAndServe(":8080", nil)
				logus.Log.CheckError(err, "failed to listen to 8080")
			}()
		}

		// May be garbage collector will help
		go func() {
			for {
				time.Sleep(time.Second * 10)
				runtime.GC()
			}
		}()

		fmt.Println("darkbot is launched. Awaiting Ctrl+C to disrupt")
		utils.SleepAwaitCtrlC()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
