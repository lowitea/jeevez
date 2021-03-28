package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/models"
	"gorm.io/gorm"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// getCurrencyRate получает и возвращает валюту из базы
func getCurrencyRate(db *gorm.DB, curPair string) (float64, error) {
	var currencyRate models.CurrencyRate
	result := db.First(&currencyRate, "name = ?", curPair)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return 0, result.Error
	} else if result.Error != nil {
		log.Printf("getting currency from db error: %s", result.Error)
		return 0, result.Error
	}
	return currencyRate.Value, nil
}

// getCurPair функция принимающая введённые пользователем названия валют и возвращающая имя пары
func getCurPair(firstCur string, secCur string) string {
	var firstElem string
	var secElem string
	if firstCur == "доллар" || firstCur == "долларов" || firstCur == "доллара" {
		firstElem = "USD"
	}
	if firstCur == "рубль" || firstCur == "рублей" || firstCur == "рубля" {
		firstElem = "RUB"
	}
	if firstCur == "евро" {
		firstElem = "EUR"
	}
	if secCur == "рубли" {
		secElem = "RUB"
	}
	if secCur == "доллары" {
		secElem = "USD"
	}
	if secCur == "евро" {
		secElem = "EUR"
	}
	if firstElem == secElem {
		return ""
	}
	return fmt.Sprintf("%s_%s", firstElem, secElem)
}

// cmdCurrencyRate команда /currency_rate показывает доступные пары валют или курс по конкретной паре
func cmdCurrencyRate(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *gorm.DB) {
	args := strings.Split(update.Message.Text, " ")
	var msgText string

	// команда пришла без параметров, отправляем список валют
	if len(args) == 1 {
		var curRates []models.CurrencyRate
		result := db.Find(&curRates)
		if result.Error != nil {
			log.Printf("getting currencies from db error: %s", result.Error)
			msgText = "Я прошу прощения. Биржа не отвечает по телефону." +
				"Попробуйте уточнить у меня список валют позднее."
		} else {
			var msgTextB bytes.Buffer
			msgTextB.WriteString("Курсы всех доступных валютных пар:\n\n")

			for _, curRate := range curRates {
				msgTextB.WriteString(curRate.Name)
				msgTextB.WriteString(fmt.Sprintf(":    %.6f", curRate.Value))
				msgTextB.WriteString("\n")
			}
			msgText = msgTextB.String()
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		msg.ReplyToMessageID = update.Message.MessageID
		_, _ = bot.Send(msg)
		return
	}

	// запросили конкретную валюту, пытаемся найти её в базе и отдаём результат
	curPair := args[1]
	curRate, err := getCurrencyRate(db, curPair)

	if err != nil {
		msgText = "К сожалению, я не смог найти курс Вашей валюты. " +
			"Попробуйте проверить список доступных валют, повторив " +
			"эту команду без параметров."
	} else {
		msgText = fmt.Sprintf("%f", curRate)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	msg.ReplyToMessageID = update.Message.MessageID
	_, _ = bot.Send(msg)
}

// CurrencyConverterHandler обрабатывает команды конвертации валют
func CurrencyConverterHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *gorm.DB) {
	if strings.HasPrefix(update.Message.Text, "/currency_rate") {
		cmdCurrencyRate(update, bot, db)
		return
	}

	expMsg := regexp.MustCompile(
		`^(\d+)\s(доллар|долларов|доллара|рубль|рублей|рубля|евро)\sв\s(рубли|доллары|евро)$`)

	tokens := expMsg.FindStringSubmatch(update.Message.Text)

	if tokens == nil {
		return
	}

	value, _ := strconv.ParseFloat(tokens[1], 64)
	currencyPair := getCurPair(tokens[2], tokens[3])

	var result float64
	if currencyPair == "" {
		result = value
	} else {
		curRate, err := getCurrencyRate(db, currencyPair)
		if err != nil {
			return
		}
		result = value * curRate
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%.2f", result))
	msg.ReplyToMessageID = update.Message.MessageID
	_, _ = bot.Send(msg)
}
