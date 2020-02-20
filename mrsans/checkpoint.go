package main

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"fmt"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	log "github.com/sirupsen/logrus"
	"math"
	"sync"
	"time"
)

type Checkpoint struct {
	Temperature float64 `json:"temp,omitempty"`
	Humidity    float64 `json:"hum,omitempty"`
	Pressure    float64 `json:"pa,omitempty"`
	PM25        float64 `json:"pm25,omitempty"`
	PM10        float64 `json:"pm10,omitempty"`
	Time        int64   `json:"time,omitempty"`
}

type CheckpointService struct {
	client *db.Client
	mux    sync.Mutex
}

func CheckpointOf(path string) string {
	return fmt.Sprintf("checkpoint/%s/%s", Config.checkpoint_base, path)
}

func DoMinuteCheckpoint(ctx context.Context, client *db.Client) error {
	log := log.WithField("job", "minute checkpoint job")
	ref := client.NewRef(CheckpointOf("minute"))
	results, err := ref.OrderByChild("time").LimitToLast(1).GetOrdered(ctx)
	if err != nil {
		log.Warnf("failed to get latest checkpoint: %v", err)
		return err
	}
	timeTruncate := time.Minute
	fromTime := int64(0)
	if len(results) == 0 {
		t, _ := time.Parse(time.RFC822, "12 Feb 20 20:00 CST")
		fromTime = t.Unix()
		log.Warnf("no checkpoint record in database, using %v as initial", t)
	} else {
		var checkpoint Checkpoint
		if err := results[0].Unmarshal(&checkpoint); err != nil {
			log.Warnf("failed to unmarshal checkpoint data: %v", err)
			return err
		}
		fromTime = checkpoint.Time
	}

	checkpointNow := time.Now()

	// count := int(GetData(Count(QueryTemperature()), checkpointNow, ctx).Value)
	// log.Warn(count)

	for {
		fromT := time.Unix(fromTime, 0).Truncate(timeTruncate)
		endT := checkpointNow.Truncate(timeTruncate)

		if fromT.Equal(endT) || fromT.After(endT) {
			break
		}

		if t := fromT.Add(time.Hour); t.Before(endT) {
			endT = t
		}

		r := v1.Range{
			Start: fromT,
			End:   endT,
			Step:  timeTruncate,
		}

		temp := GetRange(QueryTemperature(), r, ctx)
		hum := GetRange(QueryHumidity(), r, ctx)
		pa := GetRange(QueryPressure(), r, ctx)
		pm25 := GetRange(QueryPM25(), r, ctx)
		pm10 := GetRange(QueryPM10(), r, ctx)

		length := len(temp)

		if length != 0 {
			log.Infof("checkpoint %s ~ %s (not included)", fromT.Format(time.RFC3339), endT.Format(time.RFC3339))
		}

		for idx := 0; idx < length; idx++ {
			atTime := temp[idx].Timestamp.Time().Truncate(timeTruncate)
			temp := float64(temp[idx].Value)
			if math.IsNaN(temp) {
				continue
			}
			hum := float64(hum[idx].Value)
			if math.IsNaN(temp) {
				continue
			}
			pm25 := float64(pm25[idx].Value)
			if math.IsNaN(temp) {
				continue
			}
			pm10 := float64(pm10[idx].Value)
			if math.IsNaN(temp) {
				continue
			}
			pa := float64(pa[idx].Value)
			if math.IsNaN(temp) {
				continue
			}

			checkpoint := Checkpoint{
				Temperature: temp,
				Humidity:    hum,
				Pressure:    pa,
				PM25:        pm25,
				PM10:        pm10,
				Time:        atTime.Unix(),
			}

			if _, err := ref.Push(ctx, &checkpoint); err != nil {
				log.Warnf("failed when checkpoint %s: %v", atTime.Format(time.RFC1123), err)
				return err
			}
		}

		fromTime = endT.Unix()
	}

	log.Info("checkpoint complete")

	return nil
}

func DoCheckpoint(ctx context.Context, client *db.Client) error {
	log := log.WithField("job", "checkpoint job")
	t := time.Now()

	if err := DoMinuteCheckpoint(ctx, client); err != nil {
		log.Warnf("error while checkpoint minute")
		return err
	}

	checkpoint.Observe(time.Now().Sub(t).Seconds())
	return nil
}

func GetCheckpointService() *CheckpointService {
	ctx := context.Background()
	log := log.WithField("job", "checkpoint")
	conf := &firebase.Config{
		DatabaseURL: Config.firebase_url,
	}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}
	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalf("error initializing database: %v", err)
	}
	return &CheckpointService{client: client}
}

func (c *CheckpointService) RunCheckpoint() {
	c.mux.Lock()
	defer c.mux.Unlock()

	if err := DoCheckpoint(context.Background(), c.client); err != nil {
		log.Warnf("error while checkpoint")
		checkpointFailure.Add(1)
	}
}
