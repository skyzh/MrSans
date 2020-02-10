package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/fogleman/gg"
	"firebase.google.com/go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"http"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	log.info("Running!")
}