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
		b.SendMessage(message, "Не удалось запустить бота, попробуй позже")
	case "user_param_error":
		b.SendMessage(message, "Не удалось обработать сообщение, попробуй позже.")
	case "insert_food_error":
		b.SendMessage(message, "Не удалось добавить продукт в базу. Попробуй в другой раз.")
	case "count_callories_error":
		b.SendMessage(message, "Не удалось посчитать каллории этого приема пищи. Попробуй позже.")
	case "report_error":
		if strings.Contains(err.Error(), "NULL") {
			b.SendMessage(message, "Не было приемов пищи!")
			return
		}
		b.SendMessage(message, "Не удалось сформировать отчет. Попробуй позже.")
	case "find_callories_error":
		b.SendMessage(message, "Не удалось найти калорийность продукта в базе. Попробуй позже.")
	case "find_product_name_error":
		b.SendMessage(message, "Не удалось выполнить поиск продукта в базе. Попробуй позже.")
	}

}
