package models

import (
	"gorm.io/gorm"
	"log"
)

// MigrateAll выполняет автомиграцию базы данных
func MigrateAll(db *gorm.DB) error {
	log.Print("AutoMigrating has started")

	return db.AutoMigrate(
		&CurrencyRate{},
		&Chat{},
		&Subscription{},
		&ChatSubscription{},
		&CovidStat{},
	)

}
