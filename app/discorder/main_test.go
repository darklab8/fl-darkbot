package discorder

import (
	"testing"

	"github.com/darklab8/fl-darkbot/app/settings/logus"
)

func TestLogging(t *testing.T) {

	long_test := "long_test_111111111111111111"
	short_test := "short"

	logus.Log.Warn("send long msg", logus.MsgContent(long_test))
	logus.Log.Warn("send shot msg", logus.MsgContent(short_test))
}
