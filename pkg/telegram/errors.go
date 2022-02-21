package telegram

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) handleError(message *tgbotapi.Message, typeError string, err error) {
	log.Print(err)
	switch typeError {
	case "failed_start":
		b.SendMessage(message, b.msg.Errors.FailedStart)
	case "user_param_error":
		b.SendMessage(message, b.msg.Errors.UserParamError)
	case "insert_food_error":
		b.SendMessage(message, b.msg.Errors.InsertFoodError)
	case "count_callories_error":
		b.SendMessage(message, b.msg.Errors.CountCalloriesError)
	case "report_error":
		if strings.Contains(err.Error(), "NULL") {
			b.SendMessage(message, b.msg.Errors.ReportErrorNotFound)
			return
		}
		b.SendMessage(message, b.msg.Errors.ReportErrorFailed)
	case "find_callories_error":
		b.SendMessage(message, b.msg.Errors.FindCalloriesError)
	case "find_product_name_error":
		b.SendMessage(message, b.msg.Errors.FindProductNameError)
	}

}
