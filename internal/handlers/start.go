package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/structs"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
)

// UpdateChatInfoHandler обновляет информацию о чате
func UpdateChatInfoHandler(update tgbotapi.Update, db *gorm.DB) {
	var chat models.Chat
	db.First(&chat, "tg_id = ?", update.Message.Chat.ID)
	chat.TgName = update.Message.Chat.UserName
	chat.TgTitle = update.Message.Chat.Title
	chat.TgType = update.Message.Chat.Type
	db.Save(&chat)
}

// StartHandler обрабатывает команду /start добавляет чатик в базу
func StartHandler(update tgbotapi.Update, bot structs.Bot, db *gorm.DB) {
	if update.Message.Text != "/start" {
		return
	}

	msgText := "Приветствую! Я Ваш личный бот помощник. 🤵🏻\n"

	chat := models.Chat{TgID: update.Message.Chat.ID}
	if result := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&chat); result.Error != nil {
		log.Printf("create User error: %s", result.Error)
		msgText = msgText +
			"К сожалению, не получилось Вас зарегистрировать, " +
			"попробуйте пожалуйста позже, с помощью команды /start ):"
	} else {
		msgText = msgText +
			"Готов помогать всем, чем умею. Чтобы узнать, подробнее, " +
			"предлагаю использовать команду /help :)"
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	_, _ = bot.Send(msg)
}
