package base

import "darkbot/scrappy/shared/api"

func (s *BaseStorage) FixtureSetAPI(base_api api.APIinterface) {
	s.api = base_api
}
