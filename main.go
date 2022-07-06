package main

import (
	"database/sql"
	"log"
	"telegrambot/storage"
	"telegrambot/telegram"
	"time"

	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("5346332802:AAFdLv75U13NY8YyBMqdxcwknSbqCSWPeCI")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	db, err := sql.Open("mysql", "telegram_bot:Mq7gJX-4VpzH3@/wordsbot")
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	storage := storage.NewWordsBotStorage(db)

	tBot := telegram.NewBot(bot, storage)
	tBot.Start()

}
