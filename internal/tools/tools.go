package tools

import (
	"fmt"
	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/scheduler/subscriptions"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
	if err := subscriptions.InitSubscriptions(db); err != nil {
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
