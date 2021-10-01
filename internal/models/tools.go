package models

import (
	"fmt"
	"gorm.io/gorm"
	"log"
)

var autoMigrate = (*gorm.DB).AutoMigrate

// MigrateAll выполняет автомиграцию базы данных
func MigrateAll(db *gorm.DB) {
	log.Print("AutoMigrating has started")

	err := autoMigrate(
		db,
		&CurrencyRate{},
		&Chat{},
		&Subscription{},
		&ChatSubscription{},
		&CovidStat{},
	)
	if err != nil {
		panic(fmt.Sprintf("AutoMigrating failed: %s", err))
	}
}
