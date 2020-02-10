package main

import (
	"context"
	"fmt"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"time"
)

func GetLatestValue(x *[]model.SamplePair) float64 {
	return float64((*x)[len(*x) - 1].Value)
}

func SenseGenerateMessage(temp *[]model.SamplePair, hum *[]model.SamplePair, pa *[]model.SamplePair, pm25 *[]model.SamplePair, pm10 *[]model.SamplePair) string {
	return fmt.Sprintf("Mr. Sans Reporting\n" +
		"%s\n\n" +
		"*Temperature* %.2f°C\n" +
		"*Humidity* %.2f%%\n" +
		"*Pressure* %.0f Pa\n" +
		"*PM 2.5* %.2f µg/m3\n" +
		"*PM 10* %.2f µg/m3\n",
		time.Now().Format("Mon Jan 2 15:04:05 2006"),
		GetLatestValue(temp), GetLatestValue(hum),
		GetLatestValue(pa),
		GetLatestValue(pm25), GetLatestValue(pm10))
}

func ReportOnce() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	log.Info("start reporting sequence...")

	log.Info("> getting data from prometheus...")
	temp := GetData(QueryTemperature(), v1.Range{
		Start: time.Now().Add(-time.Hour * 24),
		End:   time.Now(),
		Step:  time.Minute,
	}, ctx)
	hum := GetData(QueryHumidity(), v1.Range{
		Start: time.Now().Add(-time.Hour * 24),
		End:   time.Now(),
		Step:  time.Minute,
	}, ctx)
	pa := GetData(QueryPressure(), v1.Range{
		Start: time.Now().Add(-time.Hour * 24),
		End:   time.Now(),
		Step:  time.Minute,
	}, ctx)
	pm25 := GetData(QueryPM25(), v1.Range{
		Start: time.Now().Add(-time.Hour * 24),
		End:   time.Now(),
		Step:  time.Minute,
	}, ctx)
	pm10 := GetData(QueryPM10(), v1.Range{
		Start: time.Now().Add(-time.Hour * 24),
		End:   time.Now(),
		Step:  time.Minute,
	}, ctx)
	log.Info("> plotting...")
	Plot(&temp, &hum, &pa, &pm25, &pm10, "out/report.png")
	log.Info("> sending message...")
	message := SenseGenerateMessage(&temp, &hum, &pa, &pm25, &pm10)
	SensePushMessage(message, "out/report.png")
}

func RunCronTask() {
	log.Info("Scheduler running")
	c := cron.New()
	c.AddFunc("@every 1m", func() { ReportOnce() })
	log.Info(c.Entries())
	c.Run()
}

func main() {
	log.Info("loading config...")
	LoadConfig()
	InitializeTelegramBot()
	SensePushLog("Mr Sans intialized")
	RunCronTask()
}
