package config

type Config struct {
	App struct {
		Version string
	}
	Telegram struct {
		Token string
		Admin int64
	}
}
