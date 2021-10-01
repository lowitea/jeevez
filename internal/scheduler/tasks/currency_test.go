package tasks

import (
	"errors"
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
