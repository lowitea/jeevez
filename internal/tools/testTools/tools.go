package testTools

import (
	"fmt"
	"github.com/lowitea/jeevez/internal/tools"
	"gorm.io/gorm"
	"log"
	"os"
)

const testDBName = "jeevez_test"

var localhost = "localhost"

// InitTestDB инициализируем тестовую базу данных
func InitTestDB() *gorm.DB {
	host := os.Getenv("JEEVEZ_TEST_DB_HOST")
	if host == "" {
		host = localhost
	}

	// подключаемся к public, для подготовки тестовой схемы
	db, err := tools.ConnectDB(host, 5432, "test", "test", testDBName)
	if err != nil {
		log.Panicf("error init test db %s\n", err)
	}

	// очищаем существующую тестовую базу
	clearQuery := fmt.Sprintf(""+
		"select 'drop table if exists \"' || tablename || '\" cascade;' from pg_tables "+
		"where schemaname = '%s'", testDBName)
	db.Exec(clearQuery)

	// настраиваем базу
	tools.SetupDB(db)

	return db
}
