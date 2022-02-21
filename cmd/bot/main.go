package main

import (
	"log"

	"github.com/Unlites/callorie_counter/pkg/config"
	"github.com/Unlites/callorie_counter/pkg/db"
	"github.com/Unlites/callorie_counter/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	database := db.InitDB(cfg.DbUser, cfg.DbPassword)

	botApi, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}
	botApi.Debug = false

	msgComponents := telegram.NewMessageComponents("no_waiting", "", "")

	bot := telegram.NewBot(botApi, msgComponents, database, cfg.Messages)
	if err := bot.Start(); err != nil {
		log.Fatal(err)
	}
}
