package viewer

import (
	"darkbot/app/settings/utils"
	"testing"
)

func TestDebugPerformance(t *testing.T) {
	if !utils.FixtureDevEnv() {
		return
	}

}
