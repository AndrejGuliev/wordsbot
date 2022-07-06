package telegram

import (
	"fmt"

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
	msg := tgbotapi.NewMessage(message.Chat.ID, "Приветственное сообщение")
	msg.ReplyMarkup = mainKeyboard
	b.storage.AddUser(message.From.ID)
	_, err := b.bot.Send(msg)

	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Я не знаю такую команду")
	_, err := b.bot.Send(msg)

	return err
}

func (b *Bot) handleMessages(message *tgbotapi.Message) error {
	test, err := b.storage.CurrentTest(message.From.ID)
	if err != nil {
		return err
	}
	switch test {
	case "":
		return b.handleUnknownMessages(message)
	default:
		return b.handleUswers(message)
	}
}

func (b *Bot) handleUnknownMessages(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Сначала начните урок")
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
		msg := tgbotapi.NewMessage(message.Chat.ID, "✅")
		b.bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌")
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
		msg := tgbotapi.NewMessage(message.Chat.ID, "Тест закончен")
		b.bot.Send(msg)
		err = b.storage.ClearTest(message.From.ID)
		fmt.Println(err)
		return nil
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
		return b.doesntWork(callbackQuery.Message)
	default:
		return b.handlePocketsCallback(callbackQuery)
	}

}

func (b *Bot) handlePocketsCallback(callbackQuery *tgbotapi.CallbackQuery) error {
	text := fmt.Sprintln("Выбран пакет: ", callbackQuery.Data)
	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, text)
	b.bot.Send(msg)
	b.storage.UserStartTest(callbackQuery.From.ID, callbackQuery.Data)
	_, word, _, _ := b.storage.CurrentWord(callbackQuery.From.ID)
	msg = tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, word)
	b.bot.Send(msg)
	return nil

}

func (b *Bot) makeTestsKeyboard(testNames []string) tgbotapi.InlineKeyboardMarkup {
	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, name := range testNames {
		button := tgbotapi.NewInlineKeyboardButtonData(name, name)
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{button})
	}
	button := tgbotapi.NewInlineKeyboardButtonData("Добавить пакет", "Добавить пакет")
	buttons = append(buttons, []tgbotapi.InlineKeyboardButton{button})
	testsKeyboard := tgbotapi.InlineKeyboardMarkup{buttons}

	return testsKeyboard
}

func (b *Bot) handleChoosePocketsCallback(callbackQuery *tgbotapi.CallbackQuery) error {
	var text string
	testNames, err := b.storage.MakeTestsList(callbackQuery.From.ID)
	if len(testNames) < 1 {
		text = "У вас пока нет пакетов"
	} else {
		text = "Выберите пакет"
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

func (b *Bot) doesntWork(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Пока не работает")
	_, err := b.bot.Send(msg)
	return err
}
