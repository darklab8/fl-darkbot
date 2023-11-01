package utils

import (
	"darkbot/app/settings/logus"
	"os"
)

func RegenerativeTest(callback func() error) error {
	if os.Getenv("DARK_TEST_REGENERATE") != "true" {
		logus.Debug("Skipping test data regenerative code")
		return nil
	}

	return callback()
}
