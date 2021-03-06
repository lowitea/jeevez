package tasks

type subscrConf struct {
	HumanName string
	UrlTpl    string
}

var subscrUrlMap = map[string]subscrConf{
	"covid19-moscow": {
		HumanName: "Москва",
		UrlTpl:    "https://covid-api.com/api/reports?date=%s&iso=rus&region_province=Moscow",
	},
	"covid19-russia": {
		HumanName: "Россия",
		UrlTpl:    "https://covid-api.com/api/reports?date=%s&iso=rus",
	},
}
