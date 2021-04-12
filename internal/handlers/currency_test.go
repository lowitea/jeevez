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

// TestCurrencyConverterHandler проверяет обработчик команд для валют
func TestCurrencyConverterHandler(t *testing.T) {
	db, _ := testTools.InitTestDB()
	db.Create(&[...]models.CurrencyRate{
		{Value: 77.425037, Name: "USD_RUB"},
		{Value: 0.840899, Name: "USD_EUR"},
		{Value: 0.012916, Name: "RUB_USD"},
		{Value: 0.010861, Name: "RUB_EUR"},
		{Value: 1.189850, Name: "EUR_USD"},
		{Value: 92.124168, Name: "EUR_RUB"},
	})

	cases := [...]struct {
		Cmd     string
		MsgText string
	}{
		{
			"/currency_rate",
			"Курсы всех доступных валютных пар:\n\n" +
				"USD_RUB:    77.425037\n" +
				"USD_EUR:    0.840899\n" +
				"RUB_USD:    0.012916\n" +
				"RUB_EUR:    0.010861\n" +
				"EUR_USD:    1.189850\n" +
				"EUR_RUB:    92.124168\n",
		},
		{"/currency_rate RUB_USD", "0.012916"},
		{"/currency_rate EUR_USD", "1.189850"},
		{"1000 долларов в рубли", "77425.04"},
		{"42 рубля в рубли", "42.00"},
		{"42 рубля в доллары", "0.54"},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("cmd=%s", c.Cmd), func(t *testing.T) {
			update := testTools.NewUpdate(c.Cmd)
			expMsg := tgbotapi.NewMessage(update.Message.Chat.ID, c.MsgText)
			botAPIMock := testTools.NewBotAPIMock(expMsg)

			CurrencyConverterHandler(update, botAPIMock, db)

			botAPIMock.AssertExpectations(t)
		})
	}

	// проверяем дополнительные невалидные кейсы
	db.Delete(models.CurrencyRate{}, "name = ?", "RUB_USD")
	badCases := [...]string{
		"невалидное сообщение",
		"42 рубля в доллары",
	}
	for _, cmd := range badCases {
		t.Run(fmt.Sprintf("cmd=%s", cmd), func(t *testing.T) {
			update := testTools.NewUpdate(cmd)
			botAPIMock := testTools.NewBotAPIMock(tgbotapi.MessageConfig{})

			CurrencyConverterHandler(update, botAPIMock, db)

			botAPIMock.AssertNotCalled(t, "Send")
		})
	}
}
