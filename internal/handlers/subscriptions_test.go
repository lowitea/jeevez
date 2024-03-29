package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/tools/testtools"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

// TestCmdSubscriptions проверяет команду возращающую список подписок
func TestCmdSubscriptions(t *testing.T) {
	db := testtools.InitTestDB()

	cleanDB := func() {
		db.Delete(&models.ChatSubscription{}, "true")
		db.Delete(&models.Chat{}, "true")
	}

	cleanDB()
	defer cleanDB()

	hTime := "11:00"

	chat := models.Chat{TgID: 1}
	db.Create(&chat)

	for _, subscrData := range models.SubscrNameSubscrMap {
		db.Create(&models.ChatSubscription{
			ChatID:         chat.ID,
			SubscriptionID: subscrData.ID,
			HumanTime:      hTime,
			// делаем инкремент времени для проверки сортировки
			Time: subscrData.ID + 1000,
		})
	}

	update := testtools.NewUpdate("/subscriptions")
	expMsg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"Все доступные темы для подписки:\n"+
			"\n<b>covid19-moscow</b> - Дневная статистика по COViD-19 по Москве"+
			"\n<b>covid19-russia</b> - Дневная статистика по COViD-19 по России"+
			"\n<b>weather-moscow</b> - Текущая погода в Москве"+
			"\n<b>yoga-test</b> - Ежедневный тест с позами йоги"+
			"\n\nПример команды для подписки:"+
			"\n/subscribe covid19-russia 11:00"+
			"\n\n\nТемы на которые Вы подписаны:\n"+
			"\n- <b>covid19-russia [11:00]</b> - Дневная статистика по COViD-19 по России"+
			"\n- <b>covid19-moscow [11:00]</b> - Дневная статистика по COViD-19 по Москве"+
			"\n- <b>yoga-test [11:00]</b> - Ежедневный тест с позами йоги"+
			"\n- <b>weather-moscow [11:00]</b> - Текущая погода в Москве",
	)
	expMsg.ParseMode = HTML
	botAPIMock := testtools.NewBotAPIMock(expMsg)

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
	db := testtools.InitTestDB()
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
			HTML,
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("cmd=%s", c.Cmd), func(t *testing.T) {
			update := testtools.NewUpdate(c.Cmd)
			update.Message.Chat.ID = 666
			expMsg := tgbotapi.NewMessage(update.Message.Chat.ID, c.MsgText)
			expMsg.ParseMode = c.ParseMode
			botAPIMock := testtools.NewBotAPIMock(expMsg)
			cmdSubscribe(update, botAPIMock, db)
			botAPIMock.AssertExpectations(t)
		})
	}
}

// TestCmdSubscribe проверяем подписку
func TestCmdSubscribe(t *testing.T) {
	db := testtools.InitTestDB()

	update := testtools.NewUpdate("/subscribe covid19-russia 11:00")
	update.Message.Chat.ID = 777

	chat := models.Chat{TgID: update.Message.Chat.ID}
	db.Create(&chat)

	// проверяем создание новой подписки
	expMsg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"Я понял Вас :)\nБудет сделано.\n"+
			"Теперь я буду приходить и рассказывать вам новости по теме "+
			"<b>covid19-russia</b> каждый день в <b>11:00</b>.",
	)
	expMsg.ParseMode = HTML

	botAPIMock := testtools.NewBotAPIMock(expMsg)
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
	botAPIMock = testtools.NewBotAPIMock(expMsg)
	cmdSubscribe(update, botAPIMock, db)

	botAPIMock.AssertExpectations(t)

	db.Find(&chatSubscr, "chat_id = ? and subscription_id = ?", chat.ID, 1)

	assert.Equal(t, int64(85320), chatSubscr.Time)
	assert.Equal(t, "23:42", chatSubscr.HumanTime)
}

// TestCmdUnsubscribeInvalid проверяет невалидные кейсы для команды отписки
func TestCmdUnsubscribeInvalid(t *testing.T) {
	db := testtools.InitTestDB()
	cases := [...]struct {
		Cmd       string
		MsgText   string
		ParseMode string
	}{
		{
			"/unsubscribe",
			"Чтобы отписаться от темы, отправьте команду в формате:\n" +
				"/unsubscribe название_темы\n" +
				"Например, так:\n" +
				"/unsubscribe covid19-russia",
			"",
		},
		{
			"/unsubscribe no_exist_theme",
			"Не нашёл в своих записях информации, что Вы подписаны по тему <b>no_exist_theme</b> :(",
			HTML,
		},
		{
			"/unsubscribe covid19-russia",
			"Не нашёл в своих записях информации, что Вы подписаны по тему <b>covid19-russia</b> :(",
			HTML,
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("cmd=%s", c.Cmd), func(t *testing.T) {
			update := testtools.NewUpdate(c.Cmd)
			update.Message.Chat.ID = 666
			expMsg := tgbotapi.NewMessage(update.Message.Chat.ID, c.MsgText)
			expMsg.ParseMode = c.ParseMode
			botAPIMock := testtools.NewBotAPIMock(expMsg)
			cmdUnsubscribe(update, botAPIMock, db)
			botAPIMock.AssertExpectations(t)
		})
	}
}

