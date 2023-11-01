/*
Loop configurator settings, and send view to coresponding channels
*/

package viewer

import (
	"darkbot/app/settings"
	"darkbot/app/settings/logus"
)

func Run() {
	logus.Info("Viewer is now running.")

	view := NewViewer(settings.Dbpath)
	for {
		view.Update()
	}
}
