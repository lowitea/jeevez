package testTools

import (
	"fmt"
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

	// очищаем существующую тестовую базу
	clearQuery := fmt.Sprintf(""+
		"select 'drop table if exists \"' || tablename || '\" cascade;' from pg_tables "+
		"where schemaname = '%s'", testDBName)
	if result := db.Exec(clearQuery); result.Error != nil {
		log.Fatalf("error clear test db %s", result.Error)
	}

	// настраиваем базу
	tools.SetupDB(db)

	return db
}
