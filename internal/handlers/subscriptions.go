package handlers

import (
	"bytes"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/scheduler/subscriptions"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"strconv"
	"strings"
)

// cmdSubscriptions выводит список всех доступных подписок
func cmdSubscriptions(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *gorm.DB) {
	var msgTextB bytes.Buffer
	msgTextB.WriteString("Все доступные темы для подписки:\n\n")

	for subscr := range subscriptions.SubscriptionFuncMap {
		msgTextB.WriteString("<b>")
		msgTextB.WriteString(subscr.Name)
		msgTextB.WriteString("</b> - ")
		msgTextB.WriteString(subscr.Description)
		msgTextB.WriteString("\n")
	}

	msgTextB.WriteString("\nПример команды для подписки:\n/subscribe covid19-russia 11:00\n\n")

	var chat models.Chat
	db.First(&chat, "tg_id = ?", update.Message.Chat.ID)

	var chatSubscrs []models.ChatSubscription
	db.Find(&chatSubscrs, "chat_id = ?", chat.ID)

	if len(chatSubscrs) > 0 {
		msgTextB.WriteString("\nТемы на которые Вы подписаны:\n")
	}

	for _, chatSubscr := range chatSubscrs {
		var subscr models.Subscription
		if result := db.First(&subscr, chatSubscr.SubscriptionID); result.Error != nil {
			continue
		}
		msgTextB.WriteString("\n")
		msgTextB.WriteString("<b>")
		msgTextB.WriteString(subscr.Name)
		msgTextB.WriteString("</b> - ")
		msgTextB.WriteString(subscr.Description)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgTextB.String())
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = "HTML"
	_, _ = bot.Send(msg)
}

// parseTime парсит строку в формате 23:59 и возвращает секунды
func parseTime(timeStr string) (int64, error) {
	timeTokens := strings.Split(timeStr, ":")
	if len(timeTokens) != 2 {
		return 0, errors.New("tokens count error")
	}
	hoursStr, minutesStr := timeTokens[0], timeTokens[1]

	var hours int64
	var minutes int64
	var err error
	if hours, err = strconv.ParseInt(hoursStr, 10, 64); err != nil {
		return 0, errors.New("parse hours error")
	}
	if minutes, err = strconv.ParseInt(minutesStr, 10, 64); err != nil {
		return 0, errors.New("parse minutes error")
	}

	secs := hours*3600 + minutes*60

	// проверяем что временной диапазон в рамках допустимого
	if secs > 23*60*60+59*60 || secs < 0 {
		return 0, errors.New("time interval error")
	}

	return secs, nil
}

// getSubscription возвращает нужную подписку по имени
func getSubscription(name string) (models.Subscription, error) {
	for s := range subscriptions.SubscriptionFuncMap {
		if s.Name == name {
			return s, nil
		}
	}
	return models.Subscription{}, errors.New("subscription not found")
}

// cmdSubscribe выводит список всех доступных подписок
func cmdSubscribe(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *gorm.DB) {
	args := strings.Split(update.Message.Text, " ")

	if len(args) != 3 {
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"Чтобы подписаться, отправьте команду в формате:\n"+
				"/subscribe название_темы время_оповещения\n"+
				"Например, так:\n"+
				"/subscribe covid19-russia 11:00",
		)
		msg.ReplyToMessageID = update.Message.MessageID
		_, _ = bot.Send(msg)
		return
	}

	subscrName, subscrTime := args[1], args[2]

	// находим нужную подписку в мапе
	var subscr models.Subscription
	var err error

	// если не нашли, отправляем сообщение и выходим
	if subscr, err = getSubscription(subscrName); err != nil {
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"К сожалению, такой темы не существует(\n"+
				"Посмотреть доступные можно по команде /subscriptions",
		)
		msg.ReplyToMessageID = update.Message.MessageID
		_, _ = bot.Send(msg)
		return
	}

	// парсим введённое время
	subscrSeconds, err := parseTime(subscrTime)
	if err != nil {
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"Время должно быть не меньше чем 0:00 и меньше чем 24:00, без секунд.",
		)
		msg.ReplyToMessageID = update.Message.MessageID
		_, _ = bot.Send(msg)
		return
	}

	var chat models.Chat
	db.First(&chat, "tg_id = ?", update.Message.Chat.ID)

	chatSubscr := models.ChatSubscription{
		ChatID:         chat.ID,
		SubscriptionID: subscr.ID,
		Time:           subscrSeconds,
	}

	// Создаём объект связи чата с подпиской
	var msgText string
	clauses := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "chat_id"}, {Name: "subscription_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"time", "created_at"}),
	})
	if result := clauses.Create(&chatSubscr); result.Error != nil {
		log.Printf("create ChatSubscription error: %s", result.Error)
		msgText = msgText +
			"К сожалению, не получилось Вас подписать на тему," +
			"попробуйте пожалуйста позже ):"
	} else {
		msgText = msgText + fmt.Sprintf(
			"Я понял Вас :)\nБудет сделано.\n"+
				"Теперь я буду приходить и рассказывать вам новости по теме "+
				"<b>%s</b> каждый день в <b>%s</b>.", subscrName, subscrTime,
		)
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = "HTML"
	_, _ = bot.Send(msg)
}

