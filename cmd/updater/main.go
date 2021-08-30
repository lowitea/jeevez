package main

import (
	"github.com/lowitea/jeevez/internal/config"
	"github.com/lowitea/jeevez/internal/scheduler/tasks"
	"github.com/lowitea/jeevez/internal/tools"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func initApp(initCfgFunc func() (*config.Config, error)) *cli.App {
	// инициализируем конфиг
	cfg, err := initCfgFunc()
	if err != nil {
		log.Fatalf("env parse error %s", err)
	}

	db, err := tools.ConnectDB(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)
	if err != nil {
		log.Fatalf("db connect error %s", err)
	}

	app := &cli.App{Usage: "A cli app for update date in Jeevez"}

	app.Commands = []*cli.Command{
		{
			Name:  "covid",
			Usage: "update covid19 stat",
			Action: func(c *cli.Context) error {
				tasks.CovidTask(db)
				return nil
			},
		},
		{
			Name:  "currency",
			Usage: "update currency stat",
			Action: func(c *cli.Context) error {
				tasks.CurrencyTask(db)
				return nil
			},
		},
	}

	return app
}

func main() {
	app := initApp(config.InitConfig)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
