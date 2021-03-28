package tools

import (
	"fmt"
	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/scheduler/subscriptions"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(host string, port int, user string, pwd string, dbName string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", host, user, pwd, dbName, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// миграция моделей
	if err := models.MigrateAll(db); err != nil {
		return nil, err
	}

	// инициализация вариантов подписок
	if err := subscriptions.InitSubscriptions(db); err != nil {
		return nil, err
	}
	return db, err
}
