package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/scheduler/subscriptions"
	"github.com/lowitea/jeevez/internal/structs"
	"strings"
)

func YogaCallbackHandler(update tgbotapi.Update, bot structs.Bot) {
	if !strings.HasPrefix(update.CallbackQuery.Data, "/yoga") {
		return
	}

	_, _ = bot.Request(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))

	nextMsg := "\n\nДля продолжения нажми: /yoga"

	var msg tgbotapi.MessageConfig
	switch {
	case update.CallbackQuery.Data == "/yoga valid":
		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Ура ответ, верный! 🎉 🎉 🎉"+nextMsg)
	case strings.HasPrefix(update.CallbackQuery.Data, "/yoga invalid"):
		validAnswer := strings.SplitN(update.CallbackQuery.Data, " ", 3)
		msg = tgbotapi.NewMessage(
			update.CallbackQuery.Message.Chat.ID,
			fmt.Sprintf("Эх, ошибка 🤷\nПравильный ответ:\n%s 🧘%s", validAnswer[2], nextMsg),
		)
	}
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	_, _ = bot.Send(msg)
}

func YogaHandler(update tgbotapi.Update, bot structs.Bot) {
	if update.Message.Text != "/yoga" {
		return
	}
	subscriptions.YogaTestTask(bot, nil, models.SubscrNameSubscrMap["yoga-test"], update.Message.Chat.ID)
}
