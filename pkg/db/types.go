package types

type TokenDb interface {
	SaveRequestToken(chatId int, token string) error
	SaveAccessToken(chatId int, token string) error
	RequestToken(chatId int) (string, error)
	AccessToken(chatId int) (string, error)
}
