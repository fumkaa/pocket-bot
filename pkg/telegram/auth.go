package telegram

import (
	"fmt"
	"log"

	dbMysql "github.com/fumkaa/pocket-bot/pkg/db/MySQL"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) generateAuthorizationLink(chatId int) (string, error) {
	redirectUrl := fmt.Sprintf("%s?chat_id=%d", b.redirectUrl, chatId)
	log.Printf("[generateAuthorizationLink]redirect url: %s", redirectUrl)
	reqToken, err := b.pocketClient.RequestToken(redirectUrl)
	if err != nil {
		return "", fmt.Errorf("[generateAuthorizationLink]can't get request token: %w", err)
	}
	log.Printf("[generateAuthorizationLink]request token: %s", reqToken)
	dbase := dbMysql.NewDatabase(b.db)

	err = dbase.CreateTable(chatId)
	log.Print(err)
	if err != nil {
		return "", fmt.Errorf("[generateAuthorizationLink]can't create table: %w", err)
	}
	err = dbase.SaveRequestToken(chatId, reqToken)
	if err != nil {
		return "", fmt.Errorf("[generateAuthorizationLink]can't save request token: %w", err)
	}
	return b.pocketClient.AuthorizationUrl(redirectUrl, reqToken)
}

func (b *Bot) intitAuthorizationProcess(message *tgbotapi.Message) (string, error) {
	authLink, err := b.generateAuthorizationLink(int(message.Chat.ID))
	if err != nil {
		return "", fmt.Errorf("[intitAuthorizationProcess]can't generate link: %w", err)
	}
	return authLink, nil
}

func (b *Bot) isUserAuthorization(message *tgbotapi.Message) bool {
	db := dbMysql.NewDatabase(b.db)
	accTok, _ := db.AccessToken(int(message.Chat.ID))

	return accTok != ""
}
