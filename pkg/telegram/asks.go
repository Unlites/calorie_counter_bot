package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Asking user for productName
func (b *Bot) AddFoodNoWaiting(message *tgbotapi.Message) {
	b.SendMessage(message, "Какой продукт ты хочешь добавить?")
}

// Setting incoming productName, asking for product callories
func (b *Bot) AddFoodProductNameWaiting(message *tgbotapi.Message) string {
	err := b.db.SetProductName(message.Chat.ID, message.Text)
	if err != nil {
		b.handleError(message, "user_param_error", err)
		return "error"
	}
	if b.mc.productName, err = b.db.ProductName(message.Chat.ID); err != nil {
		b.handleError(message, "user_param_error", err)
		return "error"
	}
	if CheckText(message.Text) {
		b.SendMessage(message, "Текст должен содержать только назавание одного продукта или блюда, без знаков препинания или спецсимволов.")
		return "invalid"
	}
	if b.CheckExists(message) != "unknown_product" {
		b.SendMessage(message, "Такой продукт уже есть в базе!")
		return "invalid"
	}
	b.SendMessage(message, "Сколько в этом продукте киллокаллорий?")
	return "ok"
}

// Setting incoming callories and added product to DB
func (b *Bot) AddFoodCallorieWaiting(message *tgbotapi.Message) string {
	if CheckDigits(message.Text) {
		b.SendMessage(message, "Текст должен содержать только цифру с количеством киллокаллорий в продукте или блюде.")
		return "invalid"
	}
	err := b.db.SetCallories(message.Chat.ID, message.Text)
	if err != nil {
		b.handleError(message, "user_param_error", err)
		return "error"
	}
	if b.mc.callories, err = b.db.Callories(message.Chat.ID); err != nil {
		b.handleError(message, "user_param_error", err)
		return "error"
	}

	if err := b.db.InsertFood(b.mc.productName, b.mc.callories); err != nil {
		b.handleError(message, "insert_food_error", err)
		return "error"
	}

	b.SendMessage(message, "Добавил продукт в базу!")
	return "ok"

}

// Starting write lunch handle, asking for lunch composition
func (b *Bot) WriteLunchNoWaiting(message *tgbotapi.Message) {
	b.SendMessage(message, "Окей, давай уточним состав приема пищи. Перечисляй продукты или блюда из базы по одному.")
}

// Getting productName, find its callories in Db, increase current lunch callories or return with search result
func (b *Bot) WriteLunchStartWaiting(message *tgbotapi.Message) string {
	if CheckText(message.Text) {
		b.SendMessage(message, "Текст должен содержать только назавание одного продукта или блюда, без знаков препинания или спецсимволов.")
		return "invalid"
	}
	result := b.FindCallories(message)
	if result == "ok" {
		return result
	}
	err := b.db.IncreaseCurrentCallories(message.Chat.ID, b.mc.callories)
	if err != nil {
		b.handleError(message, "count_callories_error", err)
		return "error"
	}
	b.SendMessage(message, fmt.Sprintf("Прибавил продукт %s с %s ккал! Есть еще продукты? Пиши его название, если да, либо \"нет\", если таковых не осталось.", b.mc.productName, b.mc.callories))
	return "ok"
}

// Same for rest products, also provides stop-word No for quit handle and add lunch to DB
func (b *Bot) WriteLunchRestWaiting(message *tgbotapi.Message) string {
	if message.Text == "нет" {
		countedCallories, err := b.db.SelectCurrentCallories(message.Chat.ID)
		if err != nil {
			b.handleError(message, "count_callories_error", err)
			return "error"
		}
		b.SendMessage(message, fmt.Sprintf("Прием пищи зафиксирован! Ты потребила %s киллокаллорий", countedCallories))
		return "stop"
	}
	if result := b.FindCallories(message); result != "ok" {
		return result
	}
	if err := b.db.IncreaseCurrentCallories(message.Chat.ID, b.mc.callories); err != nil {
		b.handleError(message, "count_callories_error", err)
		return "error"
	}
	b.SendMessage(message, fmt.Sprintf("Прибавил продукт %s с %s килокаллорий! Есть еще продукты? Пиши его название, если да, либо \"нет\", если таковых не осталось.", b.mc.productName, b.mc.callories))
	return "ok"
}

