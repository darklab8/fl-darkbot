package views

import (
	"darkbot/app/settings/types"
	"strings"
)

func TagContains(name string, tags []types.Tag) bool {
	for _, tag := range tags {
		if strings.Contains(name, string(tag)) {
			return true
		}
	}
	return false
}
