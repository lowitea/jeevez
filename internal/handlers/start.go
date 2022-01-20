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
	chat.TgID = update.Message.Chat.ID
	chat.TgName = update.Message.Chat.UserName
	chat.TgFirstName = update.Message.Chat.FirstName
	chat.TgLastName = update.Message.Chat.LastName
	chat.TgTitle = update.Message.Chat.Title
	chat.TgType = update.Message.Chat.Type
	clauses := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "tg_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"tg_name", "tg_first_name", "tg_last_name", "tg_title", "tg_type",
		}),
	})
	clauses.Create(&chat)
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
