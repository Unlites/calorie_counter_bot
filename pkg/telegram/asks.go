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
		b.SendMessage(message, b.msg.Responses.InvalidText)
		return "invalid"
	}
	if b.CheckExists(message) != "unknown_product" {
		b.SendMessage(message, b.msg.Responses.ProductExists)
		return "invalid"
	}
	b.SendMessage(message, b.msg.HowMuchCallories)
	return "ok"
}

// Setting incoming callories and added product to DB
func (b *Bot) AddFoodCallorieWaiting(message *tgbotapi.Message) string {
	if CheckDigits(message.Text) {
		b.SendMessage(message, b.msg.Responses.InvalidText)
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

	b.SendMessage(message, b.msg.Responses.ProductAdded)
	return "ok"

}

// Starting write lunch handle, asking for lunch composition
func (b *Bot) WriteLunchNoWaiting(message *tgbotapi.Message) {
	b.SendMessage(message, b.msg.Responses.StartLunch)
}

// Getting productName, find its callories in Db, increase current lunch callories or return with search result
func (b *Bot) WriteLunchStartWaiting(message *tgbotapi.Message) string {
	if CheckText(message.Text) {
		b.SendMessage(message, b.msg.Responses.InvalidText)
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
	b.SendMessage(message, fmt.Sprintf(b.msg.Responses.RestLunch, b.mc.productName, b.mc.callories))
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
		b.SendMessage(message, fmt.Sprintf(b.msg.Responses.StopLunch, countedCallories))
		return "stop"
	}
	if result := b.FindCallories(message); result != "ok" {
		return result
	}
	if err := b.db.IncreaseCurrentCallories(message.Chat.ID, b.mc.callories); err != nil {
		b.handleError(message, "count_callories_error", err)
		return "error"
	}
	b.SendMessage(message, fmt.Sprintf(b.msg.Responses.RestLunch, b.mc.productName, b.mc.callories))
	return "ok"
}

// If we got unknown_product in previous operation, user getting ask for tell callories of new product.
// AddFoodCallorieWaiting adds product to db, WriteLunchUnknownProductWaiting increased current_callories.
func (b *Bot) WriteLunchUnknownProductWaiting(message *tgbotapi.Message) string {
	if err := b.db.IncreaseCurrentCallories(message.Chat.ID, b.mc.callories); err != nil {
		b.handleError(message, "count_callories_error", err)
		return "error"
	}
	b.SendMessage(message, fmt.Sprintf(b.msg.Responses.RestLunch, b.mc.productName, b.mc.callories))
	return "ok"
}

// Method for start operation to show callories of existing product in DB, asking product name
func (b *Bot) ShowProductCalloriesNoWaiting(message *tgbotapi.Message) {
	b.SendMessage(message, b.msg.Responses.WhatProduct)
}

// Check product in DB, if exists - send its callories.
func (b *Bot) ShowProductCalloriesNameWaiting(message *tgbotapi.Message) string {
	if result := b.FindCallories(message); result != "ok" {
		if result == "unknown_product" {
			b.SendMessage(message, b.msg.Responses.ProductNotExists)
			return "invalid"
		}
		return "error"
	}
	b.SendMessage(message, fmt.Sprintf(b.msg.Responses.ProductCallories, b.mc.productName, b.mc.callories))
	return "ok"
}

// Reports

func (b *Bot) DayReport(message *tgbotapi.Message) {
	sum, avg, err := b.db.SelectDayCallories(message.Chat.ID)
	if err == nil {
		b.handleError(message, "report_error", err)
		return
	}
	b.SendMessage(message, fmt.Sprintf(b.msg.Responses.DayReport, sum, avg))
	return
}

func (b *Bot) WeekReport(message *tgbotapi.Message) {
	sum, avg, err := b.db.SelectWeekCallories(message.Chat.ID)
	if err == nil {
		b.handleError(message, "report_error", err)
		return
	}
	b.SendMessage(message, fmt.Sprintf(b.msg.Responses.WeekReport, sum, avg))
	return
}

func (b *Bot) MonthReport(message *tgbotapi.Message) {
	sum, avg, err := b.db.SelectMonthCallories(message.Chat.ID)
	if err == nil {
		b.handleError(message, "report_error", err)
		return
	}
	b.SendMessage(message, fmt.Sprintf(b.msg.Responses.MonthReport, sum, avg))
	return
}

// Show that we can do
func (b *Bot) ShowAsks(message *tgbotapi.Message) {
	b.SendMessage(message, b.msg.Responses.AskList)
}

// Using methods

func (b *Bot) FindCallories(message *tgbotapi.Message) string {
	if CheckText(message.Text) {
		b.SendMessage(message, b.msg.Responses.InvalidText)
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
