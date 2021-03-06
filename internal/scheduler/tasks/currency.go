package tasks

import (
	"errors"
	"fmt"
	"github.com/lowitea/jeevez/internal/app"
	"github.com/lowitea/jeevez/internal/config"
	"github.com/lowitea/jeevez/internal/models"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	url2 "net/url"
	"regexp"
	"strconv"
)

// доступные валютные пары, по которым запрашиваются данные
var CurrencyPairs = [...]string{
	"USD_RUB",
	"USD_EUR",

	"RUB_USD",
	"RUB_EUR",

	"EUR_USD",
	"EUR_RUB",
}

func getCurrencyRate(url string) (float64, error) {
	var body []byte

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error get data: %s", err)
		return 0, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error get data: %s", err)
		return 0, err
	}

	tpl := regexp.MustCompile(`{"\w{3}_\w{3}":(\d+.\d+)}`)
	rateStr := tpl.FindStringSubmatch(string(body))[1]

	rate, err := strconv.ParseFloat(rateStr, 64)
	if err != nil {
		log.Printf("Error get data: %s", err)
		return 0, err
	}

	return rate, nil
}

// CurrencyTask таска обновляющая курсы валют в базе
func CurrencyTask(db *gorm.DB) {
	log.Printf("CurrencyTask has started")
	baseUrl := url2.URL{
		Scheme: config.CurrencyApiScheme,
		Host:   config.CurrencyApiHost,
		Path:   config.CurrencyApiPath,
	}

	for _, curPair := range CurrencyPairs {

		curUrl := baseUrl
		token := app.Config.CurrencyAPI.Token
		curUrl.RawQuery = fmt.Sprintf("q=%s&compact=ultra&apiKey=%s", curPair, token)

		curRate, err := getCurrencyRate(curUrl.String())
		if err != nil {
			log.Printf("getting currency rate error: %s", err)
			return
		}

		// получаем валюту из базы, или создаём новую, если не нашли
		var currencyRate models.CurrencyRate
		result := db.First(&currencyRate, "name = ?", curPair)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			currencyRate = models.CurrencyRate{Name: curPair}
			_ = db.Create(&currencyRate)
		} else if result.Error != nil {
			log.Printf("getting currency from db error: %s", result.Error)
			return
		}

		// обновляем данные в базе
		currencyRate.Value = curRate
		db.Save(&currencyRate)
	}
}
