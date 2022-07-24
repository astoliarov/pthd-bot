package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func InitBot(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return bot, nil
}
