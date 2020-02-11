package main

import (
	"github.com/pelletier/go-toml"
	log "github.com/sirupsen/logrus"
)

var Config struct {
	bluesense_url string
	bluesense_job string
	telegram_bot_token string
	telegram_chat_id int64
	telegram_log_chat_id int64
	plot_fontface string
	instant_push bool
	site_name string
	prometheus_addr string
}

func LoadConfig() {
	config, err := toml.LoadFile("config.toml")
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}
	Config.bluesense_url = config.Get("bluesense.prometheus").(string)
	Config.bluesense_job = config.Get("bluesense.job").(string)
	Config.telegram_bot_token = config.Get("telegram.bot_token").(string)
	Config.telegram_chat_id = config.Get("telegram.chat_id").(int64)
	Config.telegram_log_chat_id = config.Get("telegram.log_chat_id").(int64)
	Config.plot_fontface = config.Get("plot.fontface").(string)
	Config.instant_push = config.Get("bluesense.instant_push").(bool)
	Config.site_name = config.Get("bluesense.site_name").(string)
	Config.prometheus_addr = config.Get("exporter.addr").(string)
}
