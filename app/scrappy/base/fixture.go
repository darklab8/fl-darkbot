package base

import "github.com/darklab8/fl-darkbot/app/scrappy/shared/api"

func FixtureBaseApiMock() BaseApi {
	return api.DarkstatHttp
}

func NewBaseApi() BaseApi {
	return api.DarkstatHttp
}
