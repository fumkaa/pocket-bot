package telegram

import (
	"fmt"
	"log"
	"net/url"

	dbMysql "github.com/fumkaa/pocket-bot/pkg/db/MySQL"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	cmdStart = "start"
	cmdHelp  = "help"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case cmdStart:
		return b.handleStartCommand(message)
	case cmdHelp:
		return b.handleHelpCommand(message)
	default:
		return b.handleUnknowCommand(message)
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	_, err := url.ParseRequestURI(message.Text)

	if err != nil {
		return errInvalidUrl
	}

	db := dbMysql.NewDatabase(b.db)
	if !b.isUserAuthorization(message) {
		return errUnauthorized
	}
	accessToken, err := db.AccessToken(int(message.Chat.ID))
	if err != nil {
		return errUnauthorized
	}

	err = b.pocketClient.Add(message.Text, accessToken)
	if err != nil {
		return errUnableToSave
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, b.config.SavesSuccessfully)
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	if b.isUserAuthorization(message) {
		log.Print("[handleStartCommand]проверка на авторизацию прошла: юзер авторизирован")
		msg := tgbotapi.NewMessage(message.Chat.ID, b.config.AlreadyAuthorized)

		_, err := b.bot.Send(msg)
		if err != nil {
			return fmt.Errorf("[handleStartCommand]cmdStart can't send message")
		}
		return nil
	}
	log.Print("[handleStartCommand]проверка на авторизацию прошла: юзер не авторизирован")
	authLink, err := b.intitAuthorizationProcess(message)
	if err != nil {
		return fmt.Errorf("[handleStartCommand]%w", err)
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(b.config.Start, authLink))

	_, err = b.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("[handleStartCommand]cmdStart can't send message")
	}
	return nil

}

func (b *Bot) handleHelpCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.config.Help)

	_, err := b.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("[handleCommand]cmdHelp can't send message")
	}
	return nil
}

func (b *Bot) handleUnknowCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.config.UnknowCommand)
	_, err := b.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("[handleCommand]cmdHelp can't send message")
	}
	return nil
}
