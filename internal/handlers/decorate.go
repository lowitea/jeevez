package handlers

import (
	"bytes"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lowitea/jeevez/internal/structs"
	"github.com/wayneashleyberry/eeemo/pkg/zalgo"
	"sort"
	"strings"
)

func decorStrth(text string) string {
	var decoratedTextB bytes.Buffer
	for _, letter := range text {
		decoratedTextB.WriteString(string(letter))
		decoratedTextB.WriteString("\u0336")
	}
	return decoratedTextB.String()
}

var LettersInvertedMap = map[string]string{
	// rus
	"ё": "ǝ", "й": "ņ", "ц": "ǹ", "у": "ʎ", "к": "ʞ", "е": "ǝ", "г": "ɹ", "ш": "m", "щ": "m", "з": "ε", "ъ": "q",
	"ф": "ȸ", "ы": "ıq", "в": "ʚ", "а": "ɐ", "п": "u", "р": "d", "л": "v", "д": "ɓ", "э": "є", "я": "ʁ", "ч": "Һ",
	"с": "ɔ", "м": "w", "и": "и", "т": "ɯ", "ь": "q", "б": "ƍ", "ю": "oı",

	// eng
	"q": "ᕹ", "w": "ʍ", "e": "ǝ", "r": "ɹ", "t": "ʇ", "y": "ʎ", "u": "n", "i": "ı", "p": "d", "a": "ɐ", "d": "p",
	"f": "ɟ", "g": "ƃ", "h": "ɥ", "j": "ɾ", "k": "ʞ", "l": "ן", "c": "ɔ", "v": "ʌ", "b": "q", "n": "u",
	"m": "ɯ",

	// spec symbols
	".": "˙", "!": "¡", "?": "¿", "_": "‾", "'": ",", "\"": "„", "&": "⅋", "(": ")", ")": "(", "{": "}", "}": "{",
	"<": ">", ">": "<", "[": "]", "∴": "∵",

	// digits
	"1": "Ɩ", "2": "↊", "3": "Ɛ", "4": "ᔭ", "5": "ϛ", "6": "9", "7": "ㄥ", "9": "6",
}

func decorInvert(text string) string {
	var decoratedTextB bytes.Buffer
	lowerText := strings.ToLower(text)

	var letters []string
	for _, letter := range lowerText {
		letters = append([]string{string(letter)}, letters...)
	}

	for _, letter := range letters {
		if decoratedLetter, ok := LettersInvertedMap[letter]; ok {
			decoratedTextB.WriteString(decoratedLetter)
		} else {
			decoratedTextB.WriteString(letter)
		}
	}

	return decoratedTextB.String()
}

func decorReverse(text string) string {
	var letters []string
	for _, letter := range text {
		letters = append([]string{string(letter)}, letters...)
	}
	var decoratedTextB bytes.Buffer
	for _, letter := range letters {
		decoratedTextB.WriteString(letter)
	}
	return decoratedTextB.String()
}

func decorZalgo(text string) string {
	return zalgo.Generate(text, "maxi", true, true, true)
}

var decorateFuncMap = map[string]struct {
	Func        func(text string) string
	Description string
}{
	"strth": {
		Func:        decorStrth,
		Description: "Зачёркнутый текст",
	},
	"reverse": {
		Func:        decorReverse,
		Description: "Текст в обратную сторону",
	},
	"invert": {
		Func:        decorInvert,
		Description: "Перевёрнутый текст",
	},
	"zalgo": {
		Func:        decorZalgo,
		Description: "Зальгофикация текста",
	},
}

// DecorateTextHandler команда /decorate для разных украшательств текста
func DecorateTextHandler(update tgbotapi.Update, bot structs.Bot) {
	if !strings.HasPrefix(update.Message.Text, "/decorate") {
		return
	}

	args := strings.SplitN(update.Message.Text, " ", 3)

	if len(args) != 3 {
		var commands []string
		for cmdName := range decorateFuncMap {
			commands = append(commands, cmdName)
		}
		sort.Strings(commands)
		var listCommandsB bytes.Buffer
		for _, cmdName := range commands {
			listCommandsB.WriteString("\n  <b>")
			listCommandsB.WriteString(cmdName)
			listCommandsB.WriteString("</b> - ")
			listCommandsB.WriteString(decorateFuncMap[cmdName].Description)
		}
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			fmt.Sprintf("Чтобы попросить меня декорировать текст, отправьте команду в формате:\n"+
				"/decorate название_преобразования Ваш текст\n"+
				"Например, так:\n"+
				"/decorate strth Текст, который нужно зачеркнуть\n\n"+
				"Все варианты декорирования:\n%s", listCommandsB.String()),
		)
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ParseMode = HTML
		_, _ = bot.Send(msg)
		return
	}

	decorName, text := args[1], args[2]

	decorateConf, ok := decorateFuncMap[decorName]
	var msgText string
	if !ok {
		msgText = "К сожалению, такому меня не учили ):"
	} else {
		msgText = decorateConf.Func(text)
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	msg.ReplyToMessageID = update.Message.MessageID
	_, _ = bot.Send(msg)
}
