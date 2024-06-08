package viewer

import (
	"testing"

	"github.com/darklab8/fl-darkbot/app/settings"
)

func TestDebugPerformance(t *testing.T) {
	if !settings.Env.IsDevEnv {
		return
	}

}
