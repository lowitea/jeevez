package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

// TestInitConfig тестирует инициализацию конфига
func TestInitConfig(t *testing.T) {
	// устанавливаем необходимые переменные в окружение
	_ = os.Setenv("JEEVEZ_APP_VERSION", "1.2.3")
	_ = os.Setenv("JEEVEZ_TELEGRAM_TOKEN", "telegram_token")
	_ = os.Setenv("JEEVEZ_TELEGRAM_ADMIN", "654865")
	_ = os.Setenv("JEEVEZ_DB_USER", "test_user")
	_ = os.Setenv("JEEVEZ_DB_PASSWORD", "db_password")
	_ = os.Setenv("JEEVEZ_CURRENCYAPI_TOKEN", "currency_token")

	expCfg := Config{
		App: struct {
			Version string `default:"x.x.x (dev)"`
		}{"1.2.3"},
		Telegram: struct {
			Token string `required:"true"`
			Admin int64  `required:"true"`
		}{"telegram_token", 654865},
		DB: struct {
			Host     string `default:"jeevez-database"`
			Port     int    `default:"5432"`
			User     string `required:"true"`
			Password string `required:"true"`
			DBName   string `default:"jeevez"`
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
