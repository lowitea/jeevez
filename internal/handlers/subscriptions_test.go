package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/tools/testTools"
	"testing"
)

// TestCmdSubscriptions проверяет команду возращающую список подписок
func TestCmdSubscriptions(t *testing.T) {
	db, _ := testTools.InitTestDB()
	db.Exec("DELETE FROM chat_subscriptions")
	db.Exec("DELETE FROM chats")

	hTime := "11:00"

	chat := models.Chat{TgID: 1}
	db.Create(&chat)
	for _, subscrData := range models.SubscrNameSubscrMap {
		db.Create(&models.ChatSubscription{
			ChatID:         chat.ID,
			SubscriptionID: subscrData.ID,
			HumanTime:      hTime,
			// делаем инкемент времени для проверки сортировки
			Time: subscrData.ID + 1000,
		})
	}

	update := testTools.NewUpdate("/subscriptions")
	expMsg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"Все доступные темы для подписки:"+
			"\n\n<b>covid19-russia</b> - Дневная статистика по COViD-19 по России"+
			"\n<b>covid19-moscow</b> - Дневная статистика по COViD-19 по Москве"+
			"\n\nПример команды для подписки:"+
			"\n/subscribe covid19-russia 11:00"+
			"\n\n\nТемы на которые Вы подписаны:"+
			"\n\n- <b>covid19-russia [11:00]</b> - Дневная статистика по COViD-19 по России"+
			"\n- <b>covid19-moscow [11:00]</b> - Дневная статистика по COViD-19 по Москве",
	)
	expMsg.ParseMode = "HTML"
	botAPIMock := testTools.NewBotAPIMock(expMsg)

	cmdSubscriptions(update, botAPIMock, db)

	botAPIMock.AssertExpectations(t)
}
