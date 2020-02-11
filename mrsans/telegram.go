package main

import (
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"time"
)

func InitializeTelegramBot(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	c := make(chan struct{})
	go func() {
		bot, err := tgbotapi.NewBotAPI(Config.telegram_bot_token)
		if err != nil {
			log.Fatal("failed to initialize bot API", err)
		}
		log.Infof("authorized as account %s", bot.Self.UserName)
		c <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		log.Fatal("failed to initialize telegram bot: ", ctx.Err())
	case <-c:
	}
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
