package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	log "github.com/sirupsen/logrus"
	"time"
)

func GetSenseClient() v1.API {
	cfg := api.Config{
		Address:      Config.bluesense_url,
		RoundTripper: nil,
	}
	client, err := api.NewClient(cfg)
	if err != nil {
		log.Fatalf("failed to create prometheus api client: %v", err)
	}
	return v1.NewAPI(client)
}

func QueryBlueSense(label string) string {
	return fmt.Sprintf("%s{job='pushgateway', exported_job='%s'}", label, Config.bluesense_job)
}

func QueryTemperature() string {
	return QueryBlueSense("temp")
}

func QueryHumidity() string {
	return QueryBlueSense("hum")
}

func QueryPressure() string {
	return QueryBlueSense("pa")
}

func QueryPM10() string {
	return QueryBlueSense("pm10")
}

func QueryPM25() string {
	return QueryBlueSense("pm25")
}

func Count(query string) string {
	return fmt.Sprintf("count_over_time(%s[1y])", query)
}

func GetRange(query string, r v1.Range, ctx context.Context) [] model.SamplePair {
	v1api := GetSenseClient()

	result, warnings, err := v1api.QueryRange(ctx, query, r)
	if err != nil {
		log.Fatalf("failed to query data: %v", err)
	}
	if warnings != nil {
		log.Warnf("query warning: %v", warnings)
	}
	mat, ok := result.(model.Matrix)
	if !ok {
		log.Fatal("failed to cast data")
	}
	if mat.Len() > 1 {
		log.Fatalf("more than 1 query result")
	}
	if mat.Len() == 0 {
		return make([] model.SamplePair, 0)
	}
	values := mat[0].Values
	return values
}

func GetData(query string, atTime time.Time, ctx context.Context) *model.Sample {
	v1api := GetSenseClient()

	result, warnings, err := v1api.Query(ctx, query, atTime)
	if err != nil {
		log.Fatalf("failed to query data: %v", err)
	}
	if warnings != nil {
		log.Warnf("query warning: %v", warnings)
	}
	vec, ok := result.(model.Vector)
	if !ok {
		log.Fatal("failed to cast data")
	}
	if vec.Len() != 1 {
		log.Fatal("more than 1 query result")
	}
	return vec[0]
}
