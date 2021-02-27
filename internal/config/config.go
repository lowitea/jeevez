package config

type Config struct {
	App struct {
		Version string `default:"x.x.x (dev)"`
	}
	Telegram struct {
		Token string `required:"true"`
		Admin int64  `required:"true"`
	}
	DB struct {
		Host     string `default:"jeevez-database"`
		Port     int    `default:"5432"`
		User     string `required:"true"`
		Password string `required:"true"`
		DBName   string `default:"jeevez"`
	}
	CurrencyAPI struct {
		Token string `required:"true"`
	}
}
