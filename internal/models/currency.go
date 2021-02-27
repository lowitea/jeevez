package models

import "time"

type CurrencyRate struct {
	ID        uint      `gorm:"primaryKey"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Value     float64
	Name      string
}
