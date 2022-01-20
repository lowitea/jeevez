package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

// TestInitConfig тестирует инициализацию конфига
func TestInitConfig(t *testing.T) {
	// ансетим обязательную переменную
	_ = os.Unsetenv("JEEVEZ_DB_USER")
	_, err := InitConfig()
	assert.Error(t, err)

	// устанавливаем необходимые переменные в окружение
	t.Setenv("JEEVEZ_APP_VERSION", "1.2.3")
	t.Setenv("JEEVEZ_TELEGRAM_TOKEN", "telegram_token")
	t.Setenv("JEEVEZ_TELEGRAM_ADMIN", "654865")
	t.Setenv("JEEVEZ_DB_USER", "test_user")
	t.Setenv("JEEVEZ_DB_PASSWORD", "db_password")
	t.Setenv("JEEVEZ_CURRENCYAPI_TOKEN", "currency_token")
	t.Setenv("JEEVEZ_TELEGRAM_BOTNAME", "test_bot")

	expCfg := Config{
		App: struct {
			Version string `default:"x.x.x (dev)"`
		}{"1.2.3"},
		Telegram: struct {
			BotName string `required:"true"`
			Token   string `required:"true"`
			Admin   int64  `required:"true"`
		}{"test_bot", "telegram_token", 654865},
		DB: struct {
			Host     string `default:"jeevez-database"`
			Port     int    `default:"5432"`
			User     string `required:"true"`
			Password string `required:"true"`
			Name     string `default:"jeevez"`
		}{"jeevez-database", 5432, "test_user", "db_password", "jeevez"},
		CurrencyAPI: struct {
			Token string `required:"true"`
		}{"currency_token"},
	}

	cfg, err := InitConfig()
	require.NoError(t, err)

	// проверяем вернувшийся конфиг
	assert.Equal(t, &expCfg, cfg)

	// проверяем глобальный конфиг
	assert.Equal(t, expCfg, Cfg)
}
