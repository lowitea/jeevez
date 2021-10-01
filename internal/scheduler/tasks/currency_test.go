package tasks

import (
	"errors"
	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/tools/testTools"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
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
		return &http.Response{Body: fakeBody{content: `{"RUB_EUR":0.01183}`}}, nil
	}
	rate, err = getCurrencyRate("")
	assert.Equal(t, 0.01183, rate)
	assert.NoError(t, err)
}

// TestCurrencyTask тестирование таски на получение данных по валютам
func TestCurrencyTask(t *testing.T) {
	db := testTools.InitTestDB()
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
		return &http.Response{Body: fakeBody{content: `{"RUB_EUR":0.01183}`}}, nil
	}
	CurrencyTask(db)

	var rates []models.CurrencyRate
	db.Find(&rates)

	assert.Len(t, rates, 6)
	for _, r := range rates {
		assert.Equal(t, 0.01183, r.Value)
	}

	// проверяем обновление данных
	httpGet = func(_ string) (*http.Response, error) {
		return &http.Response{Body: fakeBody{content: `{"RUB_EUR":0.42}`}}, nil
	}
	CurrencyTask(db)

	db.Find(&rates)

	assert.Len(t, rates, 6)
	for _, r := range rates {
		assert.Equal(t, 0.42, r.Value)
	}

	// ломаем базу
	db.Exec("drop table currency_rates")
	assert.NotPanics(t, func() { CurrencyTask(db) })
}
