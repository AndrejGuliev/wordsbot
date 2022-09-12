package telegram

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	commandStart = "start"
)

var mainKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Пакеты слов", "Пакеты слов"),
		tgbotapi.NewInlineKeyboardButtonData("Случайные слова", "Случайные слова"),
	),
)

var randomWordsKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Далее", "Далее"),
		tgbotapi.NewInlineKeyboardButtonData("Стоп", "Стоп"),
	),
)
var pocketsMenuKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Выбрать пакет", "Выбрать пакет")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Добавить пакет", "Добавить пакет")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Удалить пакет", "Удалить пакет")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("<<", "<<")),
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:

		return b.handleStartCommand(message)
	default:

		return b.handleUnknownCommand(message)
	}

}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.Start)
	msg.ReplyMarkup = mainKeyboard
	b.storage.AddUser(message.From.ID)
	messageConfig, err := b.bot.Send(msg)
	b.storage.SetMenuMessageID(message.From.ID, messageConfig.MessageID)
	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Errors.UnknownCommand)
	_, err := b.bot.Send(msg)

	return err
}

func (b *Bot) handleMessages(message *tgbotapi.Message) error {
	position, err := b.storage.GetCurrentPosition(message.From.ID)
	if err != nil {
		return err
	}
	switch position {
	case 0: //Standart Position
		return b.handleUnknownMessages(message)
	case 1: //Chosen Test
		return b.handleUswers(message)
	case 2: //Test Words Add
		return b.handleNewWordList(message)
	case 3: //Test Name Add
		return b.handleTestName(message)
	case 4: //Random Answers
		return b.handleRandomUswers(message)
	default:
		return b.handleUnknownMessages(message)
	}
}

func (b *Bot) handleUnknownMessages(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Errors.StartLesson)
	msg.ReplyMarkup = mainKeyboard
	messageConfig, err := b.bot.Send(msg)
	b.storage.SetMenuMessageID(message.From.ID, messageConfig.MessageID)

	return err
}

func (b *Bot) handleUswers(message *tgbotapi.Message) error {
	msg, err := b.checkUswers(message)
	if err != nil {
		return err
	}
	b.bot.Send(msg)

	err = b.storage.EncCurrentWordNum(message.From.ID)
	if err != nil {
		return err
	}
	wordID, word, _, err := b.storage.GetCurrentWord(message.From.ID)
	if err != nil {
		return err
	}
	if wordID == 0 {
		return b.stopTest(message.Chat.ID, message.From.ID)
	}

	msg = tgbotapi.NewMessage(message.Chat.ID, word)
	_, err = b.bot.Send(msg)
	return err

}

func (b *Bot) handleCallBacks(callbackQuery *tgbotapi.CallbackQuery) error {
	callback := tgbotapi.NewCallback(callbackQuery.ID, callbackQuery.Data)
	if _, err := b.bot.Request(callback); err != nil {
		return err
	}
	switch callbackQuery.Data {
	case "Пакеты слов":
		return b.handlePocketsMenuCallack(callbackQuery)
	case "Случайные слова":
		return b.handleRandomWordsCallback(callbackQuery)
	case "Добавить пакет":
		return b.handleAddTestCallback(callbackQuery)
	case "Удалить пакет":
		return b.handleChoosePocketsCallback(callbackQuery, 6) // Position 6 - Delete Pocket
	case "Далее":
		return b.handleNextCallback(callbackQuery)
	case "Стоп":
		return b.stopTest(callbackQuery.Message.Chat.ID, callbackQuery.From.ID)
	case "Выбрать пакет":
		return b.handleChoosePocketsCallback(callbackQuery, 5) // Position 5 - Choose Pocket
	case "<<":
		return b.handleBackButtonCallback(callbackQuery)
	default:
		return b.handlePocketsCallback(callbackQuery)
	}

}

