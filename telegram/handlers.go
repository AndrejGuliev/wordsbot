package telegram

import (
	"fmt"
	"strings"

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
	_, err := b.bot.Send(msg)

	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Errors.UnknownCommand)
	_, err := b.bot.Send(msg)

	return err
}

func (b *Bot) handleMessages(message *tgbotapi.Message) error {
	position, err := b.storage.CurrentPosition(message.From.ID)
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
	default:
		return b.handleUnknownMessages(message)
	}
}

func (b *Bot) handleUnknownMessages(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Errors.StartLesson)
	msg.ReplyMarkup = mainKeyboard
	_, err := b.bot.Send(msg)

	return err
}

func (b *Bot) handleUswers(message *tgbotapi.Message) error {
	wordID, _, translation, err := b.storage.CurrentWord(message.From.ID)
	if err != nil {
		return err
	}
	if message.Text == translation {
		msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.WordDone)
		b.bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.WordMiss)
		b.bot.Send(msg)
	}

	err = b.storage.EncCurrentWordNum(message.From.ID, wordID)
	if err != nil {
		return err
	}
	wordID, word, _, err := b.storage.CurrentWord(message.From.ID)
	if err != nil {
		return err
	}
	if wordID == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.TestDone)
		b.bot.Send(msg)
		b.storage.SetPosition(message.From.ID, 0)
		err = b.storage.EndTest(message.From.ID)
		return err
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, word)
	b.bot.Send(msg)
	return nil

}

func (b *Bot) handleCallBacks(callbackQuery *tgbotapi.CallbackQuery) error {
	callback := tgbotapi.NewCallback(callbackQuery.ID, callbackQuery.Data)
	if _, err := b.bot.Request(callback); err != nil {
		return err
	}
	switch callbackQuery.Data {
	case "Пакеты слов":
		return b.handleChoosePocketsCallback(callbackQuery)
	case "Случайные слова":
		return b.doesntWork(callbackQuery.Message)
	case "Добавить пакет":
		return b.handleAddTestCallback(callbackQuery)
	default:
		return b.handlePocketsCallback(callbackQuery)
	}

}

func (b *Bot) handlePocketsCallback(callbackQuery *tgbotapi.CallbackQuery) error {
	text := fmt.Sprintln("Выбран пакет: ", callbackQuery.Data)
	b.storage.SetPosition(callbackQuery.From.ID, 1)
	fmt.Println(callbackQuery.Message.From.ID)
	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, text)
	b.bot.Send(msg)
	b.storage.SetTest(callbackQuery.From.ID, callbackQuery.Data)
	_, word, _, _ := b.storage.CurrentWord(callbackQuery.From.ID)
	msg = tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, word)
	b.bot.Send(msg)
	return nil

}

func (b *Bot) handleChoosePocketsCallback(callbackQuery *tgbotapi.CallbackQuery) error {
	var text string
	testNames, err := b.storage.MakeTestsList(callbackQuery.From.ID)
	if len(testNames) < 1 {
		text = b.messages.Responses.NoPackages
	} else {
		text = b.messages.Responses.CoosePackage
	}
	if err != nil {
		return err
	}
	testsKeyboard := b.makeTestsKeyboard(testNames)
	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, text)
	msg.ReplyMarkup = testsKeyboard
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) makeTestsKeyboard(testNames []string) tgbotapi.InlineKeyboardMarkup {
	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, name := range testNames {
		button := tgbotapi.NewInlineKeyboardButtonData(name, name)
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{button})
	}
	button := tgbotapi.NewInlineKeyboardButtonData("Добавить пакет", "Добавить пакет")
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
		b.bot.Send(msg)
		return nil
	}
	for _, v := range pairs {
		pair := strings.Split(v, ":")
		if len(pair) == 2 && strings.TrimSpace(pair[0]) != "" && strings.TrimSpace(pair[1]) != "" {
			b.storage.AddNewPair(message.From.ID, pair)
		} else {
			msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Errors.EmptyStrings)
			b.bot.Send(msg)
			return nil
		}
	}
	b.storage.SetPosition(message.From.ID, 3)
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.InsertPackageName)
	b.bot.Send(msg)
	return nil
}

func (b *Bot) handleTestName(message *tgbotapi.Message) error {
	if exist, err := b.storage.ValidateName(message.From.ID, message.Text); err != nil {
		return err
	} else if exist {
		msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Errors.AlredyExist)
		b.bot.Send(msg)
		return nil
	}
	b.storage.AddNewTestName(message.From.ID, message.Text)
	b.storage.SetPosition(message.From.ID, 0)
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.AddedPackage)
	b.bot.Send(msg)
	return nil
}

func (b *Bot) doesntWork(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Errors.DoesntWork)
	_, err := b.bot.Send(msg)
	return err
}
