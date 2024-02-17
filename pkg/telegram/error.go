package telegram

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	errInvalidUrl   = errors.New("invalid url")
	errUnauthorized = errors.New("user is not authorized")
	errUnableToSave = errors.New("can't save url")
)

func (b *Bot) handleError(err error, message *tgbotapi.Message) {
	switch err {
	case errInvalidUrl:
		msg := tgbotapi.NewMessage(message.Chat.ID, b.config.InvalidUrl)
		b.bot.Send(msg)
	case errUnauthorized:
		msg := tgbotapi.NewMessage(message.Chat.ID, b.config.Unauthorized)
		b.bot.Send(msg)
	case errUnableToSave:
		msg := tgbotapi.NewMessage(message.Chat.ID, b.config.UnableToSave)
		b.bot.Send(msg)
	default:
		msg := tgbotapi.NewMessage(message.Chat.ID, b.config.Default)
		b.bot.Send(msg)
	}
}
