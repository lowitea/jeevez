package tasks

import (
	"errors"
	"net/http"
	"testing"

	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/tools/testtools"
	"github.com/stretchr/testify/assert"
)

const (
	FakeRespBody = `
{
  "result": "success",
  "documentation": "https://www.exchangerate-api.com/docs",
  "terms_of_use": "https://www.exchangerate-api.com/terms",
  "time_last_update_unix": 1666224001,
  "time_last_update_utc": "Thu, 20 Oct 2022 00:00:01 +0000",
  "time_next_update_unix": 1666310401,
  "time_next_update_utc": "Fri, 21 Oct 2022 00:00:01 +0000",
  "base_code": "EUR",
  "target_code": "GBP",
  "conversion_rate": 0.8713
}
`

	ExpectRate = 0.8713
)

// TestGetCurrencyRate проверяет функцию получения валют
func TestGetCurrencyRate(t *testing.T) {
	defer func(f func(url string) (resp *http.Response, err error)) { httpGet = f }(httpGet)

	expErr := errors.New("err http get")
	httpGet = func(_ string) (*http.Response, error) { return nil, expErr }
	rate, err := getCurrencyRate("")
	assert.Equal(t, float64(0), rate)
	assert.Equal(t, expErr, err)

	expErr = errors.New("err read body")
	httpGet = func(_ string) (*http.Response, error) { return &http.Response{Body: fakeBody{err: expErr}}, nil }
	rate, err = getCurrencyRate("")
	assert.Equal(t, float64(0), rate)
	assert.Equal(t, expErr, err)

	httpGet = func(_ string) (*http.Response, error) { return &http.Response{Body: fakeBody{}}, nil }
	rate, err = getCurrencyRate("")
	assert.Equal(t, float64(0), rate)
	assert.Equal(t, errors.New("error parsed currency api"), err)

	httpGet = func(_ string) (*http.Response, error) {
		return &http.Response{Body: fakeBody{content: FakeRespBody}}, nil
	}
	rate, err = getCurrencyRate("")
	assert.Equal(t, ExpectRate, rate)
	assert.NoError(t, err)
}

// TestCurrencyTask тестирование таски на получение данных по валютам
func TestCurrencyTask(t *testing.T) {
	db := testtools.InitTestDB()
	db.Delete(&models.CurrencyRate{}, "true")
	defer func(f func(url string) (resp *http.Response, err error)) {
		httpGet = f
		db.Delete(&models.CurrencyRate{}, "true")
	}(httpGet)

	assert.NotPanics(t, func() { CurrencyTask(nil) })

	httpGet = func(_ string) (*http.Response, error) { return nil, errors.New("test") }
	assert.NotPanics(t, func() { CurrencyTask(db) })

	// проверяем на пустой базе
	httpGet = func(_ string) (*http.Response, error) {
		return &http.Response{Body: fakeBody{content: FakeRespBody}}, nil
	}
	CurrencyTask(db)

	var rates []models.CurrencyRate
	db.Find(&rates)

	expRatesCount := len(currencies) * (len(currencies) - 1)
	assert.Len(t, rates, expRatesCount)
	for _, r := range rates {
		assert.Equal(t, ExpectRate, r.Value)
	}

	// проверяем обновление данных
	httpGet = func(_ string) (*http.Response, error) {
		return &http.Response{Body: fakeBody{content: FakeRespBody}}, nil
	}
	CurrencyTask(db)

	db.Find(&rates)

	assert.Len(t, rates, expRatesCount)
	for _, r := range rates {
		assert.Equal(t, ExpectRate, r.Value)
	}

	// ломаем базу
	db.Exec("drop table currency_rates")
	assert.NotPanics(t, func() { CurrencyTask(db) })
}
