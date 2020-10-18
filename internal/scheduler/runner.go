package scheduler

import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jasonlvhit/gocron"
	"io/ioutil"
	"log"
	"net/http"
)

type CovidStat struct {
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

type ApiResp struct {
	Data []CovidStat
}

func getData(url string) CovidStat {

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

	var data ApiResp
	_ = json.Unmarshal(body, &data)

	return data.Data[0]
}

func task(bot *tgbotapi.BotAPI) func() {
	return func() {
		log.Printf("Task has started")

		data := "2020-10-16"

		MoscowStat := getData(fmt.Sprintf("https://covid-api.com/api/reports?date=%s&iso=rus&region_province=Moscow", data))

		msg_text := fmt.Sprintf(
			"\U0001F9A0 <b>COVID-19 Статистика [Москва]</b>\n"+
				"%s\n\n"+
				"Подтверждённые: %d (+%d)\n"+
				"Смерти: %d (+%d)\n"+
				"Выздоровевшие: %d (+%d)\n"+
				"Болеющие: %d (+%d)\n"+
				"Летальность: %f\n",
			data,
			MoscowStat.Confirmed,
			MoscowStat.Confirmed_diff,
			MoscowStat.Deaths,
			MoscowStat.Deaths_diff,
			MoscowStat.Recovered,
			MoscowStat.Recovered_diff,
			MoscowStat.Active,
			MoscowStat.Active_diff,
			MoscowStat.Fatality_rate,
		)

		msg := tgbotapi.NewMessage(159096094, msg_text)
		msg.ParseMode = "HTML"
		msg.DisableNotification = true

		_, _ = bot.Send(msg)
	}
}

func Run(bot *tgbotapi.BotAPI) {
	log.Printf("Scheduler has started")
	s := gocron.NewScheduler()
	_ = s.Every(10).Seconds().Do(task(bot))
	//_ = s.Every(1).Day().At("03:26").Do(task(bot))
	<-s.Start()
}
