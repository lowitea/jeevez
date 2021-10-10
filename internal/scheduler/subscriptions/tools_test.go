package subscriptions

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/structs"
	"github.com/lowitea/jeevez/internal/tools/testtools"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
	"time"
)

// TestGetNowTimeInterval проверяет функцию получения временного интервала
func TestGetNowTimeInterval(t *testing.T) {
	cases := [...]struct {
		Time       string
		expMinTime int
		expMaxTime int
	}{
		{"15:04:05", 54000, 54590},
		{"11:00:00", 39600, 40190},
		{"20:10:00", 72600, 73190},
		{"0:00:00", 0, 590},
		{"23:59:59", 85800, 86390},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("time=%s", c.Time), func(t *testing.T) {
			testTime, _ := time.Parse(time.RFC3339, fmt.Sprintf("2021-06-06T%sZ", c.Time))
			actualMinTime, actualMaxTime := getNowTimeInterval(testTime)
			assert.Equal(t, c.expMinTime, actualMinTime)
			assert.Equal(t, c.expMaxTime, actualMaxTime)
		})
	}
}

// TestSend проверяет функцию отправки сообщений в соответствии с подписками
func TestSend(t *testing.T) {
	db := testtools.InitTestDB()

	defer func(m map[models.Subscription]func(
		bot structs.Bot, db *gorm.DB, subscr models.Subscription, chatTgId int64,
	)) {
		db.Exec("DELETE FROM chat_subscriptions")
		db.Exec("DELETE FROM chats")
		SubscriptionFuncMap = m
	}(SubscriptionFuncMap)

	subscr := models.SubscrNameSubscrMap["covid19-russia"]
	minTime, _ := getNowTimeInterval(time.Now())

	chat := models.Chat{TgID: 666}
	db.Create(&chat)
	db.Create(&models.ChatSubscription{
		ChatID:         chat.ID,
		SubscriptionID: subscr.ID,
		Time:           int64(minTime + 10),
		HumanTime:      "13:66",
	})

	expMsg := tgbotapi.NewMessage(
		chat.TgID,
		"🦠 <b>COVID-19 Статистика [Test Stat]</b>\n10-02-2010\n\nПодтверждённые: 100 (+110)\n"+
			"Смерти: 10001 (+1101)\nБолеющие: 1 (+101)\n"+
			"Летальность: 111.000000\n\nhttps://yandex.ru/covid19/stat",
	)
	expMsg.ParseMode = HTML
	expMsg.DisableWebPagePreview = true
	expMsg.DisableNotification = true
	botAPIMock := testtools.NewBotAPIMock(expMsg)

	// проверяет корректную отправку сообщения
	assert.NotPanics(t, func() { Send(botAPIMock, db) })

	// проверяем ошибку ненайденной функции в карте
	SubscriptionFuncMap = map[models.Subscription]TaskFunc{}
	assert.NotPanics(t, func() { Send(botAPIMock, db) })
}
