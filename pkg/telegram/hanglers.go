package telegram

import (
	"log"
	"reflect"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) handleMessage(message *tgbotapi.Message) string {
	if reflect.TypeOf(message.Text).Kind() != reflect.String {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Я понимаю только текст.")
		b.bot.Send(msg)
		return b.reset()
	}
	message.Text = strings.ToLower(message.Text)
	switch b.mc.Waiting {
	case "no_waiting":
		b.mc.Waiting = b.callAsk(message)
		return b.mc.Waiting
	case "waiting_product_name":
		if result := b.AddFoodProductNameWaiting(message); result != "ok" {
			return b.reset()
		}
		b.mc.Waiting = "waiting_callorie"
		return b.mc.Waiting
	case "waiting_callorie":
		if result := b.AddFoodCallorieWaiting(message); result != "ok" {
			return b.reset()
		}
		return b.reset()
	case "waiting_composition_start":
		if result := b.WriteLunchStartWaiting(message); result != "ok" {
			if result == "unknown_product" {
				b.SendMessage(message, "Сколько в этом продукте киллокаллорий?")
				b.mc.Waiting = "waiting_composition_new_product_callorie"
				return b.mc.Waiting
			}
			return b.reset()
		}
		b.mc.Waiting = "waiting_composition_rest"
		return b.mc.Waiting
	case "waiting_composition_rest":
		if result := b.WriteLunchRestWaiting(message); result != "ok" {
			if result == "unknown_product" {
				b.SendMessage(message, "Сколько в этом продукте киллокалорий?")
				b.mc.Waiting = "waiting_composition_new_product_callorie"
				return b.mc.Waiting
			}
			return b.reset()
		}
		b.mc.Waiting = "waiting_composition_rest"
		return b.mc.Waiting
	case "waiting_composition_new_product_callorie":
		if result := b.AddFoodCallorieWaiting(message); result != "ok" {
			return b.reset()
		}
		if result := b.WriteLunchUnknownProductWaiting(message); result != "ok" {
			return b.reset()
		}
		b.mc.Waiting = "waiting_composition_rest"
		return b.mc.Waiting
	case "waiting_product_to_show":
		b.ShowProductCalloriesNameWaiting(message)
		return b.reset()
	}
	return b.reset()
}

func (b *Bot) callAsk(message *tgbotapi.Message) string {
	if strings.Contains(message.Text, "добавь продукт в базу") {
		b.AddFoodNoWaiting(message)
		return "waiting_product_name"
	}
	if strings.Contains(message.Text, "зафиксируй прием пищи") {
		b.WriteLunchNoWaiting(message)
		return "waiting_composition_start"
	}
	if strings.Contains(message.Text, "сколько калорий в продукте") {
		b.ShowProductCalloriesNoWaiting(message)
		return "waiting_product_to_show"
	}
	if strings.Contains(message.Text, "отчет за день") {
		b.DayReport(message)
		return b.reset()
	}
	if strings.Contains(message.Text, "отчет за неделю") {
		b.WeekReport(message)
		return b.reset()
	}
	if strings.Contains(message.Text, "отчет за месяц") {
		b.MonthReport(message)
		return b.reset()
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, "Не понял. Можешь запустить команду /help, чтобы узнать, что я умею.")
	b.bot.Send(msg)
	return "no_waiting"
}

func (b *Bot) handleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case "help":
		b.ShowAsks(message)
	}
}

func (b *Bot) reset() string {
	b.mc.ProductName = ""
	b.mc.Callories = ""
	b.mc.Waiting = "no_waiting"
	if err := b.db.ResetCurrentCallories(); err != nil {
		log.Print(err)
	}
	return b.mc.Waiting
}
