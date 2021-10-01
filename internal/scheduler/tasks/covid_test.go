package tasks

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

// TestCovidStat_Update проверяет функцию обновления объекта статистики
func TestCovidStat_Update(t *testing.T) {
	stat := covidStat{
		Confirmed: 1, Deaths: 2, Recovered: 3, ConfirmedDiff: 4, DeathsDiff: 5,
		RecoveredDiff: 6, LastUpdate: "9.01.20", Active: 7, ActiveDiff: 8, FatalityRate: 0.5,
	}
	statNew := covidStat{
		Confirmed: 11, Deaths: 12, Recovered: 13, ConfirmedDiff: 14, DeathsDiff: 15,
		RecoveredDiff: 16, LastUpdate: "10.01.20", Active: 17, ActiveDiff: 18, FatalityRate: 0.8,
	}
	stat.Update(statNew)
	statExp := covidStat{
		Confirmed: 12, Deaths: 14, Recovered: 16, ConfirmedDiff: 18, DeathsDiff: 20,
		RecoveredDiff: 22, LastUpdate: "10.01.20", Active: 24, ActiveDiff: 26, FatalityRate: 0.65,
	}
	assert.Equal(t, statExp, stat)

	stat.FatalityRate = 0
	stat.Update(statNew)
	assert.Equal(t, 0.8, stat.FatalityRate)
}

type fakeBody struct {
	err     error
	content string
}

func (b fakeBody) Close() error {
	return nil
}

func (b fakeBody) Read(p []byte) (n int, err error) {
	c := []byte(b.content)
	copy(p, c)
	if b.err == nil {
		err = io.EOF
	} else {
		err = b.err
	}
	return len(c), err
}

// TestGetData тестирует функцию получения данных из апи
func TestGetData(t *testing.T) {
	defer func(f func(url string) (resp *http.Response, err error)) { httpGet = f }(httpGet)

	errExp := errors.New("test")
	httpGet = func(_ string) (*http.Response, error) { return nil, errExp }
	stat, err := getData("")
	assert.Nil(t, stat)
	assert.Equal(t, errExp, err)

	resp := http.Response{Body: fakeBody{err: errExp}}
	httpGet = func(_ string) (*http.Response, error) { return &resp, nil }
	stat, err = getData("")
	assert.Nil(t, stat)
	assert.Equal(t, errExp, err)

	resp = http.Response{Body: fakeBody{content: "invalid json"}}
	httpGet = func(_ string) (*http.Response, error) { return &resp, nil }
	stat, err = getData("")
	assert.Nil(t, stat)
	assert.IsType(t, &json.SyntaxError{}, err)

	resp = http.Response{Body: fakeBody{
		content: `{"data":[{"date":"2021-02-27","confirmed":976739,"deaths":15007,"recovered":895879,` +
			`"confirmed_diff":1825,"deaths_diff":41,"recovered_diff":2008,"last_update":"2021-02-28 05:22:20",` +
			`"active":65853,"active_diff":-224,"fatality_rate":0.0154,"region":{"iso":"RUS","name":"Russia",` +
			`"province":"Moscow","lat":"55.7504461","long":"37.6174943","cities":[]}}]}`,
	}}
	httpGet = func(_ string) (*http.Response, error) { return &resp, nil }
	stat, err = getData("")
	assert.Nil(t, err)
	expStat := []covidStat{
		{
			Confirmed:     976739,
			Deaths:        15007,
			Recovered:     895879,
			ConfirmedDiff: 1825,
			DeathsDiff:    41,
			RecoveredDiff: 2008,
			LastUpdate:    "2021-02-28 05:22:20",
			Active:        65853,
			ActiveDiff:    -224,
			FatalityRate:  0.0154,
		},
	}
	assert.Equal(t, expStat, stat)
}
