package scheduler

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jasonlvhit/gocron"
	"github.com/lowitea/jeevez/internal/scheduler/tasks"
	"log"
)

func Run(bot *tgbotapi.BotAPI) {
	log.Printf("Scheduler has started")
	s := gocron.NewScheduler()

	_ = s.Every(1).Day().At("10:00").Do(tasks.CovidTask(bot))

	<-s.Start()
}
