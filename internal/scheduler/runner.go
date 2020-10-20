package scheduler

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jasonlvhit/gocron"
	"github.com/lowitea/jeevez/internal/scheduler/tasks"
	"log"
	"time"
)

func Run(bot *tgbotapi.BotAPI) {
	log.Printf("Scheduler has started")
	s := gocron.NewScheduler()
	loc, _ := time.LoadLocation("Europe/Moscow")

	_ = s.Every(1).Day().At("12:30").Loc(loc).Do(tasks.CovidTask(bot))

	<-s.Start()
}
