package main

import (
	"database/sql"
	"log"
	"os"
	"telegrambot/storage"
	"telegrambot/telegram"
	"time"

	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("botAPIKey"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	db, err := sql.Open("mysql", os.Getenv("mysql"))
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
