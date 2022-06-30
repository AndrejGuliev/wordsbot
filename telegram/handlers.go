package telegram

import (
	//"log"
	"fmt"
	"telegrambot/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	commandStart = "start"
)

var startKeyboard = tgbotapi.NewInlineKeyboardMarkup(
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
	msg.ReplyMarkup = startKeyboard
	storage.AddUser(message.From.ID)
	_, err := b.bot.Send(msg)

	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Я не знаю такую команду")
	_, err := b.bot.Send(msg)

	return err
}

func (b *Bot) handleUnknownMessages(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Сначала начните урок")
	msg.ReplyMarkup = startKeyboard
	_, err := b.bot.Send(msg)

	return err
}

func (b *Bot) handleCallBacks(callbackQuery *tgbotapi.CallbackQuery) {
	callback := tgbotapi.NewCallback(callbackQuery.ID, callbackQuery.Data)
	fmt.Println("              ", callback, "                      ")
	if _, err := b.bot.Request(callback); err != nil {
		panic(err)
	}
	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Я понимаю кнопки")
	b.bot.Send(msg)

}
