package tasks

import (
	"encoding/json"
	"errors"
	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/tools/testTools"
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

var validResp = http.Response{Body: fakeBody{
	content: `{"data":[{"date":"2021-02-27","confirmed":976739,"deaths":15007,"recovered":895879,` +
		`"confirmed_diff":1825,"deaths_diff":41,"recovered_diff":2008,"last_update":"2021-02-28 05:22:20",` +
		`"active":65853,"active_diff":-224,"fatality_rate":0.0154,"region":{"iso":"RUS","name":"Russia",` +
		`"province":"Moscow","lat":"55.7504461","long":"37.6174943","cities":[]}}]}`,
}}

var validStat = covidStat{
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

	httpGet = func(_ string) (*http.Response, error) { return &validResp, nil }
	stat, err = getData("")
	assert.NoError(t, err)
	expStat := []covidStat{validStat}
	assert.Equal(t, expStat, stat)
}

// TestGetStat тестирует функцию получения статистики
func TestGetStat(t *testing.T) {
	defer func(f func(url string) (resp *http.Response, err error)) { httpGet = f }(httpGet)

	httpGet = func(_ string) (*http.Response, error) { return nil, errors.New("test") }
	stat, err := getStat("")
	assert.Nil(t, stat)
	assert.EqualError(t, err, "getting covid stats error: test")

	resp := http.Response{Body: fakeBody{content: `{"data":[]}`}}
	httpGet = func(_ string) (*http.Response, error) { return &resp, nil }
	stat, err = getStat("")
	assert.Nil(t, stat)
	assert.EqualError(t, err, "getting covid stats error")

	httpGet = func(_ string) (*http.Response, error) { return &validResp, nil }
	stat, err = getStat("")
	assert.Equal(t, &validStat, stat)
	assert.NoError(t, err)
}

// TestCovidTask тестирование таски на получение данных по ковиду
func TestCovidTask(t *testing.T) {
	db := testTools.InitTestDB()
	db.Delete(&models.CovidStat{}, "true")
	defer func(f func(url string) (resp *http.Response, err error)) {
		httpGet = f
		db.Delete(&models.CovidStat{}, "true")
	}(httpGet)

	assert.NotPanics(t, func() { CovidTask(nil) })

	httpGet = func(_ string) (*http.Response, error) { return nil, errors.New("test") }
	assert.NotPanics(t, func() { CovidTask(db) })

	// проверяем запись в пустую базу
	httpGet = func(_ string) (*http.Response, error) { return &validResp, nil }
	CovidTask(db)

	var stats []models.CovidStat
	db.Find(&stats)

	assert.Len(t, stats, 2)
	for _, s := range stats {
		assert.Equal(t, validStat.Confirmed, s.Confirmed)
		assert.Equal(t, validStat.Deaths, s.Deaths)
		assert.Equal(t, validStat.Recovered, s.Recovered)
		assert.Equal(t, validStat.ConfirmedDiff, s.ConfirmedDiff)
		assert.Equal(t, validStat.DeathsDiff, s.DeathsDiff)
		assert.Equal(t, validStat.RecoveredDiff, s.RecoveredDiff)
		assert.Equal(t, validStat.LastUpdate, s.LastUpdate)
		assert.Equal(t, validStat.Active, s.Active)
		assert.Equal(t, validStat.ActiveDiff, s.ActiveDiff)
		assert.Equal(t, validStat.FatalityRate, s.FatalityRate)
	}

	resp := http.Response{Body: fakeBody{
		content: `{"data":[{"date":"2021-03-27","confirmed":0,"deaths":0,"recovered":0,` +
			`"confirmed_diff":0,"deaths_diff":0,"recovered_diff":0,"last_update":"2021-03-28 05:22:20",` +
			`"active":0,"active_diff":0,"fatality_rate":0,"region":{"iso":"RUS","name":"Russia",` +
			`"province":"Moscow","lat":"0","long":"0","cities":[]}}]}`,
	}}

	// проверяем обновление данных
	httpGet = func(_ string) (*http.Response, error) { return &resp, nil }
	CovidTask(db)

	db.Find(&stats)

	assert.Len(t, stats, 2)
	for _, s := range stats {
		assert.Equal(t, int64(0), s.Confirmed)
		assert.Equal(t, int64(0), s.Deaths)
		assert.Equal(t, int64(0), s.Recovered)
		assert.Equal(t, int64(0), s.ConfirmedDiff)
		assert.Equal(t, int64(0), s.DeathsDiff)
		assert.Equal(t, int64(0), s.RecoveredDiff)
		assert.Equal(t, "2021-03-28 05:22:20", s.LastUpdate)
		assert.Equal(t, int64(0), s.Active)
		assert.Equal(t, int64(0), s.ActiveDiff)
		assert.Equal(t, float64(0), s.FatalityRate)
	}

	// ломаем базу
	db.Exec("drop table covid_stats")
	assert.NotPanics(t, func() { CovidTask(db) })
}
