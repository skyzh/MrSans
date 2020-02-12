package main

import (
	"context"
	"fmt"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func RunCronTask() {
	log.Info("Scheduler running")
	c := cron.New()
	_, err := c.AddFunc("@hourly", func() { ReportHourlyOnce() })
	if err != nil {
		log.Fatal("failed to add hourly task")
	}
	_, err = c.AddFunc("@daily", func() { ReportDailyOnce() })
	if err != nil {
		log.Fatal("failed to add daily task")
	}
	_, err = c.AddFunc("0 12 * * *", func() { ReportDailyOnce() })
	if err != nil {
		log.Fatal("failed to add daily task 2")
	}
	_, err = c.AddFunc("@every 5m", func() {
		for idx, entry := range c.Entries() {
			log.Infof("Cron tasks #%d is scheduled at %s", idx, entry.Next.Format("Mon Jan 2 15:04"))
		}
	})
	if err != nil {
		log.Fatal("failed to add cron debug task")
	}
	checkpointService := GetCheckpointService()
	go checkpointService.RunCheckpoint()
	_, err = c.AddFunc("@every 1m", func() {
		checkpointService.RunCheckpoint()
	})
	if err != nil {
		log.Fatal("failed to add checkpoint service")
	}
	c.Run()
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.Info("loading config...")
	LoadConfig()
	log.Info("initialize telegram bot...")
	InitializeTelegramBot(context.Background())
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("failed to get hostname: %v", err)
	}
	go RunPrometheus()
	go RunGrafanaWebhook()
	go SensePushLog(fmt.Sprintf("Mr Sans intialized\n"+
		"%s\n"+
		"*host* `%s`\n"+
		"*bluesense_host* `%s`\n"+
		"*bluesense_job* `%s`\n"+
		"*site_name* `%s`",
		time.Now().Format("Mon Jan 2 15:04:05 MST 2006"),
		hostname, Config.bluesense_url, Config.bluesense_job, Config.site_name))
	if Config.instant_push {
		go ReportDailyOnce()
		go ReportHourlyOnce()
	}
	RunCronTask()
}
