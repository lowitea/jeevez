package models

import (
	"time"
)

// Chat модель зарегистрированного чата
type Chat struct {
	ID            int64 `gorm:"primaryKey"`
	TgID          int64 `gorm:"uniqueIndex"`
	TgName        string
	TgTitle       string
	TgType        string
	RegisteredAt  time.Time      `gorm:"autoCreateTime"`
	Subscriptions []Subscription `gorm:"many2many:chat_subscriptions;"`
}
