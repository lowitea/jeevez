package main

import (
	"github.com/lowitea/jeevez/internal/scheduler/tasks"
	"github.com/urfave/cli/v2"
)

func covid(_ *cli.Context) error {
	tasks.CovidTask(db)
	return nil
}

func currency(_ *cli.Context) error {
	tasks.CurrencyTask(db)
	return nil
}
