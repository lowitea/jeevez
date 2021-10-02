package tasks

import "net/http"

type subscrConf struct {
	HumanName string
	URLTpl    string
}

var subscrURLMap = map[string]subscrConf{
	"covid19-moscow": {
		HumanName: "Москва",
		URLTpl:    "https://covid-api.com/api/reports?date=%s&iso=rus&region_province=Moscow",
	},
	"covid19-russia": {
		HumanName: "Россия",
		URLTpl:    "https://covid-api.com/api/reports?date=%s&iso=rus",
	},
}

var httpGet = http.Get
