package config

type Config struct {
	App struct {
		Version string `default:"x.x.x (dev)"`
	}
	Telegram struct {
		Token string `required:"true"`
		Admin int64  `required:"true"`
	}
}
