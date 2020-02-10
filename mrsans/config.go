package main

import (
	"github.com/pelletier/go-toml"
	log "github.com/sirupsen/logrus"
)

var Config struct {
	prometheus_url string
	bluesense_url string
	bluesense_job string
	telegram_bot_token string
	telegram_chat_id int64
	telegram_log_chat_id int64
}

func LoadConfig() {
	config, err := toml.LoadFile("config.toml")
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}
	Config.prometheus_url = config.Get("mrsans.prometheus_url").(string)
	Config.bluesense_url = config.Get("mrsans.bluesense_url").(string)
	Config.bluesense_job = config.Get("mrsans.bluesense_job").(string)
	Config.telegram_bot_token = config.Get("mrsans.telegram_bot_token").(string)
	Config.telegram_chat_id = config.Get("mrsans.telegram_chat_id").(int64)
	Config.telegram_log_chat_id = config.Get("mrsans.telegram_log_chat_id").(int64)
}
