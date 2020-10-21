package handlers

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/config"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

func getCurrencyRate(curPair string) (float64, error) {
	resp, err := http.Get(fmt.Sprintf("https://free.currconv.com/api/v7/convert?"+
		"q=%s&compact=ultra&apiKey=d65168e35590aedbdcc5", curPair))
	if err != nil {
		log.Printf("Error get data: %s", err)
		return 0, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error get data: %s", err)
		return 0, err
	}

	tpl := regexp.MustCompile(`{"\w{3}_\w{3}":(\d+.\d+)}`)
	rateStr := tpl.FindStringSubmatch(string(body))[1]

	rate, err := strconv.ParseFloat(rateStr, 64)
	if err != nil {
		log.Printf("Error get data: %s", err)
		return 0, err
	}

	return rate, nil
}

func getCurPair(firstCur string, secCur string) string {
	var firstElem string
	var secElem string
	if firstCur == "доллар" || firstCur == "долларов" {
		firstElem = "USD"
	}
	if firstCur == "рубль" || firstCur == "рублей" {
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

func CurrencyConverterHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI, _ *config.Config) {
	// выходим сразу, если сообщения нет
	if update.Message == nil {
		return
	}

	expMsg := regexp.MustCompile(
		`^(\d+)\s(доллар|долларов|рубль|рублей|евро)\sв\s(рубли|доллары|евро)$`)

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
		curRate, err := getCurrencyRate(currencyPair)
		if err != nil {
			return
		}
		result = value * curRate
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%.2f", result))
	msg.ReplyToMessageID = update.Message.MessageID
	_, _ = bot.Send(msg)
}
