package utils

import (
	"fmt"
	"os"
)

func RegenerativeTest(callback func() error) error {
	if os.Getenv("DARK_TEST_REGENERATE") != "true" {
		fmt.Println("Skipping test data regenerative code")
		return nil
	}

	return callback()
}
