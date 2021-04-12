package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/tools/testTools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestGetCurPair тестирование функции определения валютной пары
func TestGetCurPair(t *testing.T) {
	cases := [...]struct {
		firstCur string
		secCur   string
		expPair  string
	}{
		{"доллар", "рубли", "USD_RUB"},
		{"доллара", "рубли", "USD_RUB"},
		{"долларов", "евро", "USD_EUR"},
		{"евро", "рубли", "EUR_RUB"},
		{"евро", "доллары", "EUR_USD"},
		{"рублей", "доллары", "RUB_USD"},
		{"рубль", "евро", "RUB_EUR"},
		{"рубля", "евро", "RUB_EUR"},
		{"доллар", "доллары", ""},
		{"рублей", "рубли", ""},
		{"евро", "евро", ""},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("firstCur=%s;secCur=%s", c.firstCur, c.secCur), func(t *testing.T) {
			assert.Equal(t, getCurPair(c.firstCur, c.secCur), c.expPair)
		})
	}
}

// TestGetCurrencyRate тестирование функции получения валюты из базы
func TestGetCurrencyRate(t *testing.T) {
	curName := "USD_EUR"

	db, _ := testTools.InitTestDB()
	db.Where(fmt.Sprintf("name = %s", curName)).Delete(&models.CurrencyRate{})

	rate, err := getCurrencyRate(db, "USD_EUR")

	assert.Errorf(t, err, "record not found")
	assert.Equal(t, 0.0, rate)

	expValue := 100500.42
	_ = db.Create(&models.CurrencyRate{
		Value: expValue,
		Name:  curName,
	})
	rate, err = getCurrencyRate(db, "USD_EUR")

	require.NoError(t, err)
	assert.Equal(t, expValue, rate)
}

// TestGetMsgAllCurrencies тестирует функцию формирования сообщения со всеми доступными валютами
func TestGetMsgAllCurrencies(t *testing.T) {
	db, _ := testTools.InitTestDB()

	msg, err := getMsgAllCurrencies(db)
	assert.Errorf(t, err, "none rates")
	assert.Equal(t, "", msg)

	rates := [...]models.CurrencyRate{
		{Value: 42, Name: "USD_EUR"},
		{Value: 100500, Name: "EUR_USD"},
	}

	db.Create(&rates)

	msg, err = getMsgAllCurrencies(db)
	assert.NoError(t, err)
	assert.Equal(
		t,
		"Курсы всех доступных валютных пар:\n\n"+
			"USD_EUR:    42.000000\n"+
			"EUR_USD:    100500.000000\n",
		msg,
	)
}

// TestCmdCurrencyRateAllRates получение списка всех валют
func TestCmdCurrencyRateAllRates(t *testing.T) {
	db, _ := testTools.InitTestDB()
	update := testTools.NewUpdate("/currency_rate")

	// сначала проверяем при пустой базе
	db.Exec("DELETE FROM currency_rates")
	expMsg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"Я прошу прощения. Биржа не отвечает по телефону. "+
			"Попробуйте уточнить у меня список валют позднее.",
	)
	botAPIMock := testTools.NewBotAPIMock(expMsg)

	cmdCurrencyRate(update, botAPIMock, db)

	botAPIMock.AssertExpectations(t)

	// теперь проверяем при наличии валют в базе
	db.Create(&models.CurrencyRate{Value: 100, Name: "RUB_USD"})
	expMsg = tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"Курсы всех доступных валютных пар:\n\n"+
			"RUB_USD:    100.000000\n",
	)
	botAPIMock = testTools.NewBotAPIMock(expMsg)

	cmdCurrencyRate(update, botAPIMock, db)

	botAPIMock.AssertExpectations(t)
}

// TestCmdCurrencyRateOneRate получение значение одной валютной пары
func TestCmdCurrencyRateOneRate(t *testing.T) {
	db, _ := testTools.InitTestDB()
	update := testTools.NewUpdate("/currency_rate RUB_USD")

	// пробуем получить несуществующую валюту
	db.Exec("DELETE FROM currency_rates")
	expMsg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"К сожалению, я не смог найти курс Вашей валюты. "+
			"Попробуйте проверить список доступных валют, повторив "+
			"эту команду без параметров.",
	)
	botAPIMock := testTools.NewBotAPIMock(expMsg)

	cmdCurrencyRate(update, botAPIMock, db)

	botAPIMock.AssertExpectations(t)

	// теперь получаем существующую
	db.Create(&models.CurrencyRate{Value: 42, Name: "RUB_USD"})
	expMsg = tgbotapi.NewMessage(update.Message.Chat.ID, "42.000000")
	botAPIMock = testTools.NewBotAPIMock(expMsg)

	cmdCurrencyRate(update, botAPIMock, db)

	botAPIMock.AssertExpectations(t)
}
