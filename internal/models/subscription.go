package models

import (
	"gorm.io/gorm"
	"time"
)

// Subscription модель доступной подписки
type Subscription struct {
	ID          int64 `gorm:"primaryKey"`
	Name        string
	Description string
}

// ChatSubscription модель связи чата с подписками
type ChatSubscription struct {
	Chat           Chat
	Subscription   Subscription
	ChatID         int64     `gorm:"primaryKey"`
	SubscriptionID int64     `gorm:"primaryKey"`
	CreatedAt      time.Time `gorm:"autoUpdateTime"`
	Time           int64     `gorm:"index"`
	HumanTime      string
}

func (ChatSubscription) BeforeCreate(db *gorm.DB) error {
	return db.SetupJoinTable(&Chat{}, "Subscriptions", &ChatSubscription{})
}

var SubscrNameSubscrMap = map[string]Subscription{
	"covid19-russia": {
		ID:          1,
		Name:        "covid19-russia",
		Description: "Дневная статистика по COViD-19 по России",
	},
	"covid19-moscow": {
		ID:          2,
		Name:        "covid19-moscow",
		Description: "Дневная статистика по COViD-19 по Москве",
	},
	"yoga-test": {
		ID:          3,
		Name:        "yoga-test",
		Description: "Ежедневный тест с позами йоги",
	},
	"weather-moscow": {
		ID:          4,
		Name:        "weather-moscow",
		Description: "Текущая погода в Москве",
	},
}