// cmdUnsubscribe отписывает пользователя
func cmdUnsubscribe(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *gorm.DB) {
	args := strings.Split(update.Message.Text, " ")

	if len(args) != 2 {
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"Чтобы отписаться от темы, отправьте команду в формате:\n"+
				"/unsubscribe название_темы\n"+
				"Например, так:\n"+
				"/unsubscribe covid19-russia",
		)
		msg.ReplyToMessageID = update.Message.MessageID
		_, _ = bot.Send(msg)
		return
	}

	subscrName := args[1]

	var chat models.Chat
	db.First(&chat, "tg_id = ?", update.Message.Chat.ID)

	var subscr models.Subscription
	db.First(&subscr, "name = ?", subscrName)

	var chatSubscr models.ChatSubscription
	result := db.First(&chatSubscr, "chat_id = ? AND subscription_id = ?", chat.ID, subscr.ID)
	if result.Error != nil {
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			fmt.Sprintf("Не нашёл в своих записях информации, что Вы подписаны по тему %s :(", subscrName),
		)
		msg.ReplyToMessageID = update.Message.MessageID
		_, _ = bot.Send(msg)
		return
	}
	if result := db.Delete(&chatSubscr); result.Error != nil {
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"Произошёл пожар в картотеке, не смог откорректировать свои записи :(\n"+
				"Попробуйте, пожалуйста, позднее.",
		)
		msg.ReplyToMessageID = update.Message.MessageID
		_, _ = bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintf("Успешно отписал Вас от темы с именем %s\nНа здоровье)", subscrName),
	)
	msg.ReplyToMessageID = update.Message.MessageID
	_, _ = bot.Send(msg)
}

// cmdSubscription возвращает данные по подписке, без подписки
func cmdSubscription(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *gorm.DB) {
	args := strings.Split(update.Message.Text, " ")

	if len(args) != 2 {
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"Чтобы получить информацию по теме, отправьте команду в формате:\n"+
				"/subscription название_темы\n"+
				"Например, так:\n"+
				"/subscription covid19-russia",
		)
		msg.ReplyToMessageID = update.Message.MessageID
		_, _ = bot.Send(msg)
		return
	}

	subscrName := args[1]

	for subscr, sFunc := range subscriptions.SubscriptionFuncMap {
		if subscr.Name == subscrName {
			sFunc(bot, db, subscr)
			return
		}
	}

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"К сожалению, мне не удалось найти в своих записях такую тему :(",
	)
	msg.ReplyToMessageID = update.Message.MessageID
	_, _ = bot.Send(msg)
}

// BaseSubscriptionsHandler обработчик для команд подписок
func BaseSubscriptionsHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *gorm.DB) {
	if update.Message.Text == "/subscriptions" {
		cmdSubscriptions(update, bot, db)
	}
	if strings.HasPrefix(update.Message.Text, "/subscribe") {
		cmdSubscribe(update, bot, db)
	}
	if strings.HasPrefix(update.Message.Text, "/unsubscribe") {
		cmdUnsubscribe(update, bot, db)
	}
	if strings.HasPrefix(update.Message.Text, "/subscription") {
		cmdSubscription(update, bot, db)
	}
}
