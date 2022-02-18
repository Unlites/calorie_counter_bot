package telegram

import (
	"errors"
	"log"
	"reflect"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	if reflect.TypeOf(message.Text).Kind() != reflect.String {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Я понимаю только текст.")
		b.bot.Send(msg)
		b.reset(message)
		return
	}
	message.Text = strings.ToLower(message.Text)
	err := errors.New("")
	b.mc.Waiting, err = b.db.Waiting(message.Chat.ID)
	if err != nil {
		log.Print(err)
		b.SendMessage(message, "Не удалось обработать сообщение, попробуй позже.")
	}
	switch b.mc.Waiting {
	case "no_waiting":
		b.callAsk(message)
		return

	case "waiting_product_name":
		if result := b.AddFoodProductNameWaiting(message); result != "ok" {
			b.reset(message)
			return
		}
		if err := b.db.SetWaiting(message.Chat.ID, "waiting_callorie"); err != nil {
			log.Print(err)
		}
		return
	case "waiting_callorie":
		if result := b.AddFoodCallorieWaiting(message); result != "ok" {
			b.reset(message)
			return
		}
		b.reset(message)
		return
	case "waiting_composition_start":
		if result := b.WriteLunchStartWaiting(message); result != "ok" {
			if result == "unknown_product" {
				b.SendMessage(message, "Сколько в этом продукте киллокаллорий?")
				if err := b.db.SetWaiting(message.Chat.ID, "waiting_composition_new_product_callorie"); err != nil {
					log.Print(err)
				}
				return
			}
			b.reset(message)
			return
		}
		if err := b.db.SetWaiting(message.Chat.ID, "waiting_composition_rest"); err != nil {
			log.Print(err)
		}
		return
	case "waiting_composition_rest":
		if result := b.WriteLunchRestWaiting(message); result != "ok" {
			if result == "unknown_product" {
				b.SendMessage(message, "Сколько в этом продукте киллокалорий?")
				if err := b.db.SetWaiting(message.Chat.ID, "waiting_composition_new_product_callorie"); err != nil {
					log.Print(err)
				}
				return
			}
			b.reset(message)
			return
		}
		if err := b.db.SetWaiting(message.Chat.ID, "waiting_composition_rest"); err != nil {
			log.Print(err)
		}
		return
	case "waiting_composition_new_product_callorie":
		if result := b.AddFoodCallorieWaiting(message); result != "ok" {
			b.reset(message)
			return
		}
		if result := b.WriteLunchUnknownProductWaiting(message); result != "ok" {
			b.reset(message)
			return
		}
		if err := b.db.SetWaiting(message.Chat.ID, "waiting_composition_rest"); err != nil {
			log.Print(err)
		}
		return
	case "waiting_product_to_show":
		b.ShowProductCalloriesNameWaiting(message)
		b.reset(message)
		return
	}
	b.reset(message)
	return
}

func (b *Bot) callAsk(message *tgbotapi.Message) {
	if strings.Contains(message.Text, "добавь продукт в базу") {
		b.AddFoodNoWaiting(message)
		b.db.SetWaiting(message.Chat.ID, "waiting_product_name")
		return
	}
	if strings.Contains(message.Text, "зафиксируй прием пищи") {
		b.WriteLunchNoWaiting(message)
		b.db.SetWaiting(message.Chat.ID, "waiting_composition_start")
		return
	}
	if strings.Contains(message.Text, "сколько калорий в продукте") {
		b.ShowProductCalloriesNoWaiting(message)
		b.db.SetWaiting(message.Chat.ID, "waiting_product_to_show")
		return
	}
	if strings.Contains(message.Text, "отчет за день") {
		b.DayReport(message)
		b.reset(message)
		return
	}
	if strings.Contains(message.Text, "отчет за неделю") {
		b.WeekReport(message)
		b.reset(message)
		return
	}
	if strings.Contains(message.Text, "отчет за месяц") {
		b.MonthReport(message)
		b.reset(message)
		return
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, "Не понял. Можешь запустить команду /help, чтобы узнать, что я умею.")
	b.bot.Send(msg)
	b.reset(message)
	return
}

func (b *Bot) handleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case "help":
		b.ShowAsks(message)
	}
}

func (b *Bot) reset(message *tgbotapi.Message) {
	if err := b.db.SetProductName(message.Chat.ID, ""); err != nil {
		log.Print(err)
	}
	if err := b.db.SetCallories(message.Chat.ID, ""); err != nil {
		log.Print(err)
	}
	if err := b.db.SetWaiting(message.Chat.ID, "no_waiting"); err != nil {
		log.Print(err)
	}
	if err := b.db.ResetCurrentCallories(); err != nil {
		log.Print(err)
	}
	return
}
