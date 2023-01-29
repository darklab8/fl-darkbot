/*
Loop configurator settings, and send view to coresponding channels
*/

package viewer

import (
	"darkbot/settings"
	"darkbot/utils"
)

func Run() {
	utils.LogInfo("Viewer is now running.")

	for {
		NewViewer(settings.Dbpath).Update()
	}
}
