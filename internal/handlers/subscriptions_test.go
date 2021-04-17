package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/tools/testTools"
	"github.com/stretchr/testify/assert"
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
		"Все доступные темы для подписки:\n"+
			"\n<b>covid19-moscow</b> - Дневная статистика по COViD-19 по Москве"+
			"\n<b>covid19-russia</b> - Дневная статистика по COViD-19 по России"+
			"\n\nПример команды для подписки:"+
			"\n/subscribe covid19-russia 11:00"+
			"\n\n\nТемы на которые Вы подписаны:\n"+
			"\n- <b>covid19-russia [11:00]</b> - Дневная статистика по COViD-19 по России"+
			"\n- <b>covid19-moscow [11:00]</b> - Дневная статистика по COViD-19 по Москве",
	)
	expMsg.ParseMode = "HTML"
	botAPIMock := testTools.NewBotAPIMock(expMsg)

	cmdSubscriptions(update, botAPIMock, db)

	botAPIMock.AssertExpectations(t)
}

// TestParseTimeValid проверяет функцию парсинга времени на валидных кесах
func TestParseTimeValid(t *testing.T) {
	cases := []struct {
		hTime string
		exp   int64
	}{
		{"11:00", 39600},
		{"1:00", 3600},
		{"0:07", 420},
		{"13:59", 50340},
		{"00:00", 0},
		{"23:59", 86340},
		{"00:01", 60},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("time=%s", c.hTime), func(t *testing.T) {
			actualTime, err := parseTime(c.hTime)
			assert.NoError(t, err)
			assert.Equal(t, c.exp, actualTime)
		})
	}
}

// TestParseTimeInvalid проверяет функцию парсинга времени на невалидных кейсах
func TestParseTimeInvalid(t *testing.T) {
	cases := []struct {
		hTime     string
		expErrMsg string
	}{
		{"11:00:00", "tokens count error"},
		{"00:00:00", "tokens count error"},
		{"00:00:01", "tokens count error"},
		{"24:00", "time interval error"},
		{"24:05", "time interval error"},
		{"-01:05", "time interval error"},
		{"-1:00", "time interval error"},
		{"-1", "tokens count error"},
		{"1", "tokens count error"},
		{"a", "tokens count error"},
		{"1:v", "parse minutes error"},
		{"d:v", "parse hours error"},
		{"d:05", "parse hours error"},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("time=%s", c.hTime), func(t *testing.T) {
			actualTime, err := parseTime(c.hTime)
			assert.EqualError(t, err, c.expErrMsg)
			assert.Equal(t, int64(0), actualTime)
		})
	}
}

// TestCndSubscribeInvalid проверяем невалидные кейсы для команды подписки
func TestCndSubscribeInvalid(t *testing.T) {
	db, _ := testTools.InitTestDB()
	cases := [...]struct {
		Cmd       string
		MsgText   string
		ParseMode string
	}{
		{
			"/subscribe theme",
			"Чтобы подписаться, отправьте команду в формате:\n" +
				"/subscribe название_темы время_оповещения\n" +
				"Например, так:\n" +
				"/subscribe covid19-russia 11:00",
			"",
		},
		{
			"/subscribe no_exist_theme 11:00",
			"К сожалению, такой темы не существует(\n" +
				"Посмотреть доступные можно по команде /subscriptions",
			"",
		},
		{
			"/subscribe covid19-russia 25:00",
			"Время должно быть не меньше чем 0:00 и меньше чем 24:00, без секунд.",
			"",
		},
		{
			"/subscribe covid19-russia 11:00",
			"К сожалению, не получилось Вас подписать на тему, попробуйте пожалуйста позже ):",
			"HTML",
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("cmd=%s", c.Cmd), func(t *testing.T) {
			update := testTools.NewUpdate(c.Cmd)
			expMsg := tgbotapi.NewMessage(update.Message.Chat.ID, c.MsgText)
			expMsg.ParseMode = c.ParseMode
			botAPIMock := testTools.NewBotAPIMock(expMsg)
			cmdSubscribe(update, botAPIMock, db)
			botAPIMock.AssertExpectations(t)
		})
	}
}

// TestCndSubscribe проверяем подписку
func TestCndSubscribe(t *testing.T) {
	db, _ := testTools.InitTestDB()

	update := testTools.NewUpdate("/subscribe covid19-russia 11:00")

	chat := models.Chat{TgID: update.Message.Chat.ID}
	db.Create(&chat)

	// проверяем создание новой подписки
	expMsg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"Я понял Вас :)\nБудет сделано.\n"+
			"Теперь я буду приходить и рассказывать вам новости по теме "+
			"<b>covid19-russia</b> каждый день в <b>11:00</b>.",
	)
	expMsg.ParseMode = "HTML"

	botAPIMock := testTools.NewBotAPIMock(expMsg)
	cmdSubscribe(update, botAPIMock, db)
	botAPIMock.AssertExpectations(t)

	var chatSubscr models.ChatSubscription
	db.Last(&chatSubscr)

	assert.Equal(t, chat.ID, chatSubscr.ChatID)
	assert.Equal(t, int64(1), chatSubscr.SubscriptionID)
	assert.Equal(t, int64(39600), chatSubscr.Time)
	assert.Equal(t, "11:00", chatSubscr.HumanTime)

	// проверяем апдейт времени существующей подписки
	update.Message.Text = "/subscribe covid19-russia 23:42"
	expMsg.Text = "Я понял Вас :)\nБудет сделано.\n" +
		"Теперь я буду приходить и рассказывать вам новости по теме " +
		"<b>covid19-russia</b> каждый день в <b>23:42</b>."
	botAPIMock = testTools.NewBotAPIMock(expMsg)
	cmdSubscribe(update, botAPIMock, db)

	botAPIMock.AssertExpectations(t)

	db.Find(&chatSubscr, "chat_id = ? and subscription_id = ?", chat.ID, 1)

	assert.Equal(t, int64(85320), chatSubscr.Time)
	assert.Equal(t, "23:42", chatSubscr.HumanTime)
}
