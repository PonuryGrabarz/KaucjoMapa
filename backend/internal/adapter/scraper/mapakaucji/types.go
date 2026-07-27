package mapakaucji

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

type OpeningHours struct {
	MonSun string `json:"mon-sun"`
	MonSat string `json:"mon-sat"`
	MonFri string `json:"mon-fri"`
	Sat    string `json:"sat"`
	Sun    any    `json:"sun"`
}

type MapPoint struct {
	Type     string `json:"type"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		StoreID                string       `json:"storeId"`
		Street                 string       `json:"street"`
		Town                   string       `json:"town"`
		Lat                    float64      `json:"lat"`
		Lon                    float64      `json:"lon"`
		SalesZoneDirector      string       `json:"salesZoneDirector"`
		SalesZoneDirectorEmail string       `json:"salesZoneDirectorEmail"`
		SalesZoneID            string       `json:"salesZoneId"`
		Active                 string       `json:"active"`
		SalesZoneDirectorID    string       `json:"salesZoneDirectorId"`
		LocationID             string       `json:"locationId"`
		TownID                 string       `json:"townId"`
		IsAgency               bool         `json:"isAgency"`
		IsVisible              bool         `json:"isVisible"`
		FirstOpeningDate       string       `json:"firstOpeningDate"`
		StoreURL               string       `json:"storeUrl"`
		RelativeStoreURL       string       `json:"relativeStoreUrl"`
		OpeningHours           OpeningHours `json:"openingHours"`
		LocatorNewMonth        bool         `json:"locatorNewMonth"`
		LocatorNewTwoWeeks     bool         `json:"locatorNewTwoWeeks"`
		LocatorMerrychef       bool         `json:"locatorMerrychef"`
	} `json:"properties"`
}
