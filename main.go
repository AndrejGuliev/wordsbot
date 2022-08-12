package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/AndrejGuliev/wordsbot/pkg/config"
	"github.com/AndrejGuliev/wordsbot/pkg/storage"
	"github.com/AndrejGuliev/wordsbot/pkg/telegram"

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

	messages, err := config.InitCfg()
	if err != nil {
		panic(err)
	}

	tBot := telegram.NewBot(bot, storage, *messages)
	tBot.Start()

}
