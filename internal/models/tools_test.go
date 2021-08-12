package models

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

// TestMigrateAll смоук тест накатки миграций
func TestMigrateAll(t *testing.T) {
	autoMigrate = func(g *gorm.DB, i ...interface{}) error { return nil }
	assert.NotPanics(t, func() { MigrateAll(&gorm.DB{}) })

	autoMigrate = func(g *gorm.DB, i ...interface{}) error { return errors.New("test") }
	assert.Panicsf(t, func() { MigrateAll(&gorm.DB{}) }, "")
}
