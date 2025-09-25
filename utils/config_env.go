// Package utils make utils for load enviroment variables, bot telegram notify, ...
package utils

import (
	"time"

	"github.com/spf13/viper"
)

type EnvironmentVariables struct {
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ENVIRONMENT         string        `mapstructure:"ENVIRONMENT"`
	HTTPServerAddress   string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	TimeExpiredToken    time.Duration `mapstructure:"TIME_EXPIRED_TOKEN"`
	EmailAddressSender  string        `mapstructure:"EMAIL_ADDRESS_SENDER"`
	EmailPasswordSender string        `mapstructure:"EMAIL_PASSWORD_SENDER"`
	EmailUsernameSender string        `mapstructure:"EMAIL_USERNAME_SENDER"`
	RedisServerAddress  string        `mapstructure:"REDIS_ADDRESS_SERVER"`
	RedisServerPassword string        `mapstructure:"REDIS_PASSWORD_SERVER"`
	TelegramBotToken    string        `mapstructure:"TELEGRAM_BOT_TOKEN"`
	TelegramChatID      string        `mapstructure:"TELEGRAM_CHAT_ID"`
}

func LoadEnviromentVariables(path string) (config EnvironmentVariables, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env") // json, xml, ...

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
