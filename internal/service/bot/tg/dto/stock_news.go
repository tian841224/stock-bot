package dto

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type StockNewsMessage struct {
	Text                 string
	InlineKeyboardMarkup *tgbotapi.InlineKeyboardMarkup
}
