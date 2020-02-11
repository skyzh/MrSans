package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var (
	hourlyReport = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "mrsans_hourly_report_duration_seconds",
		Help: "Duration of generating one Mr. Sans hourly report",
	})

	dailyReport = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "mrsans_daily_report_duration_seconds",
		Help: "Duration of generating one Mr. Sans daily report",
	})
)

func RunPrometheus() {
	log.Infof("setting up prometheus server at %s", Config.prometheus_addr)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(Config.prometheus_addr, nil)
}
