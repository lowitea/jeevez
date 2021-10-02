package subscriptions

import (
	"bytes"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/structs"
	"github.com/lowitea/jeevez/internal/tools"
	"gorm.io/gorm"
	"log"
	"text/template"
)

// GetMessage вернуть строку для отправки в мессенджер
func GetMessage(stat models.CovidStat) string {
	type Ctx struct {
		Stat models.CovidStat
	}

	msgTplString := "" +
		"\U0001F9A0 <b>COVID-19 Статистика [{{ .Stat.HumanName }}]</b>\n" +
		"{{ .Stat.LastUpdate }}\n\n" +
		"Подтверждённые: {{ .Stat.Confirmed }} " +
		"({{if gt .Stat.ConfirmedDiff 0}}+{{end}}{{ .Stat.ConfirmedDiff }})\n" +
		"Смерти: {{ .Stat.Deaths }} " +
		"({{if gt .Stat.DeathsDiff 0}}+{{end}}{{ .Stat.DeathsDiff }})\n" +
		"Выздоровевшие: {{ .Stat.Recovered }} " +
		"({{if gt .Stat.RecoveredDiff 0}}+{{end}}{{ .Stat.RecoveredDiff }})\n" +
		"Болеющие: {{ .Stat.Active }} " +
		"({{if gt .Stat.ActiveDiff 0}}+{{end}}{{ .Stat.ActiveDiff }})\n" +
		"Летальность: {{ printf \"%.6f\" .Stat.FatalityRate }}\n\n" +
		"https://yandex.ru/covid19/stat"

	msgTpl := template.Must(
		template.New("msgTpl").Parse(msgTplString))

	msg := bytes.Buffer{}
	tools.Check(msgTpl.Execute(&msg, Ctx{stat}))

	return msg.String()
}

// CovidTask таска рассылающая статистику по ковиду
func CovidTask(bot structs.Bot, db *gorm.DB, subscr models.Subscription, chatTgID int64) {
	var stat models.CovidStat

	if result := db.First(&stat, "subscription_name = ?", subscr.Name); result.Error != nil {
		log.Printf("getting CovidStat error %s", result.Error)
		return
	}

	msg := tgbotapi.NewMessage(chatTgID, GetMessage(stat))
	msg.ParseMode = HTML
	msg.DisableNotification = true
	msg.DisableWebPagePreview = true

	_, _ = bot.Send(msg)
}
