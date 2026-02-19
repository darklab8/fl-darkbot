package base

/*
Stuff with things to utilize darkstat api
*/

import "github.com/darklab8/fl-darkstat/darkstat/configs_export"

func NewBase2(base Base) *configs_export.PoB {
	return &configs_export.PoB{
		PoBCore: configs_export.PoBCore{
			Name:        base.Name,
			FactionName: &base.Affiliation,
			Health:      &base.Health,
		},
	}
}
