package scheduler

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jasonlvhit/gocron"
	"github.com/lowitea/jeevez/internal/scheduler/subscriptions"
	"github.com/lowitea/jeevez/internal/scheduler/tasks"
	"time"

	"gorm.io/gorm"
	"log"
)

func Run(bot *tgbotapi.BotAPI, db *gorm.DB) {
	log.Printf("Scheduler has started")
	// TODO: Заменить на github.com/go-co-op/gocron
	s := gocron.NewScheduler()
	loc, _ := time.LoadLocation("Europe/Moscow")
	s.ChangeLoc(loc)

	// таска на рассылку данных по ковид-19
	_ = s.Every(1).Day().At("10:00").Do(func() { tasks.CovidTask(db) })

	// таска на обновление курсов валют в базе
	_ = s.Every(1).Day().At("1:00").Do(func() { tasks.CurrencyTask(db) })

	// таска на рассылку подписок
	_ = s.Every(10).Minutes().Do(func() { subscriptions.Send(bot, db) })

	<-s.Start()
}
