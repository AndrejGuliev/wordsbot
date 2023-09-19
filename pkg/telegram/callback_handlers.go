package telegram

import (
	"fmt"
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleCallBacks(callbackQuery *tgbotapi.CallbackQuery) error {
	callback := tgbotapi.NewCallback(callbackQuery.ID, "")
	if _, err := b.bot.Request(callback); err != nil {
		return err
	}
	switch callbackQuery.Data {
	case "packages":
		return b.handlePackagesMenuCallack(callbackQuery)
	case "random_words":
		return b.handleRandomWordsCallback(callbackQuery)
	case "add_package":
		return b.handleAddPackageCallback(callbackQuery)
	case "delete_package":
		return b.handleChoosePackagesCallback(callbackQuery, 6) // Position 6 - Delete Pocket
	case "next":
		return b.handleNextCallback(callbackQuery)
	case "stop":
		return b.stopTest(callbackQuery.Message.Chat.ID, callbackQuery.From.ID)
	case "choose_package":
		return b.handleChoosePackagesCallback(callbackQuery, 5) // Position 5 - Choose Pocket
	case "back":
		return b.handleBackButtonCallback(callbackQuery)
	default:
		position, err := b.storage.GetCurrentPosition(callbackQuery.From.ID)
		if err != nil {
			return err
		}
		switch position {
		case 5:
			return b.handleStartPocketCallback(callbackQuery)
		case 6:
			return b.handleDeletePackageCallback(callbackQuery)
		default:
			return b.backToMainMenu(callbackQuery)
		}
	}
}

func (b *Bot) handleRandomWordsCallback(callbackQuery *tgbotapi.CallbackQuery) error {
	b.storage.SetPosition(callbackQuery.From.ID, 4)
	b.nextRandomWord(callbackQuery)
	_, word, _, _ := b.storage.GetCurrentWord(callbackQuery.From.ID)
	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, word)
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handlePackagesMenuCallack(callbackQuery *tgbotapi.CallbackQuery) error {
	menuMessageID := callbackQuery.Message.MessageID
	reply := tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, menuMessageID, packagesMenuKeyboard)
	_, err := b.bot.Send(reply)
	return err
}

func (b *Bot) handleBackButtonCallback(callbackQuery *tgbotapi.CallbackQuery) error {
	position, err := b.storage.GetCurrentPosition(callbackQuery.From.ID)
	if err != nil {
		return err
	}
	if position == 5 || position == 6 {
		if err := b.handlePackagesMenuCallack(callbackQuery); err != nil {
			return err
		}
		return b.storage.SetPosition(callbackQuery.From.ID, 0)
	}

	return b.backToMainMenu(callbackQuery)
}

func (b *Bot) backToMainMenu(callbackQuery *tgbotapi.CallbackQuery) error {
	menuMessageID := callbackQuery.Message.MessageID
	reply := tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, menuMessageID, mainKeyboard)
	_, err := b.bot.Send(reply)
	return err
}

func (b *Bot) handleStartPocketCallback(callbackQuery *tgbotapi.CallbackQuery) error {
	text := fmt.Sprintln(b.messages.Responses.ChosedPackage, callbackQuery.Data)
	b.storage.SetPosition(callbackQuery.From.ID, 1)
	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, text)
	if _, err := b.bot.Send(msg); err != nil {
		return err
	}
	if err := b.storage.SetTest(callbackQuery.From.ID, callbackQuery.Data); err != nil {
		return err
	}
	_, word, _, err := b.storage.GetCurrentWord(callbackQuery.From.ID)
	if err != nil {
		return err
	}
	msg = tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, word)
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleDeletePackageCallback(callbackQuery *tgbotapi.CallbackQuery) error {
	if err := b.storage.DeletePackage(callbackQuery.From.ID, callbackQuery.Data); err != nil {
		return err
	}
	if err := b.storage.SetPosition(callbackQuery.From.ID, 0); err != nil {
		return err
	}
	text := fmt.Sprintln(b.messages.Responses.DeletedPackage, callbackQuery.Data)
	menuMessageID := callbackQuery.Message.MessageID
	reply := tgbotapi.NewEditMessageTextAndMarkup(callbackQuery.Message.Chat.ID, menuMessageID, text, mainKeyboard)
	_, err := b.bot.Send(reply)
	return err
}

func (b *Bot) handleChoosePackagesCallback(callbackQuery *tgbotapi.CallbackQuery, position int) error {
	if err := b.storage.SetPosition(callbackQuery.From.ID, position); err != nil {
		return err
	}
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
		menuMessageID := callbackQuery.Message.MessageID
		reply := tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, menuMessageID, testsKeyboard)
		_, err = b.bot.Send(reply)
		return err
	}
}

func (b *Bot) handleAddPackageCallback(callbackQuery *tgbotapi.CallbackQuery) error {
	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, b.messages.Responses.InsertPackage)
	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}
	return b.storage.SetPosition(callbackQuery.From.ID, 2)

}

func (b *Bot) nextRandomWord(callbackQuery *tgbotapi.CallbackQuery) error {
	testNames, err := b.storage.MakeTestsList(callbackQuery.From.ID)
	if err != nil {
		return err
	}
	if len(testNames) < 1 {
		msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, b.messages.Responses.NoPackages)
		b.storage.SetPosition(callbackQuery.From.ID, 0)
		_, err = b.bot.Send(msg)
		return err
	}
	rand.Seed(time.Now().Unix())
	randTest := testNames[rand.Intn(len(testNames))]
	min, max, err := b.storage.GetTestIdRange(callbackQuery.From.ID, randTest)
	if err != nil {
		return err
	}
	wordID := min + rand.Intn(max-min+1)
	err = b.storage.SetRandomWord(callbackQuery.From.ID, wordID)
	return err
}

func (b *Bot) handleNextCallback(callbackQuery *tgbotapi.CallbackQuery) error {
	if err := b.nextRandomWord(callbackQuery); err != nil {
		return err
	}
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
		return err
	}

	if err := b.storage.EndTest(userID); err != nil {
		return err
	}

	text := fmt.Sprintf("%s %d", b.messages.Responses.TestDone, currentAnswNum)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = mainKeyboard
	_, err = b.bot.Send(msg)
	return err
}
