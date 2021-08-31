package main

import (
	"errors"
	"github.com/lowitea/jeevez/internal/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// TestInitApp смоук тест для функции инициализации cli приложения
func TestRunner(t *testing.T) {
	testCfg := config.Config{
		Telegram: struct {
			Token string `required:"true"`
			Admin int64  `required:"true"`
		}{"1", 1},
		DB: struct {
			Host     string `default:"jeevez-database"`
			Port     int    `default:"5432"`
			User     string `required:"true"`
			Password string `required:"true"`
			Name     string `default:"jeevez"`
		}{"", 5432, "test", "test", "jeevez_test"},
		CurrencyAPI: struct {
			Token string `required:"true"`
		}{"1"},
	}
	if val, ok := os.LookupEnv("JEEVEZ_TEST_DB_HOST"); ok == true {
		testCfg.DB.Host = val
	}
	if val, ok := os.LookupEnv("JEEVEZ_TEST_DB_USER"); ok == true {
		testCfg.DB.User = val
	}
	if val, ok := os.LookupEnv("JEEVEZ_TEST_DB_PASSWORD"); ok == true {
		testCfg.DB.Password = val
	}
	if val, ok := os.LookupEnv("JEEVEZ_TEST_DB_NAME"); ok == true {
		testCfg.DB.Name = val
	}

	assert.NotPanics(t, func() { initApp(func() (*config.Config, error) { return &testCfg, nil }) })

	// ошибка инициализации конфига
	assert.PanicsWithValue(
		t,
		"env parse error test",
		func() { initApp(func() (*config.Config, error) { return nil, errors.New("test") }) },
	)

	// ошибка инициализации базы
	testCfg.DB.Host = `not_host`
	assert.Panics(t, func() { initApp(func() (*config.Config, error) { return &testCfg, nil }) })
}
