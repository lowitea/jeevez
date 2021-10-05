package subscriptions

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/structs"
	"gorm.io/gorm"
	"math/rand"
	"path"
	"time"
)

//goland:noinspection SpellCheckingInspection
var yogaPoses = map[string]string{
	"Джатхара Паривартанасана": "1.jpg",
	"Падангуштхасана":          "2.jpg",
	"Парипурна Навасана":       "3.jpg",
	"Упавиштха Конасана":       "4.jpg",
	"Прасарита Падоттанасана":  "5.jpg",
	"Випарита Карани":          "6.jpg",
	"Пашчимоттанасана":         "7.jpg",
}

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

	imgPath := path.Join("internal", "static", "yoga", yogaPoses[validPose])

	msg := tgbotapi.NewPhotoUpload(chatTgID, imgPath)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(variants...)
	_, _ = bot.Send(msg)
}
