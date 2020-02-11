package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func RunPrometheus() {
	log.Infof("setting up prometheus server at %s", Config.prometheus_addr)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(Config.prometheus_addr, nil)
}
