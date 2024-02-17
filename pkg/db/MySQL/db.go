package dbMysql

import (
	"database/sql"
	"fmt"
	"log"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{db: db}
}

func (d *Database) SaveRequestToken(chatId int, token string) error {
	row := d.db.QueryRow("SELECT access_token FROM user WHERE chat_id = ?;", chatId)
	log.Printf("[SaveRequestToken]chat id: %d", chatId)
	var acctoken string
	if err := row.Scan(&acctoken); err != nil {
		return fmt.Errorf("[SaveRequestToken]can't get access token: %w", err)
	}
	if acctoken == "" {
		_, err := d.db.Exec("UPDATE user SET request_token = ? WHERE chat_id = ?;", token, chatId)
		if err != nil {
			return fmt.Errorf("[SaveRequestToken]can't paste request token: %w", err)

		}
	}

	return nil
}

func (d *Database) SaveAccessToken(chatId int, token string) error {
	_, err := d.db.Exec("UPDATE user SET access_token = ? WHERE chat_id = ?;", token, chatId)
	if err != nil {
		return fmt.Errorf("[SaveAccessToken]can't paste access token: %w", err)
	}
	return nil
}

func (d *Database) RequestToken(chatId int) (string, error) {
	var token string
	row := d.db.QueryRow("SELECT request_token FROM user WHERE chat_id = ?;", chatId)

	if err := row.Scan(&token); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("[RequestToken]ChatId %d: no such user", chatId)
		}
		return "", fmt.Errorf("[RequestToken]can't get request token: %w", err)
	}
	return token, nil
}

func (d *Database) AccessToken(chatId int) (string, error) {
	var token string
	row := d.db.QueryRow("SELECT access_token FROM user WHERE chat_id = ?;", chatId)

	if err := row.Scan(&token); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("[AccessToken]ChatId %d: no such user", chatId)
		}
		return "", fmt.Errorf("[AccessToken]can't get access token: %w", err)
	}
	return token, nil
}

func (d *Database) CreateTable(chatId int) error {
	_, err := d.db.Exec("USE pocket_bot;")
	if err != nil {
		return fmt.Errorf("[CreateTable]can't use database: %w", err)
	}
	_, err = d.db.Exec("CREATE TABLE IF NOT EXISTS user (chat_id INT NOT NULL, request_token VARCHAR(225) DEFAULT \"\" NOT NULL, access_token VARCHAR(225) DEFAULT \"\" NOT NULL, PRIMARY KEY (`chat_id`));")
	if err != nil {
		return fmt.Errorf("[CreateTable]can't create table: %w", err)
	}

	row := d.db.QueryRow("SELECT chat_id FROM user ;")

	err = row.Scan()
	if err == sql.ErrNoRows {
		_, err = d.db.Exec("INSERT INTO user (chat_id) VALUES (?);", chatId)
		if err != nil {
			return fmt.Errorf("[CreateTable]can't paste chat id: %w", err)
		}
	}

	return nil
}
