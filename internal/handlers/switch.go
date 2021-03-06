package handlers

import (
	"bytes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

var LayoutRus = [...]string{
	"ё", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "-", "=", "й", "ц", "у", "к", "е", "н", "г", "ш", "щ", "з",
	"х", "ъ", "ф", "ы", "в", "а", "п", "р", "о", "л", "д", "ж", "э", "\\", "я", "ч", "с", "м", "и", "т", "ь", "б", "ю",
	".", "Ё", "!", "\"", "№", ";", "%", ":", "?", "*", "(", ")", "_", "+", "Х", "Ъ", "Ж", "Э", "/", "Б", "Ю", ",",
}
var LayoutEng = [...]string{
	"`", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "-", "=", "q", "w", "e", "r", "t", "y", "u", "i", "o", "p",
	"[", "]", "a", "s", "d", "f", "g", "h", "j", "k", "l", ";", "'", "\\", "z", "x", "c", "v", "b", "n", "m", ",", ".",
	"/", "~", "!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "_", "+", "{", "}", ":", "\"", "|", "<", ">", "?",
}

var RusEngMap map[string]string
var EngRusMap map[string]string

// switchLetterToEng переводит букву в альтернативную раскладку
func switchLetterToEng(letter string) string {
	lowerLetter := strings.ToLower(letter)

	swLetter, ok := RusEngMap[lowerLetter]
	if ok {
		if letter != lowerLetter {
			return strings.ToUpper(swLetter)
		}
		return swLetter
	}

	return letter
}

// switchLetterToRus переводит букву в альтернативную раскладку
func switchLetterToRus(letter string) string {
	lowerLetter := strings.ToLower(letter)

	swLetter, ok := EngRusMap[lowerLetter]
	if ok {
		if letter != lowerLetter {
			return strings.ToUpper(swLetter)
		}
		return swLetter
	}

	return letter
}

// getSwitchFunc возвращает функцию для перевода текста
func getSwitchFunc(text string) func(letter string) string {
	// создаём мапу для кэша
	if len(RusEngMap) == 0 {
		RusEngMap = make(map[string]string)
		for i, symbol := range LayoutRus {
			RusEngMap[symbol] = LayoutEng[i]
		}
	}
	if len(EngRusMap) == 0 {
		EngRusMap = make(map[string]string)
		for i, symbol := range LayoutEng {
			EngRusMap[symbol] = LayoutRus[i]
		}
	}

	var rusCount int
	var engCount int
	for _, symbol := range text {
		if _, ok := RusEngMap[string(symbol)]; ok {
			rusCount++
		} else if _, ok := EngRusMap[string(symbol)]; ok {
			engCount++
		}
	}

	var switchLetterFunc func(letter string) string

	if rusCount > engCount {
		switchLetterFunc = switchLetterToEng
	} else {
		switchLetterFunc = switchLetterToRus
	}
	return switchLetterFunc
}

// cmdSwitch меняет раскладку текста
func SwitchHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if !strings.HasPrefix(update.Message.Text, "/switch") {
		return
	}

	tokens := strings.SplitN(update.Message.Text, " ", 2)

	if len(tokens) != 2 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"Пришлите текст для изменения его раскладки.\n"+
				"Пример команды:\n"+
				"/switch L;bdbc - cfvsq kexibq ,jn!)",
		)
		msg.ReplyToMessageID = update.Message.MessageID
		_, _ = bot.Send(msg)
	}

	text := tokens[1]

	switchLetterFunc := getSwitchFunc(text)

	var msgTextB bytes.Buffer
	for _, symbol := range text {
		msgTextB.WriteString(switchLetterFunc(string(symbol)))
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgTextB.String())
	msg.ReplyToMessageID = update.Message.MessageID
	_, _ = bot.Send(msg)
}
