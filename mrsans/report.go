package main

import (
	"context"
	"fmt"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	log "github.com/sirupsen/logrus"
	"time"
)

func GetLatestValue(x *[]model.SamplePair) float64 {
	return float64((*x)[len(*x)-1].Value)
}

func SenseGenerateMessage(msg string, temp *[]model.SamplePair, hum *[]model.SamplePair, pa *[]model.SamplePair, pm25 *[]model.SamplePair, pm10 *[]model.SamplePair) string {
	return fmt.Sprintf("🤖 *Mr. Sans Reporting* %s\n"+
		"%s\n"+
		"%s\n\n"+
		"*Temperature* %.2f°C\n"+
		"*Humidity* %.2f%%\n"+
		"*Pressure* %.0f Pa\n"+
		"*PM2.5* %.2f µg/m3\n"+
		"*PM10* %.2f µg/m3\n",
		msg,
		Config.site_name,
		time.Now().Format("Mon Jan 2 15:04 MST 2006"),
		GetLatestValue(temp), GetLatestValue(hum),
		GetLatestValue(pa),
		GetLatestValue(pm25), GetLatestValue(pm10))
}

func ReportHourlyOnce() {
	t := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	log := log.WithField("job", "hourly report")
	log.Info("start reporting hourly sequence...")

	log.Info("> getting data from prometheus...")
	r := v1.Range{
		Start: time.Now().Add(-time.Hour * 24),
		End:   time.Now(),
		Step:  time.Minute,
	}
	temp := GetRange(QueryTemperature(), r, ctx)
	hum := GetRange(QueryHumidity(), r, ctx)
	pa := GetRange(QueryPressure(), r, ctx)
	pm25 := GetRange(QueryPM25(), r, ctx)
	pm10 := GetRange(QueryPM10(), r, ctx)
	log.Info("> plotting...")
	msg := fmt.Sprintf("Hourly @ %s", Config.site_name)
	Plot(msg, time.Hour, 0, &temp, &hum, &pa, &pm25, &pm10, "out/report_hourly.png")
	log.Info("> sending message...")
	message := SenseGenerateMessage("#Hourly", &temp, &hum, &pa, &pm25, &pm10)
	SensePushMessage(message, "out/report_hourly.png")
	log.Info("> done")
	hourlyReport.Observe(time.Now().Sub(t).Seconds())
}

func ReportDailyOnce() {
	t := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	log := log.WithField("job", "daily report")
	log.Info("start reporting daily sequence...")

	log.Info("> getting data from prometheus...")
	r := v1.Range{
		Start: time.Now().Add(-time.Hour * 7 * 24),
		End:   time.Now(),
		Step:  time.Minute * 10,
	}
	temp := GetRange(QueryTemperature(), r, ctx)
	hum := GetRange(QueryHumidity(), r, ctx)
	pa := GetRange(QueryPressure(), r, ctx)
	pm25 := GetRange(QueryPM25(), r, ctx)
	pm10 := GetRange(QueryPM10(), r, ctx)
	log.Info("> plotting...")
	msg := fmt.Sprintf("Daily @ %s", Config.site_name)
	_, offset := time.Now().Zone()
	Plot(msg, time.Hour*24, time.Duration(-offset)*time.Second, &temp, &hum, &pa, &pm25, &pm10, "out/report_daily.png")
	log.Info("> sending message...")
	message := SenseGenerateMessage("#Daily", &temp, &hum, &pa, &pm25, &pm10)
	SensePushMessage(message, "out/report_daily.png")
	log.Info("> done")
	dailyReport.Observe(time.Now().Sub(t).Seconds())
}
