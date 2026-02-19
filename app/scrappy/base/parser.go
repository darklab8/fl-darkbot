package base

type baseSerializer struct {
	Affiliation string  `json:"affiliation"`
	Health      float64 `json:"health"`
	Tid         int     `json:"tid"`
}
