package tasks

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lowitea/jeevez/internal/models"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var httpGet = http.Get

// covidStat структура с ответом от апи статистки по ковиду
type covidStat struct {
	Confirmed     int64   `json:"confirmed"`
	Deaths        int64   `json:"deaths"`
	Recovered     int64   `json:"recovered"`
	ConfirmedDiff int64   `json:"confirmed_diff"`
	DeathsDiff    int64   `json:"deaths_diff"`
	RecoveredDiff int64   `json:"recovered_diff"`
	LastUpdate    string  `json:"last_update"`
	Active        int64   `json:"active"`
	ActiveDiff    int64   `json:"active_diff"`
	FatalityRate  float64 `json:"fatality_rate"`
}

// Update метод для сложения структур covidStat
func (stat *covidStat) Update(updStat covidStat) *covidStat {
	stat.Confirmed += updStat.Confirmed
	stat.ConfirmedDiff += updStat.ConfirmedDiff
	stat.Deaths += updStat.Deaths
	stat.DeathsDiff += updStat.DeathsDiff
	stat.Recovered += updStat.Recovered
	stat.RecoveredDiff += updStat.RecoveredDiff
	stat.LastUpdate = updStat.LastUpdate
	stat.Active += updStat.Active
	stat.ActiveDiff += updStat.ActiveDiff
	if stat.FatalityRate != 0 && updStat.FatalityRate != 0 {
		stat.FatalityRate = (stat.FatalityRate + updStat.FatalityRate) / 2
	} else if stat.FatalityRate == 0 {
		stat.FatalityRate = updStat.FatalityRate
	}
	return stat
}

// getData получить список данных по ковид из апи
func getData(url string) ([]covidStat, error) {
	// apiResp ответ от апи статистики
	type apiResp struct {
		Data []covidStat
	}

	resp, err := httpGet(url)
	if err != nil {
		log.Printf("Error get data: %s", err)
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error get data: %s", err)
		return nil, err
	}

	var data apiResp
	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("Error get data: %s", err)
		return nil, err
	}
	return data.Data, nil
}

// getStat получения статистики по ссылке
func getStat(url string) (*covidStat, error) {
	// запрашиваем статистику по ковиду, сначала за вчера, но если не получилось, то за позавчера
	var stats []covidStat
	loc, _ := time.LoadLocation("Europe/Moscow")
	t := time.Now()

	for _, day := range [...]int{-1, -2} {
		dt := t.AddDate(0, 0, day).In(loc).Format("2006-01-02")
		var err error
		stats, err = getData(fmt.Sprintf(url, dt))
		if err != nil {
			return nil, fmt.Errorf("getting covid stats error: %s", err)
		}
		if len(stats) != 0 {
			break
		}
	}

	if len(stats) == 0 {
		return nil, errors.New("getting covid stats error")
	}

	var result covidStat
	for _, stat := range stats {
		result.Update(stat)
	}

	return &result, nil
}

// CovidTask таска рассылающая статистику по ковиду
func CovidTask(db *gorm.DB) {
	log.Printf("CovidTask has started")
	if db == nil {
		log.Printf("db is nil")
		return
	}

	for statName, statConf := range subscrUrlMap {
		respStat, err := getStat(statConf.UrlTpl)
		if err != nil {
			log.Printf("Error get data: %s", err)
			continue
		}

		stat := models.CovidStat{
			SubscriptionName: statName,
			HumanName:        statConf.HumanName,
			Confirmed:        respStat.Confirmed,
			Deaths:           respStat.Deaths,
			Recovered:        respStat.Recovered,
			ConfirmedDiff:    respStat.ConfirmedDiff,
			DeathsDiff:       respStat.DeathsDiff,
			RecoveredDiff:    respStat.RecoveredDiff,
			LastUpdate:       respStat.LastUpdate,
			Active:           respStat.Active,
			ActiveDiff:       respStat.ActiveDiff,
			FatalityRate:     respStat.FatalityRate,
		}

		// создаём или обновляем статистику по ковиду
		var statDB models.CovidStat
		result := db.First(&statDB, "subscription_name = ?", statName)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			db.Create(&stat)
		} else if result.Error != nil {
			log.Printf("getting CovidStat from db error: %s", result.Error)
			return
		} else {
			stat.ID = statDB.ID
			db.Save(&stat)
		}
	}
}
