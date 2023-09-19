package telegram

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// handleMessages handles user messages based on the user's current position.
func (b *Bot) handleMessages(message *tgbotapi.Message) error {
	position, err := b.storage.GetCurrentPosition(message.From.ID)
	if err != nil {
		return err
	}
	switch position {
	case 0: // Standard Position
		return b.handleUnknownMessages(message)
	case 1: // Chosen Test
		return b.handleAnswers(message)
	case 2: // Add package words
		return b.handleNewWordList(message)
	case 3: // Add package name
		return b.handlePackageName(message)
	case 4: // Random Answers
		return b.handleRandomAnswers(message)
	default:
		return b.handleUnknownMessages(message)
	}
}

// handleUnknownMessages handles messages when the user is in an unknown state or the message is not recognized.
func (b *Bot) handleUnknownMessages(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Errors.StartLesson)
	msg.ReplyMarkup = mainKeyboard
	_, err := b.bot.Send(msg)
	return err
}

// handleAnswers handles user answers during the test.
func (b *Bot) handleAnswers(message *tgbotapi.Message) error {
	msg, err := b.checkAnswers(message)
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

// handleRandomAnswers handles user answers during the random answers test.
func (b *Bot) handleRandomAnswers(message *tgbotapi.Message) error {
	msg, err := b.checkAnswers(message)
	if err != nil {
		return err
	}
	msg.ReplyMarkup = randomWordsKeyboard

	_, err = b.bot.Send(msg)
	return err
}

// handleNewWordList handles user input when adding new word pairs to a test.
func (b *Bot) handleNewWordList(message *tgbotapi.Message) error {
	// Split the user's input into pairs and check their validity.
	pairs := strings.Split(message.Text, "\n")
	var words [][]string
	if len(pairs) <= 4 {
		// If there are too few pairs, send an error message.
		msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Errors.SmallPackage)
		_, err := b.bot.Send(msg)
		return err
	}
	for _, v := range pairs {
		pair := strings.Split(v, ":")
		if len(pair) == 2 && strings.TrimSpace(pair[0]) != "" && strings.TrimSpace(pair[1]) != "" {
			// If the pair is valid, add it to the storage.
			words = append(words, pair)
		} else {
			// If the pair is invalid, send an error message.
			msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Errors.EmptyStrings)
			_, err := b.bot.Send(msg)
			return err
		}
	}
	// Set the user's position to indicate the need for a test name.
	b.storage.SetPosition(message.From.ID, 3)

	// Send a message requesting the test name.
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.InsertPackageName)
	_, err := b.bot.Send(msg)
	return err
}

// handlePackageName handles user input when adding a new test name.
func (b *Bot) handlePackageName(message *tgbotapi.Message) error {
	// Check if the test name already exists.
	if exist, err := b.storage.ValidateName(message.From.ID, message.Text); err != nil {
		return err
	} else if exist {
		// If the test name exists, send an error message.
		msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Errors.AlreadyExist)
		_, err := b.bot.Send(msg)
		return err
	}
	// Add the new test name to the storage.
	b.storage.AddNewTestName(message.From.ID, message.Text)
	// Set the user's position to the standard position.
	b.storage.SetPosition(message.From.ID, 0)

	// Send a message confirming the addition of the test name.
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.AddedPackage)
	_, err := b.bot.Send(msg)
	return err
}

// checkAnswers checks the user's answers and generates a response message based on correctness.
func (b *Bot) checkAnswers(message *tgbotapi.Message) (tgbotapi.MessageConfig, error) {
	// Retrieve the current word and its translation from the storage.
	_, _, translation, err := b.storage.GetCurrentWord(message.From.ID)
	if strings.EqualFold(message.Text, translation) {
		// If the answer is correct, increment the answer count and send a success message.
		b.storage.EncCurrentAnswNum(message.From.ID)
		msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.WordDone)
		return msg, err
	} else {
		// If the answer is incorrect, send an error message with the correct translation.
		text := fmt.Sprintf("%s %s", b.messages.Responses.WordMiss, translation)
		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		return msg, err
	}
}
