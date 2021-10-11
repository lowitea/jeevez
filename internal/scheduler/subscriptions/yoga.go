package subscriptions

import (
	"embed"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/structs"
	"gorm.io/gorm"
	"math/rand"
	"path"
	"time"
)

// максимальная длина названия позы - 26 символов
//goland:noinspection SpellCheckingInspection
var yogaPoses = map[string]string{
	"Джатхара паривартанасана":   "1.jpg",
	"Падангуштхасана":            "2.jpg",
	"Парипурна навасана":         "3.jpg",
	"Упавиштха конасана":         "4.jpg",
	"Прасарита падоттанасана":    "5.jpg",
	"Випарита карани":            "6.jpg",
	"Пашчимоттанасана":           "7.jpg",
	"Адхо мукха шванасана":       "8.jpg",
	"Бхуджангасана":              "9.jpg",
	"Бадха конасана":             "10.jpg",
	"Ардха чандрасана":           "11.jpg",
	"Вирабхадрасана 1":           "12.jpg",
	"Вирабхадрасана 2":           "13.jpg",
	"Вирабхадрасана 3":           "14.jpg",
	"Врикшасана":                 "15.jpg",
	"Маричиасана":                "16.jpg",
	"Дандасана":                  "17.jpg",
	"Чатуранга дандасана":        "18.jpg",
	"Джану ширшасана":            "19.jpg",
	"Шавасана":                   "20.jpg",
	"Дханурасана":                "21.jpg",
	"Супта конасана":             "22.jpg",
	"Падахастасана":              "23.jpg",
	"Париврита паршваконасана":   "24.jpg",
	"Париврита триконасана":      "25.jpg",
	"Паригхасана":                "26.jpg",
	"Паршвоттанасана":            "27.jpg",
	"Пурвоттанасана":             "28.jpg",
	"Матсиасана":                 "29.jpg",
	"Вирасана":                   "30.jpg",
	"Тадасана самастхити":        "31.jpg",
	"Саламба сарвангасана":       "32.jpg",
	"Гарудасана":                 "33.jpg",
	"Трианг мукхаикапада пашчим": "34.jpg",
	"Урдхва мукха шванасана":     "35.jpg",
	"Уттанасана":                 "36.jpg",
	"Уттхита паршваконасана":     "37.jpg",
	"Уттхита триконасана":        "38.jpg",
	"Уштрасана":                  "39.jpg",
	"Халасана":                   "40.jpg",
	"Урдхва прасарита падасана":  "41.jpg",
	"Пранамасана":                "42.jpg",
	"Парасаранасана":             "43.jpg",
	"Ашва cанчаланасана":         "44.jpg",
	"Парватанасана":              "45.jpg",
	"Аштанга намаскар":           "46.jpg",
	"Урдха хастасана":            "47.jpg",
}

//go:embed static/yoga
var yogaImages embed.FS

var allPoses = make([]string, 0, len(yogaPoses))

func init() {
	for k := range yogaPoses {
		allPoses = append(allPoses, k)
	}
}

func YogaTestTask(bot structs.Bot, _ *gorm.DB, _ models.Subscription, chatTgID int64) {
	rand.Seed(time.Now().Unix())

	poses := allPoses
	rand.Shuffle(len(poses), func(i, j int) { poses[i], poses[j] = poses[j], poses[i] })

	validPose := poses[0]

	variants := make([][]tgbotapi.InlineKeyboardButton, 0, 4)
	variants = append(variants, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(validPose, "/yoga valid")),
	)

	// берём от первого до четвёртого варианта включительно (всего три),
	// так как нулевой элемент был использован в качестве правильного ответа
	for _, p := range poses[1:4] {
		variants = append(variants, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(p, fmt.Sprintf("/yoga invalid %s", validPose)),
		))
	}

	rand.Shuffle(len(variants), func(i, j int) { variants[i], variants[j] = variants[j], variants[i] })

	imgPath := path.Join("static", "yoga", yogaPoses[validPose])
	img, _ := yogaImages.ReadFile(imgPath)

	msg := tgbotapi.NewPhotoUpload(chatTgID, tgbotapi.FileBytes{Name: yogaPoses[validPose], Bytes: img})
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(variants...)
	_, _ = bot.Send(msg)
}
