package tasks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
	"time"
)

// covidStat структура с ответом от апи статистки по ковиду
type covidStat struct {
	Confirmed     int     `json:"confirmed"`
	Deaths        int     `json:"deaths"`
	Recovered     int     `json:"recovered"`
	ConfirmedDiff int     `json:"confirmed_diff"`
	DeathsDiff    int     `json:"deaths_diff"`
	RecoveredDiff int     `json:"recovered_diff"`
	LastUpdate    string  `json:"last_update"`
	Active        int     `json:"active"`
	ActiveDiff    int     `json:"active_diff"`
	FatalityRate  float64 `json:"fatality_rate"`
}

// Update метод для сложения структур covidStat
func (stat *covidStat) Update(updStat covidStat) *covidStat {
	stat.Confirmed += updStat.Confirmed
	stat.ConfirmedDiff += updStat.ConfirmedDiff
	stat.Deaths += updStat.Deaths
	stat.DeathsDiff += updStat.DeathsDiff
	stat.Recovered += updStat.Recovered
	stat.RecoveredDiff += updStat.RecoveredDiff
	stat.LastUpdate = updStat.LastUpdate
	stat.Active += updStat.Active
	stat.ActiveDiff += updStat.ActiveDiff
	if stat.FatalityRate != 0 && updStat.FatalityRate != 0 {
		stat.FatalityRate = (stat.FatalityRate + updStat.FatalityRate) / 2
	} else if stat.FatalityRate == 0 {
		stat.FatalityRate = updStat.FatalityRate
	}
	return stat
}

// GetMessage вернуть строку для отправки в мессенджер
func (stat *covidStat) GetMessage(data string, statName string) string {
	type Ctx struct {
		Stat     *covidStat
		Data     string
		StatName string
	}

	msgTplString := "" +
		"\U0001F9A0 <b>COVID-19 Статистика [{{ .StatName }}]</b>\n" +
		"{{ .Data }}\n\n" +
		"Подтверждённые: {{ .Stat.Confirmed }} (+{{ .Stat.ConfirmedDiff }})\n" +
		"Смерти: {{ .Stat.Deaths }} (+{{ .Stat.DeathsDiff }})\n" +
		"Выздоровевшие: {{ .Stat.Recovered }} (+{{ .Stat.RecoveredDiff }})\n" +
		"Болеющие: {{ .Stat.Active }} (+{{ .Stat.ActiveDiff }})\n" +
		"Летальность: {{ printf \"%.6f\" .Stat.FatalityRate }}\n\n" +
		"https://yandex.ru/covid19/stat"

	msgTpl := template.Must(
		template.New("msgTpl").Parse(msgTplString))

	ctx := Ctx{stat, data, statName}

	msg := bytes.Buffer{}
	if err := msgTpl.Execute(&msg, ctx); err != nil {
		panic(err)
	}

	return msg.String()
}

// apiResp ответ от апи статистики
type apiResp struct {
	Data []covidStat
}

// получения статистики с апи
func getData(url string) (covidStat, error) {

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error get data: %s", err)
		return covidStat{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error get data: %s", err)
		return covidStat{}, err
	}

	var data apiResp
	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("Error get data: %s", err)
		return covidStat{}, err
	}

	var result covidStat

	for _, stat := range data.Data {
		result.Update(stat)
	}

	return result, nil
}

// CovidTask таска рассылающая статистику по ковиду
func CovidTask(bot *tgbotapi.BotAPI) func() {
	return func() {
		log.Printf("Covid task has started")

		dt := time.Now().AddDate(0, 0, -1)
		data := dt.Format("2006-01-02")

		statUrls := map[string]string{
			"Москва": "https://covid-api.com/api/reports?date=%s&iso=rus&region_province=Moscow",
			"Россия": "https://covid-api.com/api/reports?date=%s&iso=rus",
		}

		for statName, statUrl := range statUrls {
			stat, err := getData(fmt.Sprintf(statUrl, data))
			if err != nil {
				log.Printf("Error get data: %s", err)
				continue
			}
			msg := tgbotapi.NewMessage(159096094, stat.GetMessage(data, statName))
			msg.ParseMode = "HTML"
			msg.DisableNotification = true
			msg.DisableWebPagePreview = true

			_, _ = bot.Send(msg)
		}
	}
}