func (b *Bot) handleRandomWordsCallback(callbackQuery *tgbotapi.CallbackQuery) error {
	b.storage.SetPosition(callbackQuery.From.ID, 4)
	b.nextRandomWord(callbackQuery.From.ID)
	_, word, _, _ := b.storage.GetCurrentWord(callbackQuery.From.ID)
	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, word)
	_, err := b.bot.Send(msg)
	return err

}
func (b *Bot) handlePocketsMenuCallack(callbackQuery *tgbotapi.CallbackQuery) error {
	MenuMessageID, err := b.storage.GetMenuMessageID(callbackQuery.From.ID)
	if err != nil {
		return err
	}
	reply := tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, MenuMessageID, pocketsMenuKeyboard)
	_, err = b.bot.Send(reply)
	return err
}

func (b *Bot) handleBackButtonCallback(callbackQuery *tgbotapi.CallbackQuery) error {
	position, err := b.storage.GetCurrentPosition(callbackQuery.From.ID)
	if err != nil {
		return err
	}
	if position == 5 || position == 6 {
		if err := b.handlePocketsMenuCallack(callbackQuery); err != nil {
			return err
		}
		return b.storage.SetPosition(callbackQuery.From.ID, 0)
	} else {
		return b.backToMainMenu(callbackQuery)
	}
}

func (b *Bot) backToMainMenu(callbackQuery *tgbotapi.CallbackQuery) error {
	MenuMessageID, err := b.storage.GetMenuMessageID(callbackQuery.From.ID)
	if err != nil {
		return err
	}
	reply := tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, MenuMessageID, mainKeyboard)
	_, err = b.bot.Send(reply)
	return err
}

func (b *Bot) handleStartPocketCallback(callbackQuery *tgbotapi.CallbackQuery) error {
	text := fmt.Sprintln("Выбран пакет: ", callbackQuery.Data)
	b.storage.SetPosition(callbackQuery.From.ID, 1)
	fmt.Println(callbackQuery.Message.From.ID)
	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, text)
	b.bot.Send(msg)
	b.storage.SetTest(callbackQuery.From.ID, callbackQuery.Data)
	_, word, _, _ := b.storage.GetCurrentWord(callbackQuery.From.ID)
	msg = tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, word)
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handlePocketsCallback(callbackQuery *tgbotapi.CallbackQuery) error {
	position, err := b.storage.GetCurrentPosition(callbackQuery.From.ID)
	if err != nil {
		return err
	}
	switch position {
	case 5:
		return b.handleStartPocketCallback(callbackQuery)
	case 6:
		return b.handleDeletePocketCallback(callbackQuery)
	default:
		return b.backToMainMenu(callbackQuery)
	}
}

func (b *Bot) handleDeletePocketCallback(callbackQuery *tgbotapi.CallbackQuery) error {
	b.storage.DeletePocket(callbackQuery.From.ID, callbackQuery.Data)
	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, b.messages.Responses.ChoosePackageToDelete)
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleChoosePocketsCallback(callbackQuery *tgbotapi.CallbackQuery, position int) error {
	b.storage.SetPosition(callbackQuery.From.ID, position)
	testNames, err := b.storage.MakeTestsList(callbackQuery.From.ID)
	if err != nil {
		return err
	}
	if len(testNames) < 1 {
		msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, b.messages.Responses.NoPackages)
		b.storage.SetPosition(callbackQuery.From.ID, 1)
		_, err = b.bot.Send(msg)
		return err
	} else {
		testsKeyboard := b.makeTestsKeyboard(testNames)
		MenuMessageID, err := b.storage.GetMenuMessageID(callbackQuery.From.ID)
		if err != nil {
			return err
		}
		reply := tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, MenuMessageID, testsKeyboard)
		_, err = b.bot.Send(reply)
		return err
	}
}

func (b *Bot) makeTestsKeyboard(testNames []string) tgbotapi.InlineKeyboardMarkup {
	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, name := range testNames {
		button := tgbotapi.NewInlineKeyboardButtonData(name, name)
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{button})
	}
	button := tgbotapi.NewInlineKeyboardButtonData("<<", "<<")
	buttons = append(buttons, []tgbotapi.InlineKeyboardButton{button})
	testsKeyboard := tgbotapi.InlineKeyboardMarkup{InlineKeyboard: buttons}

	return testsKeyboard
}

