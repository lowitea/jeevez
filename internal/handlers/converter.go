package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/config"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

func getCurrencyRate() (float64, error) {
	resp, err := http.Get("https://free.currconv.com/api/v7/convert?q=USD_RUB&compact=ultra&apiKey=d65168e35590aedbdcc5")
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

	type apiResp struct {
		UsdRub float64 `json:"USD_RUB"`
	}

	var data apiResp
	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("Error get data: %s", err)
		return 0, err
	}
	return data.UsdRub, nil
}

func CurrencyConverterHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI, _ *config.Config) {
	// выходим сразу, если сообщения нет
	if update.Message == nil {
		return
	}

	expMsg := regexp.MustCompile(
		`(\d+)\s(доллар|долларов|рубль|рублей|евро)\sв\s(рубли|доллары|евро)`)

	tokens := expMsg.FindStringSubmatch(update.Message.Text)

	if tokens == nil {
		return
	}

	value, _ := strconv.ParseFloat(tokens[1], 64)
	//firstCur := tokens[2]
	//secCur := tokens[3]

	curRate, err := getCurrencyRate()

	if err != nil {
		return
	}

	result := value * curRate

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%.2f", result))
	msg.ReplyToMessageID = update.Message.MessageID
	_, _ = bot.Send(msg)
}
