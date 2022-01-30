package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lowitea/jeevez/internal/scheduler/subscriptions"
	"github.com/lowitea/jeevez/internal/structs"
)

func WeatherHandler(update *tgbotapi.Update, bot structs.Bot) {
	if update.Message.Text != "/weather" {
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, subscriptions.GetWeatherMessage("weather-moscow"))
	_, _ = bot.Send(msg)
}
