package handlers

import (
	"bytes"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/structs"
	"gorm.io/gorm"
)

// ChatListHandler показывает список пользователей
func ChatListHandler(update tgbotapi.Update, bot structs.Bot, db *gorm.DB) {
	if update.Message.Text != "/ul" {
		return
	}

	var chats []models.Chat
	db.Find(&chats)

	var msgTextB bytes.Buffer
	for _, c := range chats {
		msgTextB.WriteString(fmt.Sprintf("%d  - ", c.TgID))
		if c.TgTitle != "" {
			msgTextB.WriteString(fmt.Sprintf("%s ", c.TgTitle))
		}
		if c.TgName != "" {
			msgTextB.WriteString(fmt.Sprintf("@%s ", c.TgName))
		}
		msgTextB.WriteString(fmt.Sprintf("(%s)\n", c.TgType))
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgTextB.String())
	_, _ = bot.Send(msg)
}
