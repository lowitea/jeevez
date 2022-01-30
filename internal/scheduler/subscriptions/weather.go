package subscriptions

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lowitea/jeevez/internal/config"
	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/structs"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var Cities = map[string]int{
	"weather-moscow": 524894,
}

var Icons = map[string]string{
	"01d": "☀️",
	"01n": "🌖",
	"02d": "🌤",
	"02n": "🌤",
	"03d": "☁️",
	"03n": "☁️",
	"04d": "☁️",
	"04n": "☁️",
	"09d": "🌨",
	"09n": "🌨",
	"10d": "🌦",
	"10n": "🌦",
	"11d": "⛈",
	"11n": "⛈",
	"13d": "❄️",
	"13n": "❄️",
	"50d": "🌫️",
	"50n": "🌫",
}

type weatherData struct {
	Weather []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"` // давление
		Humidity  int     `json:"humidity"` // влажность
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`  // направление
		Gust  float64 `json:"gust"` // порывы
	} `json:"wind"`
	Sys struct {
		Sunrise int64 `json:"sunrise"` // восход
		Sunset  int64 `json:"sunset"`  // закат
	} `json:"sys"`
}

type weatherCacheValue struct {
	data *weatherData
	ts   int64
}

var weatherCache = map[int]weatherCacheValue{}

const weatherURL = "http://api.openweathermap.org/data/2.5/weather?id=%d&lang=ru&units=metric&appid=%s"

// getWeatherAPIData получить текущий прогноз погоды
func getWeatherAPIData(cityID int) (*weatherData, error) {
	now := time.Now()

	if cache, ok := weatherCache[cityID]; ok {
		ttl := now.Add(-10 * time.Minute).Unix()
		if cache.ts > ttl {
			return cache.data, nil
		}
	}

	resp, err := http.Get(fmt.Sprintf(weatherURL, cityID, config.Cfg.WeatherAPI.Token)) // nolint
	if err != nil {
		log.Printf("Error get data: %s", err)
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error get data: %s", err)
		return nil, err
	}

	var data weatherData
	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("Error get data: %s", err)
		return nil, err
	}

	weatherCache[cityID] = weatherCacheValue{data: &data, ts: now.Unix()}

	return &data, nil
}

func GetWeatherMessage(city string) string {
	data, err := getWeatherAPIData(Cities[city])

	if err != nil {
		return "Я приношу свои извинения, но к сожалению метеостанция не отвечает 🤷"
	}

	return fmt.Sprintf(
		"🌡 %.0f° (ощущается %.0f°)\n"+
			"%s   %s\n\n"+
			"🗜 Давление: %d мм рт. ст.\n"+
			"💧 Влажность: %d%%\n\n"+
			"💨 Ветер: %.1f м/c, порывы до %.1f м/c\n\n"+
			"🌅 Восход: %s   🌇 Закат: %s\n",
		data.Main.Temp, data.Main.FeelsLike,
		Icons[data.Weather[0].Icon], strings.Title(data.Weather[0].Description),
		data.Main.Pressure, data.Main.Humidity,
		data.Wind.Speed, data.Wind.Gust,
		time.Unix(data.Sys.Sunrise, 0).Format("15:04:05"),
		time.Unix(data.Sys.Sunset, 0).Format("15:04:05"),
	)
}

func WeatherTask(bot structs.Bot, _ *gorm.DB, subscr models.Subscription, chatTgID int64) {
	msg := tgbotapi.NewMessage(chatTgID, GetWeatherMessage(subscr.Name))
	msg.DisableNotification = true
	msg.DisableWebPagePreview = true
	_, _ = bot.Send(msg)
}
