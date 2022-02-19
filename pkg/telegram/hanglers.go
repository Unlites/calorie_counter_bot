package telegram

import (
	"errors"
	"reflect"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) handleMessage(message *tgbotapi.Message) {

	// Check type of message - it must be string only
	if reflect.TypeOf(message.Text).Kind() != reflect.String {
		b.SendMessage(message, "Я понимаю только текст")
		b.reset(message)
		return
	}

	// We will work with the lowercase text
	message.Text = strings.ToLower(message.Text)

	err := errors.New("")

	// Get general parameters from DB
	if b.mc.waiting, err = b.db.Waiting(message.Chat.ID); err != nil {
		b.handleError(message, "user_param_error", err)
		return
	}
	// Waiting - it is that bot waiting from user depending of operation type
	switch b.mc.waiting {
	// no_waiting - first stage of waitings
	case "no_waiting":
		// Method for set next Waiting depending of message text by user
		if err := b.callAsk(message); err != nil {
			b.handleError(message, "user_param_error", err)
		}
		return
	// Waiting for name of product to add it to DB
	case "waiting_product_name":
		if result := b.AddFoodProductNameWaiting(message); result != "ok" {
			b.reset(message)
			return
		}
		if err := b.db.SetWaiting(message.Chat.ID, "waiting_callorie"); err != nil {
			b.handleError(message, "user_param_error", err)
			return
		}
		return
	// Waiting for callorie of product to add it to DB
	case "waiting_callorie":
		if result := b.AddFoodCallorieWaiting(message); result != "ok" {
			b.reset(message)
			return
		}
		b.reset(message)
		return
	// First waiting for write lunch
	case "waiting_composition_start":
		if result := b.WriteLunchStartWaiting(message); result != "ok" {
			// If product not in DB, set waiting for start operation to add it to DB
			if result == "unknown_product" {
				b.SendMessage(message, "Сколько в этом продукте киллокаллорий?")
				if err := b.db.SetWaiting(message.Chat.ID, "waiting_composition_new_product_callorie"); err != nil {
					b.handleError(message, "user_param_error", err)
					return
				}
				return
			}
			b.reset(message)
			return
		}
		if err := b.db.SetWaiting(message.Chat.ID, "waiting_composition_rest"); err != nil {
			b.handleError(message, "user_param_error", err)
		}
		return

	// Waiting for rest of product to write lunch
	case "waiting_composition_rest":
		if result := b.WriteLunchRestWaiting(message); result != "ok" {
			// If product not in DB, set waiting for start operation to add it to DB
			if result == "unknown_product" {
				b.SendMessage(message, "Сколько в этом продукте киллокалорий?")
				if err := b.db.SetWaiting(message.Chat.ID, "waiting_composition_new_product_callorie"); err != nil {
					b.handleError(message, "user_param_error", err)
					return
				}
				return
			}
			b.reset(message)
			return
		}
		if err := b.db.SetWaiting(message.Chat.ID, "waiting_composition_rest"); err != nil {
			b.handleError(message, "user_param_error", err)
		}
		return

	// Waiting for add product to DB while writing the lunch
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
			b.handleError(message, "user_param_error", err)
		}
		return

	// Waiting for name of product to show its calorie
	case "waiting_product_to_show":
		b.ShowProductCalloriesNameWaiting(message)
		b.reset(message)
		return
	}
	b.reset(message)
	return
}

func (b *Bot) callAsk(message *tgbotapi.Message) error {
	if strings.Contains(message.Text, "добавь продукт в базу") {
		b.AddFoodNoWaiting(message)
		err := b.db.SetWaiting(message.Chat.ID, "waiting_product_name")
		if err != nil {
			return err
		}
		return nil
	}
	if strings.Contains(message.Text, "зафиксируй прием пищи") {
		b.WriteLunchNoWaiting(message)
		err := b.db.SetWaiting(message.Chat.ID, "waiting_composition_start")
		if err != nil {
			return err
		}
		return nil
	}
	if strings.Contains(message.Text, "сколько калорий в продукте") {
		b.ShowProductCalloriesNoWaiting(message)
		err := b.db.SetWaiting(message.Chat.ID, "waiting_product_to_show")
		if err != nil {
			return err
		}
		return nil
	}
	if strings.Contains(message.Text, "отчет за день") {
		b.DayReport(message)
		b.reset(message)
		return nil
	}
	if strings.Contains(message.Text, "отчет за неделю") {
		b.WeekReport(message)
		b.reset(message)
		return nil
	}
	if strings.Contains(message.Text, "отчет за месяц") {
		b.MonthReport(message)
		b.reset(message)
		return nil
	}

	b.SendMessage(message, "Не понял. Можешь запустить команду /help, чтобы узнать, что я умею.")
	b.reset(message)
	return nil
}

func (b *Bot) handleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		userExists, err := b.db.UserExists(message.Chat.ID)
		if err != nil {
			b.handleError(message, "failed_start", err)
			return
		}
		if !userExists {
			if err = b.db.AddNewUser(message.Chat.ID); err != nil {
				b.handleError(message, "failed_start", err)
				return
			}
			b.SendMessage(message, "Добро пожаловать! Запусти команду /help, чтобы узнать, что я умею.")
			return
		}
		b.SendMessage(message, "Привет! Запусти команду /help, чтобы узнать, что я умею.")
	case "help":
		b.ShowAsks(message)
	default:
		b.SendMessage(message, "Неизвестная команда.")
	}
}

// Method to stop current operation and set all parameters to default values.
func (b *Bot) reset(message *tgbotapi.Message) {
	b.db.SetProductName(message.Chat.ID, "")
	b.db.SetCallories(message.Chat.ID, "")
	b.db.SetWaiting(message.Chat.ID, "no_waiting")
	b.db.ResetCurrentCallories(message.Chat.ID)
	return
}
