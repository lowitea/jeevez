package testTools

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// TestInitTestDB проверяет работу функции подключения к тестовой базе
func TestInitTestDB(t *testing.T) {
	db := InitTestDB()
	assert.NotNil(t, db)

	if dbHost, ok := os.LookupEnv("JEEVEZ_TEST_DB_HOST"); ok {
		defer func(host string) { _ = os.Setenv("JEEVEZ_TEST_DB_HOST", host) }(dbHost)
	} else {
		defer func() { _ = os.Unsetenv("JEEVEZ_TEST_DB_HOST") }()
	}
	err := os.Setenv("JEEVEZ_TEST_DB_HOST", "not_exist_host")
	assert.NoError(t, err)
	assert.Panics(t, func() { InitTestDB() })

	err = os.Unsetenv("JEEVEZ_TEST_DB_HOST")
	assert.NoError(t, err)

	defer func(lh string) { localhost = lh }(localhost)
	localhost = "not_exist_host"
	assert.Panics(t, func() { InitTestDB() })
}
