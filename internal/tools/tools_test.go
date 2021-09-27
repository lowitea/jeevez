package tools

import (
	"errors"
	"github.com/lowitea/jeevez/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

// TestCheck проверяет функцию-декоратор для обработки ошибки
func TestCheck(t *testing.T) {
	assert.PanicsWithError(
		t,
		"test error",
		func() { Check(errors.New("test error")) },
	)
}

// TestInitSubscriptions проверяет функцию сохранения в базу списка подписок
func TestInitSubscriptions(t *testing.T) {
	host := os.Getenv("JEEVEZ_TEST_DB_HOST")
	if host == "" {
		host = "localhost"
	}

	db, err := ConnectDB(host, 5432, "test", "test", "jeevez_test")
	assert.NoError(t, err)

	defer func(m map[string]models.Subscription) {
		models.SubscrNameSubscrMap = m
		models.MigrateAll(db)
		InitSubscriptions(db)
	}(models.SubscrNameSubscrMap)

	err = db.Delete(&models.ChatSubscription{}, "true").Error
	assert.NoError(t, err)

	err = db.Delete(&models.Subscription{}, "true").Error
	assert.NoError(t, err)

	// проверяем что записи были созданы
	InitSubscriptions(db)
	for _, subscr := range models.SubscrNameSubscrMap {
		subscrDB := models.Subscription{}
		err := db.First(&subscrDB, "id = ? AND name = ?", subscr.ID, subscr.Name).Error
		assert.NoError(t, err)
	}

	// проверяем ошибку создания записи
	err = db.Exec("drop table subscriptions cascade").Error
	require.NoError(t, err)
	assert.Panics(t, func() { InitSubscriptions(db) })
	models.MigrateAll(db)

	// проверяем сохранение при изменённом описании
	assert.NotPanics(t, func() { InitSubscriptions(db) })
	s := models.SubscrNameSubscrMap["covid19-russia"]
	s.Description = ""
	models.SubscrNameSubscrMap["covid19-russia"] = s
	assert.NotPanics(t, func() { InitSubscriptions(db) })
}
