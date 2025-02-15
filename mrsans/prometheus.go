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
	reportFailure = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mrsans_telegram_failure_count",
		Help: "Count of Telegram push failure",
	})
	reportSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mrsans_telegram_success_count",
		Help: "Count of Telegram push success",
	})
	checkpoint = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "mrsans_checkpoint_seconds",
		Help: "Duration of checkpoint task",
	})
	checkpointFailure = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mrsans_checkpoint_failure_count",
		Help: "Duration of checkpoint failure",
	})
)

func RunPrometheus() {
	log.Infof("setting up prometheus server at %s", Config.prometheus_addr)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(Config.prometheus_addr, nil))
}
