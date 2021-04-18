package testTools

import (
	"github.com/lowitea/jeevez/internal/tools"
	"gorm.io/gorm"
	"log"
	"os"
)

const testDBName = "jeevez_test"

// InitTestDB инициализируем тестовую базу данных
func InitTestDB() *gorm.DB {
	host := os.Getenv("JEEVEZ_TEST_DB_HOST")
	if host == "" {
		host = "localhost"
	}

	// подключаемся к public, для подготовки тестовой схемы
	db, err := tools.ConnectDB(host, 5432, "test", "test", testDBName)
	if err != nil {
		log.Fatalf("error init test db %s", err)
	}

	db.Exec("DROP SCHEMA public CASCADE")
	db.Exec("CREATE SCHEMA public")

	// настраиваем базу
	if err := tools.SetupDB(db); err != nil {
		log.Fatalf("error init test db %s", err)
	}

	return db
}
