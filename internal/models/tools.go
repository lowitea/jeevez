package models

import (
	"gorm.io/gorm"
	"log"
)

// MigrateAll выполняет автомиграцию базы данных
func MigrateAll(db *gorm.DB) error {
	log.Print("AutoMigrating has started")
	if err := db.AutoMigrate(&CurrencyRate{}); err != nil {
		log.Printf("migrating CurrencyRate error: %s", err)
		return err
	}
	if err := db.AutoMigrate(&Chat{}); err != nil {
		log.Printf("migrating Chat error: %s", err)
		return err
	}
	if err := db.AutoMigrate(&Subscription{}); err != nil {
		log.Printf("migrating Subscription error: %s", err)
		return err
	}
	if err := db.AutoMigrate(&ChatSubscription{}); err != nil {
		log.Printf("migrating ChatSubscription error: %s", err)
		return err
	}
	if err := db.AutoMigrate(&CovidStat{}); err != nil {
		log.Printf("migrating CovidStat error: %s", err)
		return err
	}
	return nil
}
