/*
Loop configurator settings, and send view to coresponding channels
*/

package viewer

import (
	"darkbot/utils"
	"fmt"
)

func Run() {

	// Query all channels

	// For each channel
	// Query all Discord messages

	// Try to grab already sent message by ID, if yes, assign to found objects with message ID.
	// Render new messages (ensure preserved Message ID)
	// Edit if message ID is present.
	// Send if not present.

	fmt.Println("Viewer is now running.  Press CTRL-C to exit.")
	utils.SleepAwaitCtrlC()
	fmt.Println("gracefully closed discord conn")
}
