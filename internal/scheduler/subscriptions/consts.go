package subscriptions

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/models"
	"gorm.io/gorm"
)

var SubscrNameSubscrMap = map[string]models.Subscription{
	"covid19-russia": {
		ID:          1,
		Name:        "covid19-russia",
		Description: "Дневная статистика по COViD-19 по России",
	},
	"covid19-moscow": {
		ID:          2,
		Name:        "covid19-moscow",
		Description: "Дневная статистика по COViD-19 по Москве",
	},
}

var SubscriptionFuncMap = map[models.Subscription]func(bot *tgbotapi.BotAPI, db *gorm.DB, subscr models.Subscription, chatTgId int64){
	{
		ID:          1,
		Name:        "covid19-russia",
		Description: "Дневная статистика по COViD-19 по России",
	}: CovidTask,
	{
		ID:          2,
		Name:        "covid19-moscow",
		Description: "Дневная статистика по COViD-19 по Москве",
	}: CovidTask,
}
