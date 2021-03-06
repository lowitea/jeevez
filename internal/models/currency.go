package models

import "time"

// CurrencyRate модель валютной пары
type CurrencyRate struct {
	ID        int64     `gorm:"primaryKey"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Value     float64
	Name      string
}
