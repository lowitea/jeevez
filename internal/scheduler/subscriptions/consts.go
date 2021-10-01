package subscriptions

import (
	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/structs"
	"gorm.io/gorm"
)

var SubscriptionFuncMap = map[models.Subscription]func(bot structs.Bot, db *gorm.DB, subscr models.Subscription, chatTgId int64){
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
