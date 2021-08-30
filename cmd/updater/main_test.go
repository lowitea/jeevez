package main

import (
	"github.com/lowitea/jeevez/internal/config"
	"github.com/stretchr/testify/assert"
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

	assert.NotPanics(t, func() { initApp(func() (*config.Config, error) { return &testCfg, nil }) })
}
