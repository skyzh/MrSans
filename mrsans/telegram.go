package main

import (
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"time"
)

// InitializeTelegramBot gets the authorized account of telegram bot key
func InitializeTelegramBot(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	c := make(chan struct{})
	go func() {
		bot, err := tgbotapi.NewBotAPI(Config.telegram_bot_token)
		if err != nil {
			log.Fatalf("failed to initialize bot API: %v", err)
		}
		log.Infof("authorized as account %s", bot.Self.UserName)
		c <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		log.Fatalf("failed to initialize telegram bot: %v", ctx.Err())
	case <-c:
	}
}

// SensePushLog pushes message to telegram channel `telegram.chat_id`
func SensePushMessage(caption string, photo string) error {
	bot, err := tgbotapi.NewBotAPI(Config.telegram_bot_token)
	if err != nil {
		log.Warnf("failed to initialize bot API: %v", err)
		reportFailure.Add(1)
		return err
	}

	msg := tgbotapi.NewPhotoUpload(Config.telegram_chat_id, photo)
	msg.Caption = caption
	msg.ParseMode = "markdown"

	_, err = bot.Send(msg)

	if err != nil {
		log.Warnf("failed to send message: %v", err)
		reportFailure.Add(1)
		return err
	}

	reportSuccess.Add(1)

	return nil
}

// SensePushLog pushes message to telegram channel `telegram.log_chat_id`
func SensePushLog(message string) error {
	bot, err := tgbotapi.NewBotAPI(Config.telegram_bot_token)
	if err != nil {
		log.Warnf("failed to initialize bot API %v", err)
		reportFailure.Add(1)
		return err
	}

	msg := tgbotapi.NewMessage(Config.telegram_log_chat_id, message)
	msg.ParseMode = "markdown"
	msg.DisableNotification = true
	_, err = bot.Send(msg)

	if err != nil {
		log.Warnf("failed to send message %v", err)
		reportFailure.Add(1)
		return err
	}

	reportSuccess.Add(1)

	return nil
}
