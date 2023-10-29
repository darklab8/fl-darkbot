/*
Loop configurator settings, and send view to coresponding channels
*/

package viewer

import (
	"darkbot/settings"
	"darkbot/settings/utils/logger"
)

func Run() {
	logger.Info("Viewer is now running.")

	view := NewViewer(settings.Dbpath)
	for {
		view.Update()
	}
}