// If we got unknown_product in previous operation, user getting ask for tell callories of new product.
// AddFoodCallorieWaiting adds product to db, WriteLunchUnknownProductWaiting increased current_callories.
func (b *Bot) WriteLunchUnknownProductWaiting(message *tgbotapi.Message) string {
	if err := b.db.IncreaseCurrentCallories(message.Chat.ID, b.mc.callories); err != nil {
		b.handleError(message, "count_callories_error", err)
		return "error"
	}
	b.SendMessage(message, fmt.Sprintf("Прибавил %s продукт с %s килокаллорий! Есть еще продукты? Пиши его название, если да, либо \"нет\", если таковых не осталось.", b.mc.productName, b.mc.callories))
	return "ok"
}

// Method for start operation to show callories of existing product in DB, asking product name
func (b *Bot) ShowProductCalloriesNoWaiting(message *tgbotapi.Message) {
	b.SendMessage(message, "О каком продукте идет речь?")
}

// Check product in DB, if exists - send its callories.
func (b *Bot) ShowProductCalloriesNameWaiting(message *tgbotapi.Message) string {
	if result := b.FindCallories(message); result != "ok" {
		if result == "unknown_product" {
			b.SendMessage(message, "Такого продукта нет в базе!")
			return "invalid"
		}
		return "error"
	}
	b.SendMessage(message, fmt.Sprintf("В продукте %s %s килокалорий", b.mc.productName, b.mc.callories))
	return "ok"
}

// Reports

func (b *Bot) DayReport(message *tgbotapi.Message) {
	sum, avg, err := b.db.SelectDayCallories(message.Chat.ID)
	if err == nil {
		b.handleError(message, "report_error", err)
		return
	}
	b.SendMessage(message, fmt.Sprintf("За этот день ты потребила %s калорий. В среднем %s за прием пищи.", sum, avg))
	return
}

func (b *Bot) WeekReport(message *tgbotapi.Message) {
	sum, avg, err := b.db.SelectWeekCallories(message.Chat.ID)
	if err == nil {
		b.handleError(message, "report_error", err)
		return
	}
	b.SendMessage(message, fmt.Sprintf("За прошедшую неделю ты потребила %s калорий. В среднем %s за прием пищи.", sum, avg))
	return
}

func (b *Bot) MonthReport(message *tgbotapi.Message) {
	sum, avg, err := b.db.SelectMonthCallories(message.Chat.ID)
	if err == nil {
		b.handleError(message, "report_error", err)
		return
	}
	b.SendMessage(message, fmt.Sprintf("За прошедший месяц ты потребила %s калорий. В среднем %s за прием пищи.", sum, avg))
	return
}

// Show that we can do
func (b *Bot) ShowAsks(message *tgbotapi.Message) {
	b.SendMessage(message, "Список команд:\nДобавь продукт в базу\nЗафиксируй прием пищи\nСколько калорий в продукте\nОтчет за день\nОтчет за неделю\nОтчет за месяц")
}

// Using methods

func (b *Bot) FindCallories(message *tgbotapi.Message) string {
	if CheckText(message.Text) {
		b.SendMessage(message, "Текст должен содержать только назавание одного продукта или блюда, без знаков препинания или спецсимволов.")
		return "invalid"
	}
	err := b.db.SetProductName(message.Chat.ID, message.Text)
	if err != nil {
		b.handleError(message, "user_param_error", err)
		return "error"
	}
	if b.mc.productName, err = b.db.ProductName(message.Chat.ID); err != nil {
		b.handleError(message, "user_param_error", err)
		return "error"
	}
	if b.CheckExists(message) != "product_exists" {
		return b.CheckExists(message)
	}
	callories, err := b.db.SelectProductCallories(message.Chat.ID, b.mc.productName)
	if err != nil {
		b.handleError(message, "find_callories_error", err)
		return "error"
	}
	if err := b.db.SetCallories(message.Chat.ID, callories); err != nil {
		b.handleError(message, "user_param_error", err)
		return "error"
	}
	if b.mc.callories, err = b.db.Callories(message.Chat.ID); err != nil {
		b.handleError(message, "user_param_error", err)
		return "error"
	}
	return "ok"
}

func (b *Bot) SendMessage(message *tgbotapi.Message, text string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	b.bot.Send(msg)
}

func (b *Bot) CheckExists(message *tgbotapi.Message) string {
	isProductExists, err := b.db.SelectFood(message.Chat.ID, b.mc.productName)
	if err != nil {
		b.handleError(message, "find_product_name_error", err)
		return "error"
	}
	if b.mc.productName != isProductExists {
		return "unknown_product"
	}
	return "product_exists"
}
