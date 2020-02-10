package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	log "github.com/sirupsen/logrus"
)

func GetSenseClient() v1.API {
	cfg := api.Config{
		Address:      Config.bluesense_url,
		RoundTripper: nil,
	}
	client, err := api.NewClient(cfg)
	if err != nil {
		log.Fatal("failed to create prometheus api client: ", err)
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


func GetData(query string, r v1.Range, ctx context.Context) [] model.SamplePair {
	v1api := GetSenseClient()

	result, warnings, err := v1api.QueryRange(ctx, query, r)
	if err != nil {
		log.Fatal("failed to query data: ", err)
	}
	if warnings != nil {
		log.Warn("query warning: ", warnings)
	}
	mat, ok := result.(model.Matrix)
	if !ok {
		log.Fatal("failed to cast data")
	}
	if mat.Len() != 1 {
		log.Fatal("more than 1 query result")
	}
	values := mat[0].Values
	return values
}
