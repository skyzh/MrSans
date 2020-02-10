package main

import (
	"context"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("loading config...")
	LoadConfig()
	GetData(context.Background())
}
