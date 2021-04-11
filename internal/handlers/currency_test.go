package handlers

import (
	"fmt"
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
