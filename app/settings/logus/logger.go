package logus

import (
	_ "github.com/darklab8/fl-darkbot/app/settings" // enverant.json injection to env
	"github.com/darklab8/go-utils/typelog"
)

var Log *typelog.Logger = typelog.NewLogger("darkbot")
