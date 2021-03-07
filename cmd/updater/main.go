package main

import (
	"fmt"
	"github.com/lowitea/jeevez/internal/config"
	"github.com/lowitea/jeevez/internal/scheduler/tasks"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	// инициализируем конфиг
	cfg, err := config.InitConfig()
	if err != nil {
		log.Printf("env parse error %s", err)
		os.Exit(1)
	}

	// инициализация базы
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d",
		cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName, cfg.DB.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("failed to connect database: %s", err)
		os.Exit(1)
	}

	app := &cli.App{}

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

	if err = app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
