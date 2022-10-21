package subscriptions

import (
	"log"
	"time"

	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/structs"
	"gorm.io/gorm"
)

// getNowTimeInterval получить округлённый временной интервал на основе текущего времени.
// Например, если сейчас 17:23, вернётся 62400 и 62990,
// это 17:20:00 и 17:29:50 переведённое в минуты.
func getNowTimeInterval(now time.Time) (minTime int, maxTime int) {
	roundedMinMinutes := now.Minute() / 10 * 10
	minTime = now.Hour()*3600 + roundedMinMinutes*60
	maxTime = minTime + 590
	return
}

// Send отправляет сообщения по всем подпискам
func Send(bot structs.Bot, db *gorm.DB) {
	log.Printf("Start send subscriptions")

	loc, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(loc)
	minTime, maxTime := getNowTimeInterval(now)

	var chatSubscriptions []models.ChatSubscription
	db.Where("time BETWEEN ? AND ?", minTime, maxTime).Find(&chatSubscriptions)

	for _, chatSubscr := range chatSubscriptions {
		var subscr models.Subscription
		db.First(&subscr, chatSubscr.SubscriptionID)

		sFunc, ok := SubscriptionFuncMap[subscr]
		if !ok {
			log.Print("Subscription func not found error")
			continue
		}

		var chat models.Chat
		db.First(&chat, chatSubscr.ChatID)

		sFunc(bot, db, subscr, chat.TgID)
	}
}
