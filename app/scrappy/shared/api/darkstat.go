package api

import (
	"github.com/darklab8/fl-darkbot/app/settings"
	"github.com/darklab8/fl-darkstat/darkapis/darkhttp"
)

var DarkstatHttp = darkhttp.NewClient(settings.Env.DarkstatApiUrl)
