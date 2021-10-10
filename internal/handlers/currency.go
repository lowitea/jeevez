package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/structs"
	"gorm.io/gorm"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const (
	USD = "USD"
	RUB = "RUB"
	EUR = "EUR"
	GBP = "GBP"
)

var firstCurPatterns = map[string]string{
	"доллар":   USD,
	"долларов": USD,
	"доллара":  USD,
	"$":        USD,
	"рубль":    RUB,
	"рублей":   RUB,
	"рубля":    RUB,
	"₽":        RUB,
	"евро":     EUR,
	"€":        EUR,
	"фунт":     GBP,
	"фунтов":   GBP,
	"фунта":    GBP,
	"£":        GBP,
}

var secCurPatterns = map[string]string{
	"доллары": USD,
	"$":       USD,
	"рубли":   RUB,
	"₽":       RUB,
	"евро":    EUR,
	"€":       EUR,
	"фунты":   GBP,
	"£":       GBP,
}

var msgTemplate *regexp.Regexp

func init() {
	firstCurKeys := make([]string, 0, len(firstCurPatterns))
	for k := range firstCurPatterns {
		if k == "$" {
			k = `\$`
		}
		firstCurKeys = append(firstCurKeys, k)
	}

	secCurKeys := make([]string, 0, len(secCurPatterns))
	for k := range secCurPatterns {
		if k == "$" {
			k = `\$`
		}
		secCurKeys = append(secCurKeys, k)
	}

	msgTemplate = regexp.MustCompile(fmt.Sprintf(
		`^(\d+)\s?(%s)\sв\s(%s)$`,
		strings.Join(firstCurKeys, "|"), strings.Join(secCurKeys, "|"),
	))
}

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
	firstElem := firstCurPatterns[firstCur]
	secElem := secCurPatterns[secCur]

	if firstElem == secElem {
		return ""
	}
	return fmt.Sprintf("%s_%s", firstElem, secElem)
}

// getMsgAllCurrencies формирует сообщение со списком всех валют
func getMsgAllCurrencies(db *gorm.DB) (msgText string, err error) {
	var curRates []models.CurrencyRate
	if result := db.Find(&curRates); result.Error != nil {
		log.Printf("getting currencies from db error: %s", result.Error)
		return "", result.Error
	}

	if len(curRates) == 0 {
		log.Printf("none rates")
		return "", errors.New("none rates")
	}

	var msgTextB bytes.Buffer
	msgTextB.WriteString("Курсы всех доступных валютных пар:\n\n")

	for _, curRate := range curRates {
		msgTextB.WriteString(curRate.Name)
		msgTextB.WriteString(fmt.Sprintf(":    %.6f", curRate.Value))
		msgTextB.WriteString("\n")
	}
	return msgTextB.String(), nil
}

// cmdCurrencyRate команда /currency_rate показывает доступные пары валют или курс по конкретной паре
func cmdCurrencyRate(update tgbotapi.Update, bot structs.Bot, db *gorm.DB) {
	args := strings.Split(update.Message.Text, " ")

	var msgText string

	// команда пришла без параметров, отправляем список валют
	if len(args) == 1 {
		msgText, err := getMsgAllCurrencies(db)
		if err != nil {
			msgText = "Я прошу прощения. Биржа не отвечает по телефону. " +
				"Попробуйте уточнить у меня список валют позднее."
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
func CurrencyConverterHandler(update tgbotapi.Update, bot structs.Bot, db *gorm.DB) {
	if strings.HasPrefix(update.Message.Text, "/currency_rate") {
		cmdCurrencyRate(update, bot, db)
		return
	}

	tokens := msgTemplate.FindStringSubmatch(update.Message.Text)

	if tokens == nil {
		return
	}

	value, _ := strconv.ParseFloat(tokens[1], 64)
	currencyPair := getCurPair(tokens[2], tokens[3])

	var result float64

	// если вернулась пустая строка, значит валюты одинаковые и нужно вернуть тоже значение, что ввели
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
