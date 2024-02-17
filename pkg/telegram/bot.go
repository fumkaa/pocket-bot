package telegram

import (
	"database/sql"
	"log"

	"github.com/fumkaa/pocket-bot/pkg/config"
	client "github.com/fumkaa/pocket-bot/pkg/pocket/client/Authentication"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot          *tgbotapi.BotAPI
	pocketClient client.Client
	redirectUrl  string
	db           *sql.DB
	config       config.Message
}

func NewBot(bot *tgbotapi.BotAPI, pocketClient client.Client, redirectUrl string, db *sql.DB, config config.Message) *Bot {
	return &Bot{
		bot:          bot,
		pocketClient: pocketClient,
		redirectUrl:  redirectUrl,
		db:           db,
		config:       config,
	}
}

func (b *Bot) Start() {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates := b.initUpdatesChannel()
	b.handleUpdates(updates)
}

func (b *Bot) initUpdatesChannel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				b.handleError(err, update.Message)
			}
			continue
		}

		if err := b.handleMessage(update.Message); err != nil {
			b.handleError(err, update.Message)
		}
	}

}
