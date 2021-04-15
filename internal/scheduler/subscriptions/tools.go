package subscriptions

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/models"
	"gorm.io/gorm"
	"log"
	"time"
)

// Send отправляет сообщения по всем подпискам
func Send(bot *tgbotapi.BotAPI, db *gorm.DB) {
	log.Printf("Start send subscriptions")

	loc, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(loc)

	roundedMinMinutes := now.Minute() / 10 * 10
	minTime := now.Hour()*3600 + roundedMinMinutes*60
	maxTime := minTime + 590

	var chatSubscriptions []models.ChatSubscription

	db.Where("time BETWEEN ? AND ?", minTime, maxTime).Find(&chatSubscriptions)

	for _, chatSubscr := range chatSubscriptions {
		var subscr models.Subscription
		if result := db.First(&subscr, chatSubscr.SubscriptionID); result.Error != nil {
			log.Printf("getting Subscription error: %s", result.Error)
			continue
		}

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
