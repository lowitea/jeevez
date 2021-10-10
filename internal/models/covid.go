package models

import "time"

// covidStat структура с ответом от апи статистки по ковиду
type CovidStat struct {
	ID               int64     `gorm:"primaryKey"`
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`
	HumanName        string
	SubscriptionName string `gorm:"uniqueIndex"`
	Confirmed        int64
	Deaths           int64
	ConfirmedDiff    int64
	DeathsDiff       int64
	LastUpdate       string
	Active           int64
	ActiveDiff       int64
	FatalityRate     float64
}
