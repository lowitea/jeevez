package handlers

import (
	"bytes"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/structs"
	"gorm.io/gorm"
	"strconv"
	"strings"
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
	msg.ReplyToMessageID = update.Message.MessageID
	_, _ = bot.Send(msg)
}

// DeleteChatHandler удаляет чат
func DeleteChatHandler(update tgbotapi.Update, bot structs.Bot, db *gorm.DB) {
	args := strings.Split(update.Message.Text, " ")

	if len(args) != 2 || args[0] != "/delUsr" {
		return
	}

	chatID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return
	}

	var msgText string
	if err := db.Delete(&models.Chat{}, "tg_id = ?", chatID).Error; err != nil {
		msgText = err.Error()
	} else {
		msgText = "done"
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	msg.ReplyToMessageID = update.Message.MessageID
	_, _ = bot.Send(msg)
}

// AdminHelpHandler показывает админские команды
func AdminHelpHandler(update tgbotapi.Update, bot structs.Bot) {
	if update.Message.Text != "/hh" {
		return
	}

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"/me - показать свой id и id чата\n"+
			"/ul - показать список пользователей\n"+
			"/delUsr # - удалить пользователя\n",
	)
	msg.ReplyToMessageID = update.Message.MessageID
	_, _ = bot.Send(msg)
}
