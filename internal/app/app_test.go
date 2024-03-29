package app

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lowitea/jeevez/internal/tools/testtools"
)

// // TestInitApp смоук тест инициализации приложения
// func TestInitApp(t *testing.T) {
// 	testCfg := config.Config{
// 		Telegram: struct {
// 			BotName string `required:"true"`
// 			Token   string `required:"true"`
// 			Admin   int64  `required:"true"`
// 		}{"", "1", 1},
// 		DB: struct {
// 			Host     string `default:"jeevez-database"`
// 			Port     int    `default:"5432"`
// 			User     string `required:"true"`
// 			Password string `required:"true"`
// 			Name     string `default:"jeevez"`
// 		}{"", 5432, "test", "test", "jeevez_test"},
// 		CurrencyAPI: struct {
// 			Token string `required:"true"`
// 		}{"1"},
// 	}

// 	cases := [...]struct {
// 		Name        string
// 		DBHost      string
// 		InitCfgFunc func() (*config.Config, error)
// 		InitBotFunc func(_ string) (*tgbotapi.BotAPI, error)
// 		ErrMsg      string
// 	}{
// 		{
// 			"cfg_err",
// 			"",
// 			func() (*config.Config, error) { return nil, errors.New("test") },
// 			func(_ string) (*tgbotapi.BotAPI, error) { return nil, nil },
// 			"env parse error: test",
// 		},
// 		{
// 			"bot_err",
// 			"",
// 			func() (*config.Config, error) { return &testCfg, nil },
// 			func(_ string) (*tgbotapi.BotAPI, error) { return nil, errors.New("test") },
// 			"error connect to telegram: test",
// 		},
// 		{
// 			"db_err",
// 			"bad_host",
// 			func() (*config.Config, error) { return &testCfg, nil },
// 			func(_ string) (*tgbotapi.BotAPI, error) { return &tgbotapi.BotAPI{}, nil },
// 			"error init database: test",
// 		},
// 	}

// 	for _, c := range cases {
// 		t.Run(fmt.Sprintf("cmd=%s", c.Name), func(t *testing.T) {
// 			testCfg.DB.Host = c.DBHost
// 			_, err := initApp(c.InitBotFunc, c.InitCfgFunc)
// 			assert.Errorf(t, err, c.ErrMsg)
// 		})
// 	}

// 	testCfg.DB.Host = os.Getenv("JEEVEZ_TEST_DB_HOST")
// 	initCfgFunc := func() (*config.Config, error) { return &testCfg, nil }
// 	initBotFunc := func(_ string) (*tgbotapi.BotAPI, error) { return &tgbotapi.BotAPI{}, nil }
// 	_, err := initApp(initBotFunc, initCfgFunc)
// 	assert.NoError(t, err)
// }

// TestReleaseNotify проверяет функцию отправки сообщения о релизе
func TestReleaseNotify(t *testing.T) {
	var adminID int64 = 666
	expMsg := tgbotapi.NewMessage(adminID, "🤵🏻 Я обновился! :)\nМоя новая версия: 6.6.6")
	botAPIMock := testtools.NewBotAPIMock(expMsg)
	releaseNotify(botAPIMock, adminID, "6.6.6")
	botAPIMock.AssertExpectations(t)
}

// // TestProcessUpdate смоук тест общего запуска хендлеров
// func TestProcessUpdate(t *testing.T) {
// 	db := testtools.InitTestDB()
// 	update := testtools.NewUpdate("no_command")
// 	botAPIMock := testtools.NewBotAPIMock(tgbotapi.MessageConfig{})
// 	cfg := &config.Config{Telegram: struct {
// 		BotName string `required:"true"`
// 		Token   string `required:"true"`
// 		Admin   int64  `required:"true"`
// 	}{}}

// 	dep := &appDepContainer{
// 		db:  db,
// 		bot: botAPIMock,
// 		cfg: cfg,
// 	}

// 	assert.NotPanics(t, func() { processUpdate(update, dep) })
// 	botAPIMock.AssertNotCalled(t, "Send")

// 	assert.NotPanics(t, func() { processUpdate(tgbotapi.Update{}, dep) })
// 	botAPIMock.AssertNotCalled(t, "Send")
// }
