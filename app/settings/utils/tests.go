package utils

import (
	"darkbot/app/settings/darkbot_logus"
	"os"
)

func RegenerativeTest(callback func() error) error {
	if os.Getenv("DARK_TEST_REGENERATE") != "true" {
		darkbot_logus.Log.Debug("Skipping test data regenerative code")
		return nil
	}

	return callback()
}
