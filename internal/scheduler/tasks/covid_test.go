package tasks

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestCovidStat_Update проверяет функцию обновления объекта статистики
func TestCovidStat_Update(t *testing.T) {
	stat := covidStat{
		Confirmed: 1, Deaths: 2, Recovered: 3, ConfirmedDiff: 4, DeathsDiff: 5,
		RecoveredDiff: 6, LastUpdate: "9.01.20", Active: 7, ActiveDiff: 8, FatalityRate: 0.5,
	}
	statNew := covidStat{
		Confirmed: 11, Deaths: 12, Recovered: 13, ConfirmedDiff: 14, DeathsDiff: 15,
		RecoveredDiff: 16, LastUpdate: "10.01.20", Active: 17, ActiveDiff: 18, FatalityRate: 0.8,
	}
	stat.Update(statNew)
	statExp := covidStat{
		Confirmed: 12, Deaths: 14, Recovered: 16, ConfirmedDiff: 18, DeathsDiff: 20,
		RecoveredDiff: 22, LastUpdate: "10.01.20", Active: 24, ActiveDiff: 26, FatalityRate: 0.65,
	}
	assert.Equal(t, statExp, stat)

	stat.FatalityRate = 0
	stat.Update(statNew)
	assert.Equal(t, 0.8, stat.FatalityRate)
}
