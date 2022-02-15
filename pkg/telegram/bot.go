package telegram

import (
	"github.com/Unlites/callorie_counter/pkg/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	bot *tgbotapi.BotAPI
	mc  *MessageComponents
	db  *db.Db
}

type MessageComponents struct {
	Waiting     string
	ProductName string
	Callories   string
}

func NewMessageComponents(Waiting string, ProductName string, Callories string) *MessageComponents {
	return &MessageComponents{Waiting: "no_waiting", ProductName: ProductName, Callories: Callories}
}

func NewBot(bot *tgbotapi.BotAPI, mc *MessageComponents, db *db.Db) *Bot {
	return &Bot{bot: bot, mc: mc, db: db}
}

func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.IsCommand() {
			b.handleCommand(update.Message)
			continue
		}

		b.handleMessage(update.Message)
	}
	return nil
}
