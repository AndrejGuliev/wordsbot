package main

import (
	"log"

	"telegrambot/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("5346332802:AAFdLv75U13NY8YyBMqdxcwknSbqCSWPeCI")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	tBot := telegram.NewBot(bot)
	tBot.Start()

}
