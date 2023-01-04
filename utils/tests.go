package utils

import (
	"fmt"
	"os"
)

func RegenerativeTest(callback func() error) error {
	if len(os.Getenv("DARK_TEST_REGENERATE")) == 0 {
		fmt.Println("Skipping test data regenerative code")
		return nil
	}

	return callback()
}
