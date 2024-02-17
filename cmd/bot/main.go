package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/fumkaa/pocket-bot/pkg/config"
	dbMysql "github.com/fumkaa/pocket-bot/pkg/db/MySQL"
	pClient "github.com/fumkaa/pocket-bot/pkg/pocket/client/Authentication"
	"github.com/fumkaa/pocket-bot/pkg/server"
	"github.com/fumkaa/pocket-bot/pkg/telegram"
	"github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cnf, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}
	log.Print(cnf)
	bot, err := tgbotapi.NewBotAPI(cnf.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false

	pocketClient := pClient.New(cnf.PocketConsumerKey)

	db, err := connDatabase(cnf.DBname, cnf.User, cnf.Password, cnf.Port)
	if err != nil {
		log.Fatal(err)
	}

	telegramBot := telegram.NewBot(bot, *pocketClient, cnf.AuthServerURL, db, cnf.Message)
	authServer := server.NewAuthorizationServer(pocketClient, dbMysql.NewDatabase(db), cnf.TelegramBotURL)
	go func() {
		telegramBot.Start()
	}()

	if err := authServer.Start(); err != nil {
		log.Fatal(err)
	}
}

func connDatabase(DBName string, user string, passwd string, port string) (*sql.DB, error) {
	cfg := mysql.Config{
		User:              user,
		Passwd:            passwd,
		Net:               "tcp",
		Addr:              port,
		DBName:            DBName,
		InterpolateParams: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return &sql.DB{}, fmt.Errorf("[connDatabase]can't create pool: %w", err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return &sql.DB{}, fmt.Errorf("[connDatabase]can't connected to database: %w", pingErr)
	}
	log.Print("Connected to database!")
	return db, nil
}
