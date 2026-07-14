package kaucjomatapi

type Kaucjomat struct {
	Brand    string  `json:"brand"`
	City     string  `json:"city"`
	Address  string  `json:"address"`
	Postcode string  `json:"postcode"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	Pet      bool    `json:"pet"`
	Can      bool    `json:"can"`
	Glass    bool    `json:"glass"`
}
