package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type covidStat struct {
	Confirmed      int
	Deaths         int
	Recovered      int
	Confirmed_diff int
	Deaths_diff    int
	Recovered_diff int
	Last_update    string
	Active         int
	Active_diff    int
	Fatality_rate  float64
}

func (stat *covidStat) Update(updStat covidStat) *covidStat {
	stat.Confirmed += updStat.Confirmed
	stat.Deaths += updStat.Deaths
	stat.Recovered += updStat.Recovered
	stat.Confirmed_diff += updStat.Confirmed_diff
	stat.Deaths_diff += updStat.Deaths_diff
	stat.Recovered_diff += updStat.Recovered_diff
	stat.Last_update = updStat.Last_update
	stat.Active += updStat.Active
	stat.Active_diff += updStat.Active_diff
	if stat.Fatality_rate != 0 && updStat.Fatality_rate != 0 {
		stat.Fatality_rate = (stat.Fatality_rate + updStat.Fatality_rate) / 2
	} else if stat.Fatality_rate == 0 {
		stat.Fatality_rate = updStat.Fatality_rate
	}
	return stat
}

func (stat *covidStat) GetMessage(data string, statName string) string {
	msgTemplate := "\U0001F9A0 <b>COVID-19 Статистика [%s]</b>\n" +
		"%s\n\n" +
		"Подтверждённые: %d (+%d)\n" +
		"Смерти: %d (+%d)\n" +
		"Выздоровевшие: %d (+%d)\n" +
		"Болеющие: %d (+%d)\n" +
		"Летальность: %f\n\n" +
		"https://yandex.ru/covid19/stat"
	return fmt.Sprintf(
		msgTemplate,
		statName,
		data,
		stat.Confirmed,
		stat.Confirmed_diff,
		stat.Deaths,
		stat.Deaths_diff,
		stat.Recovered,
		stat.Recovered_diff,
		stat.Active,
		stat.Active_diff,
		stat.Fatality_rate,
	)
}

type apiResp struct {
	Data []covidStat
}

func getData(url string) covidStat {

	resp, err := http.Get(url)
	if err != nil {
		print(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}

	fmt.Print(string(body))

	var data apiResp
	_ = json.Unmarshal(body, &data)

	var result covidStat

	for _, stat := range data.Data {
		result.Update(stat)
	}

	return result
}

func CovidTask(bot *tgbotapi.BotAPI) func() {
	return func() {
		log.Printf("Task has started")

		dt := time.Now().AddDate(0, 0, -1)
		data := dt.Format("2006-01-02")

		statUrls := map[string]string{
			"Москва": "https://covid-api.com/api/reports?date=%s&iso=rus&region_province=Moscow",
			"Россия": "https://covid-api.com/api/reports?date=%s&iso=rus",
		}

		for statName, statUrl := range statUrls {
			stat := getData(fmt.Sprintf(statUrl, data))
			msg := tgbotapi.NewMessage(159096094, stat.GetMessage(data, statName))
			msg.ParseMode = "HTML"
			msg.DisableNotification = true
			msg.DisableWebPagePreview = true

			_, _ = bot.Send(msg)
		}
	}
}
