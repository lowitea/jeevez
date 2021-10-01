package subscriptions

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/tools/testTools"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testCovidStat = models.CovidStat{
	SubscriptionName: "covid19-russia",
	HumanName:        "Test Stat",
	Confirmed:        100,
	Deaths:           10001,
	Recovered:        11,
	ConfirmedDiff:    110,
	DeathsDiff:       1101,
	RecoveredDiff:    10101,
	LastUpdate:       "10-02-2010",
	Active:           1,
	ActiveDiff:       101,
	FatalityRate:     111,
}

const testCovidMsg = "\U0001F9A0 <b>COVID-19 Статистика [Test Stat]</b>\n" +
	"10-02-2010\n\n" +
	"Подтверждённые: 100 (+110)\n" +
	"Смерти: 10001 (+1101)\n" +
	"Выздоровевшие: 11 (+10101)\n" +
	"Болеющие: 1 (+101)\n" +
	"Летальность: 111.000000\n\n" +
	"https://yandex.ru/covid19/stat"

// TestGetMessage тестирует функцию формирования сообщения из объекта статистики
func TestGetMessage(t *testing.T) {
	actualMsg := GetMessage(testCovidStat)
	assert.Equal(t, testCovidMsg, actualMsg)
}

// TestCovidTask проверяет таску отправки сообщения со статистикой ковида
func TestCovidTask(t *testing.T) {
	db := testTools.InitTestDB()
	db.Exec("DELETE FROM covid_stats")

	var chatId int64 = 1010
	subscr := models.SubscrNameSubscrMap[testCovidStat.SubscriptionName]

	botAPIMock := testTools.NewBotAPIMock(tgbotapi.MessageConfig{})
	CovidTask(botAPIMock, db, subscr, chatId)
	botAPIMock.AssertNotCalled(t, "Send")

	// создаём объект статистики
	db.Create(&testCovidStat)

	expMsg := tgbotapi.NewMessage(chatId, testCovidMsg)
	expMsg.ParseMode = "HTML"
	expMsg.DisableNotification = true
	expMsg.DisableWebPagePreview = true
	botAPIMock = testTools.NewBotAPIMock(expMsg)
	CovidTask(botAPIMock, db, subscr, chatId)
	botAPIMock.AssertExpectations(t)
}
