package models

import "time"

// CurrencyRate модель валютной пары
type CurrencyRate struct {
	ID        uint      `gorm:"primaryKey"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Value     float64
	Name      string
}
