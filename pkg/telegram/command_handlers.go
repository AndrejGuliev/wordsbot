package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// commandStart is a constant representing the "start" command.
const (
	commandStart = "start"
)

// handleCommand handles incoming bot commands and routes them to specific command handlers.
func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

// handleStartCommand handles the "start" command by sending a response message and adding the user to the storage.
func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	// Create a response message for the "start" command.
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.Start)
	msg.ReplyMarkup = mainKeyboard

	// Add the user to the storage (assuming this operation is successful).
	if err := b.storage.AddUser(message.From.ID); err != nil {
		return err
	}

	// Send the response message to the user.
	_, err := b.bot.Send(msg)
	return err
}

// handleUnknownCommand handles unknown commands by sending an error message.
func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	// Create an error message for unknown commands.
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Errors.UnknownCommand)

	// Send the error message to the user.
	_, err := b.bot.Send(msg)
	return err
}
