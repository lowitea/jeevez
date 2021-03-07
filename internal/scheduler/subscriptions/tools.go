package subscriptions

import (
	"errors"
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
	maxTime := minTime + 600

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

// InitSubscriptions создаёт в базе недостающие подписки
func InitSubscriptions(db *gorm.DB) error {
	log.Print("InitSubscriptions has started")
	for subscr := range SubscriptionFuncMap {
		// пытаемся получить подписку из базы по id и name
		subscrDB := models.Subscription{}
		result := db.First(&subscrDB, "id = ? AND name = ?", subscr.ID, subscr.Name)

		// если такого не нашлось, создаём
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// так как у нас захардкожены id в коде, нужно попробовать удалить из базы запись с таким id
			_ = db.Delete(&models.ChatSubscription{}, "subscription_id = ?", subscr.ID)
			_ = db.Delete(&models.Subscription{}, subscr.ID)

			// создаём новую запись
			if result = db.Create(&subscr); result.Error != nil {
				log.Printf("create Subscription error: %s", result.Error)
				return result.Error
			}
			return nil
		} else if result.Error != nil {
			log.Printf("update Subscription error: %s", result.Error)
			return result.Error
		}

		// обновляем запись если отличаются другие поля
		if subscr.Description != subscrDB.Description {
			db.Save(&subscr)
		}
	}
	return nil
}
