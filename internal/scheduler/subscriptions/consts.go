package subscriptions

import (
	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/structs"
	"gorm.io/gorm"
)

const HTML = "HTML"

type TaskFunc = func(bot structs.Bot, db *gorm.DB, subscr models.Subscription, chatTgId int64)

var SubscriptionFuncMap = map[models.Subscription]TaskFunc{
	models.SubscrNameSubscrMap["covid19-russia"]: CovidTask,
	models.SubscrNameSubscrMap["covid19-moscow"]: CovidTask,
	models.SubscrNameSubscrMap["yoga-test"]:      YogaTestTask,
}
