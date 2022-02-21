package telegram

import (
	"log"

	"github.com/Unlites/callorie_counter/pkg/config"
	"github.com/Unlites/callorie_counter/pkg/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	bot *tgbotapi.BotAPI
	mc  *MessageComponents
	db  *db.Db
	msg config.Messages
}

type MessageComponents struct {
	waiting     string
	productName string
	callories   string
}

func NewMessageComponents(waiting string, productName string, callories string) *MessageComponents {
	return &MessageComponents{waiting: waiting, productName: productName, callories: callories}
}

func NewBot(bot *tgbotapi.BotAPI, mc *MessageComponents, db *db.Db, msg config.Messages) *Bot {
	return &Bot{bot: bot, mc: mc, db: db, msg: msg}
}

func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
		return err
	}
	for update := range updates {
		// Ignore any non-Message Updates
		if update.Message == nil {
			continue
		}
		// Command handler
		if update.Message.IsCommand() {
			b.handleCommand(update.Message)
			continue
		}
		// Message handler
		b.handleMessage(update.Message)
	}
	return nil
}
