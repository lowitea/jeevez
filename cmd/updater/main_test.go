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
	if val, ok := os.LookupEnv("JEEVEZ_TEST_DB_HOST"); ok {
		testCfg.DB.Host = val
	}
	if val, ok := os.LookupEnv("JEEVEZ_TEST_DB_USER"); ok {
		testCfg.DB.User = val
	}
	if val, ok := os.LookupEnv("JEEVEZ_TEST_DB_PASSWORD"); ok {
		testCfg.DB.Password = val
	}
	if val, ok := os.LookupEnv("JEEVEZ_TEST_DB_NAME"); ok {
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

// TestMainFunc смок тест на ошибку в основной функции запуска cli
func TestMainFunc(t *testing.T) {
	for _, envName := range [...]string{
		"JEEVEZ_TELEGRAM_TOKEN",
		"JEEVEZ_TELEGRAM_ADMIN",
		"JEEVEZ_CURRENCYAPI_TOKEN",
	} {
		if _, ok := os.LookupEnv(envName); !ok {
			_ = os.Setenv(envName, "1")
			defer func() { _ = os.Unsetenv(envName) }()
		}
	}

	if val, ok := os.LookupEnv("JEEVEZ_DB_USER"); ok {
		defer func(val string) { _ = os.Setenv("JEEVEZ_DB_USER", val) }(val)
	}

	if val, ok := os.LookupEnv("JEEVEZ_DB_PASSWORD"); ok {
		defer func(val string) { _ = os.Setenv("JEEVEZ_DB_PASSWORD", val) }(val)
	}

	if val, ok := os.LookupEnv("JEEVEZ_DB_HOST"); ok {
		defer func(val string) { _ = os.Setenv("JEEVEZ_DB_HOST", val) }(val)
	}

	if val, ok := os.LookupEnv("JEEVEZ_TEST_DB_USER"); ok {
		_ = os.Setenv("JEEVEZ_DB_USER", val)
	} else {
		_ = os.Setenv("JEEVEZ_DB_USER", "test")
	}

	if val, ok := os.LookupEnv("JEEVEZ_TEST_DB_PASSWORD"); ok {
		_ = os.Setenv("JEEVEZ_DB_PASSWORD", val)
	} else {
		_ = os.Setenv("JEEVEZ_DB_PASSWORD", "test")
	}

	if val, ok := os.LookupEnv("JEEVEZ_TEST_DB_HOST"); ok {
		_ = os.Setenv("JEEVEZ_DB_HOST", val)
	} else {
		_ = os.Setenv("JEEVEZ_DB_HOST", "localhost")
	}

	assert.Panics(t, main)
}
