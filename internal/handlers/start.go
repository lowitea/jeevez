package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
)

// StartHandler обрабатывает команду /start добавляет чатик в базу
func StartHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *gorm.DB) {
	if update.Message.Text != "/start" {
		return
	}

	msgText := "Приветствую! Я Ваш личный бот помощник. 🤵🏻\n"

	chat := models.Chat{TgID: update.Message.Chat.ID}
	if result := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&chat); result.Error != nil {
		log.Printf("create User error: %s", result.Error)
		msgText = msgText +
			"К сожалению, не получилось Вас зарегистрировать," +
			"попробуйте пожалуйста позже, с помощью команды /start ):"
	} else {
		msgText = msgText +
			"Готов помогать всем, чем умею. Чтобы узнать, подробнее, " +
			"предлагаю использовать команду /help :)"
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	_, _ = bot.Send(msg)
}