func (b *Bot) handleAddTestCallback(callbackQuery *tgbotapi.CallbackQuery) error {
	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, b.messages.Responses.InsertPackage)
	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}
	return b.storage.SetPosition(callbackQuery.From.ID, 2)

}

func (b *Bot) handleNewWordList(message *tgbotapi.Message) error {
	pairs := strings.Split(message.Text, "\n")
	fmt.Print(pairs)
	if len(pairs) <= 4 {
		msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Errors.SmallPackage)
		_, err := b.bot.Send(msg)
		return err
	}
	for _, v := range pairs {
		pair := strings.Split(v, ":")
		if len(pair) == 2 && strings.TrimSpace(pair[0]) != "" && strings.TrimSpace(pair[1]) != "" {
			b.storage.AddNewPair(message.From.ID, pair)
		} else {
			msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Errors.EmptyStrings)
			_, err := b.bot.Send(msg)
			return err
		}
	}
	b.storage.SetPosition(message.From.ID, 3)
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.InsertPackageName)
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleTestName(message *tgbotapi.Message) error {
	if exist, err := b.storage.ValidateName(message.From.ID, message.Text); err != nil {
		return err
	} else if exist {
		msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Errors.AlredyExist)
		_, err := b.bot.Send(msg)
		return err
	}
	b.storage.AddNewTestName(message.From.ID, message.Text)
	b.storage.SetPosition(message.From.ID, 0)
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.AddedPackage)
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) nextRandomWord(userID int64) error {
	fmt.Println(userID)
	testNames, err := b.storage.MakeTestsList(userID)
	if len(testNames) < 1 {
		return err
	}
	if err != nil {
		return err
	}
	fmt.Println(testNames)
	rand.Seed(time.Now().Unix())
	fmt.Println(len(testNames))
	randTest := testNames[rand.Intn(len(testNames))]
	min, max, err := b.storage.TestIdRange(userID, randTest)
	if err != nil {
		return err
	}
	wordID := min + rand.Intn(max-min+1)
	err = b.storage.SetRandomWord(userID, wordID)
	return err
}

func (b *Bot) checkUswers(message *tgbotapi.Message) (tgbotapi.MessageConfig, error) {
	_, _, translation, err := b.storage.GetCurrentWord(message.From.ID)
	if strings.EqualFold(message.Text, translation) {
		b.storage.EncCurrentAnswNum(message.From.ID)
		msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.WordDone)
		return msg, err
	} else {
		text := fmt.Sprintf("%s %s", b.messages.Responses.WordMiss, translation)
		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		return msg, err
	}

}

func (b *Bot) handleRandomUswers(message *tgbotapi.Message) error {
	msg, err := b.checkUswers(message)
	if err != nil {
		return err
	}
	msg.ReplyMarkup = randomWordsKeyboard
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleNextCallback(callbackQuery *tgbotapi.CallbackQuery) error {
	b.nextRandomWord(callbackQuery.From.ID)
	b.storage.GetCurrentWord(callbackQuery.From.ID)
	_, word, _, err := b.storage.GetCurrentWord(callbackQuery.From.ID)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, word)
	_, err = b.bot.Send(msg)
	return err

}

func (b *Bot) stopTest(chatID int64, userID int64) error {
	currentAnswNum, err := b.storage.GetCurrentAnswNum(userID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	b.storage.SetPosition(userID, 0)
	err = b.storage.EndTest(userID)
	if err != nil {
		return err
	}
	text := fmt.Sprintf("%s %d", b.messages.Responses.TestDone, currentAnswNum)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = mainKeyboard
	messageConfig, err := b.bot.Send(msg)
	b.storage.SetMenuMessageID(userID, messageConfig.MessageID)
	return err
}

/*func (b *Bot) doesntWork(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Errors.DoesntWork)
	_, err := b.bot.Send(msg)
	return err
}*/