// TestCmdUnsubscribe проверяет отписку от темы
func TestCmdUnsubscribe(t *testing.T) {
	db := testtools.InitTestDB()

	update := testtools.NewUpdate("/unsubscribe covid19-russia")

	chat := models.Chat{TgID: update.Message.Chat.ID}
	db.Create(&chat)
	chatSubscr := models.ChatSubscription{ChatID: chat.ID, SubscriptionID: 1}
	db.Create(&chatSubscr)

	expMsg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"Успешно отписал Вас от темы с именем <b>covid19-russia</b>\nНа здоровье)",
	)
	expMsg.ParseMode = HTML
	botAPIMock := testtools.NewBotAPIMock(expMsg)
	cmdUnsubscribe(update, botAPIMock, db)
	botAPIMock.AssertExpectations(t)

	// проверяем что подписка реально удалилась из базы
	result := db.First(
		&models.ChatSubscription{},
		"chat_id = ? AND subscription_id = ?",
		update.Message.Chat.ID,
		1,
	)

	assert.EqualError(t, result.Error, gorm.ErrRecordNotFound.Error())
}

// TestCmdSubscriptionInvalid проверяет невалидные кейсы для получения данных темы без подписки
func TestCmdSubscriptionInvalid(t *testing.T) {
	db := testtools.InitTestDB()
	cases := [...]struct {
		Cmd     string
		MsgText string
	}{
		{
			"/subscription",
			"Чтобы получить информацию по теме, отправьте команду в формате:\n" +
				"/subscription название_темы\n" +
				"Например, так:\n" +
				"/subscription covid19-russia",
		},
		{
			"/subscription no_exist_theme",
			"К сожалению, мне не удалось найти в своих записях такую тему :(",
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("cmd=%s", c.Cmd), func(t *testing.T) {
			update := testtools.NewUpdate(c.Cmd)
			expMsg := tgbotapi.NewMessage(update.Message.Chat.ID, c.MsgText)
			botAPIMock := testtools.NewBotAPIMock(expMsg)
			cmdSubscription(update, botAPIMock, db)
			botAPIMock.AssertExpectations(t)
		})
	}
}

// TestCmdSubscription проверяет команду получения данных темы без подписки
func TestCmdSubscription(t *testing.T) {
	db := testtools.InitTestDB()
	covidStat := models.CovidStat{
		SubscriptionName: "covid19-moscow",
		Confirmed:        10,
		Deaths:           101,
		ConfirmedDiff:    23,
		DeathsDiff:       32,
		LastUpdate:       "2021-04-18 04:20:41",
		Active:           45,
		ActiveDiff:       54,
		FatalityRate:     99.9,
	}
	db.Create(&covidStat)

	update := testtools.NewUpdate("/subscription covid19-moscow")
	expMsg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"🦠 <b>COVID-19 Статистика []</b>\n"+
			"2021-04-18 04:20:41\n\n"+
			"Подтверждённые: 10 (+23)\n"+
			"Смерти: 101 (+32)\n"+
			"Болеющие: 45 (+54)\n"+
			"Летальность: 99.900000\n\n"+
			"https://yandex.ru/covid19/stat",
	)
	expMsg.ParseMode = HTML
	expMsg.DisableWebPagePreview = true
	expMsg.DisableNotification = true
	botAPIMock := testtools.NewBotAPIMock(expMsg)
	cmdSubscription(update, botAPIMock, db)
	botAPIMock.AssertExpectations(t)
}

// TestSubscriptionsHandler тестирует общий обработчик для команд управления подписками
func TestSubscriptionsHandler(t *testing.T) {
	db := testtools.InitTestDB()
	cases := [...]struct {
		Cmd       string
		MsgText   string
		ParseMode string
	}{
		{
			"/subscriptions",
			"Все доступные темы для подписки:\n\n" +
				"<b>covid19-moscow</b> - Дневная статистика по COViD-19 по Москве\n" +
				"<b>covid19-russia</b> - Дневная статистика по COViD-19 по России\n" +
				"<b>weather-moscow</b> - Текущая погода в Москве\n" +
				"<b>yoga-test</b> - Ежедневный тест с позами йоги\n\n" +
				"Пример команды для подписки:\n" +
				"/subscribe covid19-russia 11:00\n\n",
			HTML,
		},
		{
			"/subscribe",
			"Чтобы подписаться, отправьте команду в формате:\n" +
				"/subscribe название_темы время_оповещения\n" +
				"Например, так:\n" +
				"/subscribe covid19-russia 11:00",
			"",
		},
		{
			"/unsubscribe",
			"Чтобы отписаться от темы, отправьте команду в формате:\n" +
				"/unsubscribe название_темы\n" +
				"Например, так:\n" +
				"/unsubscribe covid19-russia",
			"",
		},
		{
			"/subscription",
			"Чтобы получить информацию по теме, отправьте команду в формате:\n" +
				"/subscription название_темы\n" +
				"Например, так:\n" +
				"/subscription covid19-russia",
			"",
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("cmd=%s", c.Cmd), func(t *testing.T) {
			update := testtools.NewUpdate(c.Cmd)
			update.Message.Chat.ID = 666
			expMsg := tgbotapi.NewMessage(update.Message.Chat.ID, c.MsgText)
			expMsg.ParseMode = c.ParseMode
			botAPIMock := testtools.NewBotAPIMock(expMsg)
			SubscriptionsHandler(update, botAPIMock, db)
			botAPIMock.AssertExpectations(t)
		})
	}
}
