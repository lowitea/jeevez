package scheduler

import (
	"github.com/go-co-op/gocron"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lowitea/jeevez/internal/scheduler/subscriptions"
	"github.com/lowitea/jeevez/internal/scheduler/tasks"
	"time"

	"gorm.io/gorm"
	"log"
)

func Run(bot *tgbotapi.BotAPI, db *gorm.DB) {
	log.Printf("Scheduler has started")

	loc, _ := time.LoadLocation("Europe/Moscow")
	s := gocron.NewScheduler(loc)

	// таска на обновление данных по ковид-19
	_, _ = s.Every(1).Day().At("10:00").Do(func() { tasks.CovidTask(db) })

	// таска на обновление курсов валют в базе
	_, _ = s.Every(1).Day().At("1:00").Do(func() { tasks.CurrencyTask(db) })

	// таска на рассылку подписок
	_, _ = s.Every(10).Minutes().Do(func() { subscriptions.Send(bot, db) })

	s.StartAsync()
}
