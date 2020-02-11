package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func RunPrometheus() {
	log.Info("setting up prometheus server")
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9400", nil)
}
