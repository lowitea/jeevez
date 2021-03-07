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
	ChatID         int64     `gorm:"primaryKey"`
	SubscriptionID int64     `gorm:"primaryKey"`
	CreatedAt      time.Time `gorm:"autoUpdateTime"`
	Time           int64     `gorm:"index"`
	HumanTime      string
}

func (ChatSubscription) BeforeCreate(db *gorm.DB) error {
	return db.SetupJoinTable(&Chat{}, "Subscriptions", &ChatSubscription{})
}
