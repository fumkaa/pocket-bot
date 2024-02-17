package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken     string
	PocketConsumerKey string
	AuthServerURL     string
	User              string
	Password          string
	Port              string
	TelegramBotURL    string `mapstructure:"bot_url"`
	DBname            string `mapstructure:"pocket_bot"`

	Message Message
}

type Message struct {
	Error
	Responses
}

type Error struct {
	Default      string `mapstructure:"default"`
	InvalidUrl   string `mapstructure:"invalid_url"`
	Unauthorized string `mapstructure:"unauthorized"`
	UnableToSave string `mapstructure:"unable_to_save"`
}

type Responses struct {
	Start             string `mapstructure:"start"`
	Help              string `mapstructure:"help"`
	AlreadyAuthorized string `mapstructure:"already_authorized"`
	SavesSuccessfully string `mapstructure:"saves_successfully"`
	UnknowCommand     string `mapstructure:"unknow_command"`
}

func Init() (*Config, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("message.response", &config.Message.Responses); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("message.error", &config.Message.Error); err != nil {
		return nil, err
	}

	if err := parseEnv(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func parseEnv(config *Config) error {
	os.Setenv("TOKEN_BOT", "6834826217:AAE5AUrcscS28o26heRp2UWIluTSnEFjuPA")
	os.Setenv("CONSUMER_KEY", "110507-b116c24f0e328a441521ee0")
	os.Setenv("AUTH_SERVER_URL", "http://localhost/")
	os.Setenv("USER_DB", "root")
	os.Setenv("PASSWORD_DB", "28aprelz")
	os.Setenv("PORT_DB", "127.0.0.1:3306")
	nameEnv := []string{"TOKEN_BOT", "CONSUMER_KEY", "AUTH_SERVER_URL", "USER_DB", "PASSWORD_DB", "PORT_DB"}
	for _, n := range nameEnv {
		if err := viper.BindEnv(n); err != nil {
			return err
		}
	}

	config.TelegramToken = viper.GetString("TOKEN_BOT")
	config.PocketConsumerKey = viper.GetString("CONSUMER_KEY")
	config.AuthServerURL = viper.GetString("AUTH_SERVER_URL")
	config.User = viper.GetString("USER_DB")
	config.Password = viper.GetString("PASSWORD_DB")
	config.Port = viper.GetString("PORT_DB")
	return nil
}
