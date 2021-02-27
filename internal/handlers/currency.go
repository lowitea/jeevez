package handlers

import (
	"errors"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/models"
	"gorm.io/gorm"
	"log"
	"regexp"
	"strconv"
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

func CurrencyConverterHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *gorm.DB) {
	// выходим сразу, если сообщения нет
	if update.Message == nil {
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
