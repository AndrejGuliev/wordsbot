package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var mainKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Пакеты слов", "packages"),
		tgbotapi.NewInlineKeyboardButtonData("Случайные слова", "random_words"),
	),
)

var randomWordsKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Далее", "next"),
		tgbotapi.NewInlineKeyboardButtonData("Стоп", "stop"),
	),
)
var packagesMenuKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Выбрать пакет", "choose_package")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Добавить пакет", "add_package")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Удалить пакет", "delete_package")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("<<", "back")),
)

func (b *Bot) makeTestsKeyboard(testNames []string) tgbotapi.InlineKeyboardMarkup {
	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, name := range testNames {
		button := tgbotapi.NewInlineKeyboardButtonData(name, name)
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{button})
	}
	button := tgbotapi.NewInlineKeyboardButtonData("<<", "back")
	buttons = append(buttons, []tgbotapi.InlineKeyboardButton{button})
	testsKeyboard := tgbotapi.InlineKeyboardMarkup{InlineKeyboard: buttons}

	return testsKeyboard
}
