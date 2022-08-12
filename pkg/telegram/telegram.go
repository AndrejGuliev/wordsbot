package telegram

import (
	"log"

	"github.com/AndrejGuliev/wordsbot/pkg/config"
	"github.com/AndrejGuliev/wordsbot/pkg/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot      *tgbotapi.BotAPI
	storage  *storage.WordsBotStorage
	messages config.Messages
}

func NewBot(bot *tgbotapi.BotAPI, storage *storage.WordsBotStorage, messages config.Messages) *Bot {
	return &Bot{bot: bot, storage: storage, messages: messages}
}

func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				if err := b.handleCommand(update.Message); err != nil {
					log.Fatal(err)
				}
			} else {
				b.handleMessages(update.Message)
			}
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			continue
		} else if update.CallbackQuery != nil {
			b.handleCallBacks(update.CallbackQuery)
			continue
		}
	}
	return nil
}
