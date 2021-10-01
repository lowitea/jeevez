package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"testing"
)

func TestCovid(t *testing.T) {
	assert.NotPanics(t, func() { _ = covid(&cli.Context{}) })
}

func TestCurrency(t *testing.T) {
	assert.NotPanics(t, func() { _ = currency(&cli.Context{}) })
}
