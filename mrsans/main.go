package main

import (
	"context"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	log.Info("loading config...")
	LoadConfig()
	temp := GetData(QueryTemperature(), v1.Range{
		Start: time.Now().Add(-time.Hour * 12),
		End:   time.Now(),
		Step:  time.Minute,
	}, context.Background())
	hum := GetData(QueryHumidity(), v1.Range{
		Start: time.Now().Add(-time.Hour * 12),
		End:   time.Now(),
		Step:  time.Minute,
	}, context.Background())
	pa := GetData(QueryPressure(), v1.Range{
		Start: time.Now().Add(-time.Hour * 12),
		End:   time.Now(),
		Step:  time.Minute,
	}, context.Background())
	pm25 := GetData(QueryPM25(), v1.Range{
		Start: time.Now().Add(-time.Hour * 12),
		End:   time.Now(),
		Step:  time.Minute,
	}, context.Background())
	pm10 := GetData(QueryPM10(), v1.Range{
		Start: time.Now().Add(-time.Hour * 12),
		End:   time.Now(),
		Step:  time.Minute,
	}, context.Background())
	Plot(&temp, &hum, &pa, &pm25, &pm10)
}
