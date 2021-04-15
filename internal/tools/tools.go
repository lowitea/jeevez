package tools

import (
	"errors"
	"fmt"
	"github.com/lowitea/jeevez/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// InitSubscriptions создаёт в базе недостающие подписки
func InitSubscriptions(db *gorm.DB) error {
	log.Print("InitSubscriptions has started")
	for _, subscr := range models.SubscrNameSubscrMap {
		// пытаемся получить подписку из базы по id и name
		subscrDB := models.Subscription{}
		result := db.First(&subscrDB, "id = ? AND name = ?", subscr.ID, subscr.Name)

		// если такого не нашлось, создаём
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// так как у нас захардкожены id в коде, нужно попробовать удалить из базы запись с таким id
			_ = db.Delete(&models.ChatSubscription{}, "subscription_id = ?", subscr.ID)
			_ = db.Delete(&models.Subscription{}, subscr.ID)

			// создаём новую запись
			if result = db.Create(&subscr); result.Error != nil {
				log.Printf("create Subscription error: %s", result.Error)
				return result.Error
			}
			continue
		} else if result.Error != nil {
			log.Printf("update Subscription error: %s", result.Error)
			return result.Error
		}

		// обновляем запись если отличаются другие поля
		if subscr.Description != subscrDB.Description {
			db.Save(&subscr)
		}
	}
	return nil
}

// ConnectDB открываем коннект к базе данных
func ConnectDB(host string, port int, user string, pwd string, dbName string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", host, user, pwd, dbName, port)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// SetupDB подготавливаем данные в базе
func SetupDB(db *gorm.DB) error {
	// миграция моделей
	if err := models.MigrateAll(db); err != nil {
		return err
	}

	// инициализация вариантов подписок
	if err := InitSubscriptions(db); err != nil {
		return err
	}
	return nil
}

// InitDB инициализируем продовую базу данных
func InitDB(host string, port int, user string, pwd string, dbName string) (*gorm.DB, error) {
	db, err := ConnectDB(host, port, user, pwd, dbName)
	if err != nil {
		return nil, err
	}

	// настраиваем базу
	if err := SetupDB(db); err != nil {
		return nil, err
	}

	return db, err
}
