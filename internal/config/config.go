package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App struct {
		Version string `default:"x.x.x (dev)"`
	}
	Telegram struct {
		BotName string `required:"true"`
		Token   string `required:"true"`
	}
	DB struct {
		Host     string `default:"jeevez-database"`
		Port     int    `default:"5432"`
		User     string `required:"true"`
		Password string `required:"true"`
		Name     string `default:"jeevez"`
	}
	CurrencyAPI struct {
		Token string `required:"true"`
	}
	WeatherAPI struct {
		Token string `required:"true"`
	}
	Admin struct {
		TelegramID int64  `required:"true"`
		Email      string `required:"true"`
	}
	Mail struct {
		Host           string `required:"true"`
		Login          string `required:"true"`
		Password       string `required:"true"`
		PrimaryDomain  string `required:"true"`
		TempMailDomain string `required:"true"`
	}
}

var Cfg Config

func InitConfig() (*Config, error) {
	if err := envconfig.Process("jeevez", &Cfg); err != nil {
		return nil, err
	}
	return &Cfg, nil
}
