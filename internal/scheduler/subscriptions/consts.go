package subscriptions

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/models"
)

var SubscriptionFuncMap = map[models.Subscription]func(bot *tgbotapi.BotAPI){
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
