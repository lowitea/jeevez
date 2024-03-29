package tasks

import (
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/lowitea/jeevez/internal/config"
	"github.com/lowitea/jeevez/internal/models"
	"gorm.io/gorm"
)

var currencies = [...]string{
	"USD", "RUB", "EUR", "GBP", "PLN", "TRY",
}

// currencyPairs доступные валютные пары, по которым запрашиваются данные
var currencyPairs = make([]string, 0, len(currencies)*(len(currencies)-1))

func init() {
	for _, firstCur := range currencies {
		for _, secCur := range currencies {
			if firstCur == secCur {
				continue
			}
			currencyPairs = append(currencyPairs, fmt.Sprintf("%s_%s", firstCur, secCur))
		}
	}
}

// getCurrencyRate получает валютные пары из апи
func getCurrencyRate(targetURL string) (float64, error) {
	var body []byte

	resp, err := httpGet(targetURL)
	if err != nil {
		log.Printf("Error get data: %s", err)
		return 0, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error get data: %s", err)
		return 0, err
	}

	bodyStr := string(body)

	tpl := regexp.MustCompile(`"conversion_rate":\s?(\d+\.\d+)`)
	parsedBody := tpl.FindStringSubmatch(bodyStr)

	if len(parsedBody) < 2 {
		log.Printf("Error parsed currency api.\nBody: %s", bodyStr)
		return 0, errors.New("error parsed currency api")
	}

	rateStr := parsedBody[1]
	rate, _ := strconv.ParseFloat(rateStr, 64)

	return rate, nil
}

// CurrencyTask таска обновляющая курсы валют в базе
func CurrencyTask(db *gorm.DB) {
	log.Printf("CurrencyTask has started")
	if db == nil {
		log.Printf("db is nil")
		return
	}
	baseURL := fmt.Sprintf(
		"%s://%s/%s",
		config.CurrencyAPIScheme, config.CurrencyAPIHost, fmt.Sprintf(
			config.CurrencyAPIPath, config.Cfg.CurrencyAPI.Token,
		),
	)

	for _, curPair := range currencyPairs {
		curURL := fmt.Sprintf("%s/%s", baseURL, strings.ReplaceAll(curPair, "_", "/"))
		curRate, err := getCurrencyRate(curURL)
		if err != nil {
			log.Printf("getting currency rate error: %s", err)
			return
		}

		// получаем валюту из базы, или создаём новую, если не нашли
		var currencyRate models.CurrencyRate
		result := db.First(&currencyRate, "name = ?", curPair)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			currencyRate = models.CurrencyRate{Name: curPair}
			db.Create(&currencyRate)
		} else if result.Error != nil {
			log.Printf("getting currency from db error: %s", result.Error)
			return
		}

		// обновляем данные в базе
		currencyRate.Value = curRate
		db.Save(&currencyRate)
	}
}
