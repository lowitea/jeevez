package main

import (
	"github.com/lowitea/jeevez/internal/config"
	"github.com/lowitea/jeevez/internal/tools"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
	"log"
	"os"
)

var db *gorm.DB

func initApp(initCfgFunc func() (*config.Config, error)) *cli.App {
	// инициализируем конфиг
	cfg, err := initCfgFunc()
	if err != nil {
		log.Panicf("env parse error %s", err)
	}

	db, err = tools.ConnectDB(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)
	if err != nil {
		log.Panicf("db connect error %s", err)
	}

	app := &cli.App{Usage: "A cli app for update date in Jeevez"}

	app.Commands = []*cli.Command{
		{
			Name:   "covid",
			Usage:  "update covid19 stat",
			Action: covid,
		},
		{
			Name:   "currency",
			Usage:  "update currency stat",
			Action: currency,
		},
	}

	return app
}

func main() {
	if err := initApp(config.InitConfig).Run(os.Args); err != nil {
		log.Panicln(err)
	}
}
