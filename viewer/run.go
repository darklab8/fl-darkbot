/*
Loop configurator settings, and send view to coresponding channels
*/

package viewer

import (
	"darkbot/settings"
	"darkbot/utils/logger"
)

func Run() {
	logger.Info("Viewer is now running.")

	for {
		NewViewer(settings.Dbpath).Update()
	}
}
