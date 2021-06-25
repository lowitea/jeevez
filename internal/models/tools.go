package models

import (
	"gorm.io/gorm"
	"log"
)

// MigrateAll выполняет автомиграцию базы данных
func MigrateAll(db *gorm.DB) error {
	log.Print("AutoMigrating has started")

	err := db.AutoMigrate(
		&CurrencyRate{},
		&Chat{},
		&Subscription{},
		&ChatSubscription{},
		&CovidStat{},
	)

	if err == nil {
		return nil
	}

	log.Printf("AutoMigrating error: %s", err)
	return err
}
