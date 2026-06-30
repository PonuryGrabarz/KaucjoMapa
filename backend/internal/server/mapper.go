package server

import "strings"

type PointStyle struct {
	ShortName string
	BgColor   string
	TextColor string
}

var chainStyles = map[string]PointStyle{
	"biedronka":          {ShortName: "Bd", BgColor: "#E3000B", TextColor: "#FFFFFF"},
	"lidl":               {ShortName: "Li", BgColor: "#0050AA", TextColor: "#FFD700"},
	"żabka":              {ShortName: "Żb", BgColor: "#00923F", TextColor: "#FFFFFF"},
	"zabka":              {ShortName: "Żb", BgColor: "#00923F", TextColor: "#FFFFFF"},
	"carrefour":          {ShortName: "Ca", BgColor: "#004B9B", TextColor: "#FFFFFF"},
	"aldi":               {ShortName: "Al", BgColor: "#1A3A6B", TextColor: "#FFFFFF"},
	"netto":              {ShortName: "Ne", BgColor: "#FFCC00", TextColor: "#CC0000"},
	"kaufland":           {ShortName: "Ka", BgColor: "#CC0000", TextColor: "#FFFFFF"},
	"lewiatan":           {ShortName: "Le", BgColor: "#C8102E", TextColor: "#FFFFFF"},
	"polo market":        {ShortName: "Po", BgColor: "#003DA5", TextColor: "#FFFFFF"},
	"polomarket":         {ShortName: "Po", BgColor: "#003DA5", TextColor: "#FFFFFF"},
	"spar":               {ShortName: "Sp", BgColor: "#006F3C", TextColor: "#FFFFFF"},
	"eurospar":           {ShortName: "ES", BgColor: "#006F3C", TextColor: "#FFCC00"},
	"stokrotka":          {ShortName: "Sk", BgColor: "#6AB023", TextColor: "#FFFFFF"},
	"dino":               {ShortName: "Di", BgColor: "#E30613", TextColor: "#FFFF00"},
	"intermarche":        {ShortName: "In", BgColor: "#E30613", TextColor: "#FFFFFF"},
	"intermarcheé":       {ShortName: "In", BgColor: "#E30613", TextColor: "#FFFFFF"},
	"e.leclerc":          {ShortName: "EL", BgColor: "#003189", TextColor: "#FFCC00"},
	"leclerc":            {ShortName: "EL", BgColor: "#003189", TextColor: "#FFCC00"},
	"delikatesy centrum": {ShortName: "DC", BgColor: "#C8102E", TextColor: "#FFFFFF"},
	"makro":              {ShortName: "Mk", BgColor: "#CC0000", TextColor: "#FFFFFF"},
	"selgros":            {ShortName: "Se", BgColor: "#F47920", TextColor: "#FFFFFF"},
	"auchan":             {ShortName: "Au", BgColor: "#E3000B", TextColor: "#FFFFFF"},
	"tesco":              {ShortName: "Te", BgColor: "#EE1C25", TextColor: "#FFFFFF"},
	"butelkomat":         {ShortName: "BT", BgColor: "#1B5E20", TextColor: "#FFFFFF"},
	"recyklomat":         {ShortName: "RC", BgColor: "#2E7D32", TextColor: "#FFFFFF"},
	"automat zwrotny":    {ShortName: "AZ", BgColor: "#2E7D32", TextColor: "#FFFFFF"},
}

var defaultStyle = PointStyle{
	ShortName: "Sk",
	BgColor:   "#546E7A",
	TextColor: "#FFFFFF",
}

func ResolvePointStyle(name string) PointStyle {
	normalized := strings.ToLower(strings.TrimSpace(name))

	for key, style := range chainStyles {
		if strings.Contains(normalized, key) {
			return style
		}
	}

	return defaultStyle
}
