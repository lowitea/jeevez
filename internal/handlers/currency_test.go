package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/tools/testtools"
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

	db := testtools.InitTestDB()
	db.Where("name = ?", curName).Delete(&models.CurrencyRate{})

	rate, err := getCurrencyRate(db, curName)
	assert.EqualError(t, err, "record not found")
	assert.Equal(t, 0.0, rate)

	expValue := 100500.42
	_ = db.Create(&models.CurrencyRate{
		Value: expValue,
		Name:  curName,
	})
	rate, err = getCurrencyRate(db, curName)

	require.NoError(t, err)
	assert.Equal(t, expValue, rate)
}

// TestGetCurrencyRateDBError проверяет ошибку в бд в функции getCurrencyRate
func TestGetCurrencyRateDBError(t *testing.T) {
	db := testtools.InitTestDB()
	db.Exec("DROP TABLE currency_rates")

	rate, err := getCurrencyRate(db, "USD_EUR")
	assert.EqualError(t, err, "ERROR: relation \"currency_rates\" does not exist (SQLSTATE 42P01)")
	assert.Equal(t, 0.0, rate)
}

// TestGetMsgAllCurrencies тестирует функцию формирования сообщения со всеми доступными валютами
func TestGetMsgAllCurrencies(t *testing.T) {
	db := testtools.InitTestDB()

	// проверяем работу с пустой базой
	db.Exec("DELETE FROM currency_rates")
	msg, err := getMsgAllCurrencies(db)
	assert.EqualError(t, err, "none rates")
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

// TestGetMsgAllCurrenciesDBError проверяет ошибку в бд в функции getMsgAllCurrencies
func TestGetMsgAllCurrenciesDBError(t *testing.T) {
	db := testtools.InitTestDB()
	db.Exec("DROP TABLE currency_rates")

	rate, err := getMsgAllCurrencies(db)
	assert.EqualError(t, err, "ERROR: relation \"currency_rates\" does not exist (SQLSTATE 42P01)")
	assert.Equal(t, "", rate)
}

// TestCmdCurrencyRateAllRates получение списка всех валют
func TestCmdCurrencyRateAllRates(t *testing.T) {
	db := testtools.InitTestDB()
	update := testtools.NewUpdate("/currency_rate")

	// сначала проверяем при пустой базе
	db.Exec("DELETE FROM currency_rates")
	expMsg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"Я прошу прощения. Биржа не отвечает по телефону. "+
			"Попробуйте уточнить у меня список валют позднее.",
	)
	botAPIMock := testtools.NewBotAPIMock(expMsg)

	cmdCurrencyRate(update, botAPIMock, db)

	botAPIMock.AssertExpectations(t)

	// теперь проверяем при наличии валют в базе
	db.Create(&models.CurrencyRate{Value: 100, Name: "RUB_USD"})
	expMsg = tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"Курсы всех доступных валютных пар:\n\n"+
			"RUB_USD:    100.000000\n",
	)
	botAPIMock = testtools.NewBotAPIMock(expMsg)

	cmdCurrencyRate(update, botAPIMock, db)

	botAPIMock.AssertExpectations(t)
}

// TestCmdCurrencyRateOneRate получение значение одной валютной пары
func TestCmdCurrencyRateOneRate(t *testing.T) {
	db := testtools.InitTestDB()
	update := testtools.NewUpdate("/currency_rate RUB_USD")

	// пробуем получить несуществующую валюту
	db.Exec("DELETE FROM currency_rates")
	expMsg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"К сожалению, я не смог найти курс Вашей валюты. "+
			"Попробуйте проверить список доступных валют, повторив "+
			"эту команду без параметров.",
	)
	botAPIMock := testtools.NewBotAPIMock(expMsg)

	cmdCurrencyRate(update, botAPIMock, db)

	botAPIMock.AssertExpectations(t)

	// теперь получаем существующую
	db.Create(&models.CurrencyRate{Value: 42, Name: "RUB_USD"})
	expMsg = tgbotapi.NewMessage(update.Message.Chat.ID, "42.000000")
	botAPIMock = testtools.NewBotAPIMock(expMsg)

	cmdCurrencyRate(update, botAPIMock, db)

	botAPIMock.AssertExpectations(t)
}

// TestCurrencyConverterHandler проверяет обработчик команд для валют
func TestCurrencyConverterHandler(t *testing.T) {
	db := testtools.InitTestDB()
	db.Exec("DELETE FROM currency_rates")
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
			update := testtools.NewUpdate(c.Cmd)
			expMsg := tgbotapi.NewMessage(update.Message.Chat.ID, c.MsgText)
			botAPIMock := testtools.NewBotAPIMock(expMsg)

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
			update := testtools.NewUpdate(cmd)
			botAPIMock := testtools.NewBotAPIMock(tgbotapi.MessageConfig{})

			CurrencyConverterHandler(update, botAPIMock, db)

			botAPIMock.AssertNotCalled(t, "Send")
		})
	}
}
