package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

func InitializeTelegramBot() {
	bot, err := tgbotapi.NewBotAPI(Config.telegram_bot_token)
	if err != nil {
		log.Fatal("failed to initialize bot API", err)
	}
	log.Infof("authorized as account %s", bot.Self.UserName)
}

func SensePushMessage(caption string, photo string) {
	bot, err := tgbotapi.NewBotAPI(Config.telegram_bot_token)
	if err != nil {
		log.Fatal("failed to initialize bot API", err)
	}

	msg := tgbotapi.NewPhotoUpload(Config.telegram_chat_id, photo)
	msg.Caption = caption
	msg.ParseMode = "markdown"

	_, err = bot.Send(msg)

	if err != nil {
		log.Warn("failed to send message", err)
	}
}

func SensePushLog(message string) {
	bot, err := tgbotapi.NewBotAPI(Config.telegram_bot_token)
	if err != nil {
		log.Fatal("failed to initialize bot API", err)
	}

	msg := tgbotapi.NewMessage(Config.telegram_log_chat_id, message)
	msg.ParseMode = "markdown"
	_, err = bot.Send(msg)

	if err != nil {
		log.Warn("failed to send message", err)
	}
}
