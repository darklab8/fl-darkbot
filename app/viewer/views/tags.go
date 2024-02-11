package views

import (
	"strings"

	"github.com/darklab/fl-darkbot/app/settings/types"
)

func TagContains(name string, tags []types.Tag) bool {
	for _, tag := range tags {
		if strings.Contains(name, string(tag)) {
			return true
		}
	}
	return false
}
