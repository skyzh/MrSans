package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os/exec"
	"time"
)

type GrafanaAlert struct {
	State string `json:"state"`
	Tags map[string]string `json:"tags"`
}

func RunGrafanaWebhook() {
	log := log.WithField("job", "grafana")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var alert GrafanaAlert

		err := json.NewDecoder(r.Body).Decode(&alert)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Warn("failed to handle webhook", err)
			return
		}

		if alert.State == "alerting" {
			log.Warnf("grafana state alerting with tags %+v, restarting bluesense service", alert.Tags)
			for k := range alert.Tags {
				if k == "mrsans-do" {
					val := alert.Tags[k]
					if val == "restart-systemctl" {
						cmd := exec.Command("systemctl", "restart", "bluesense")
						err := cmd.Run()
						if err != nil {
							log.Warn("failed to run command", err)
						}
					}
					if val == "reboot" {
						go func() {
							log.Warn("restart after 5 seconds...", err)
							time.Sleep(time.Second * 5)
							cmd := exec.Command("reboot")
							err := cmd.Run()
							if err != nil {
								log.Warn("failed to run command reboot", err)
							}
						} ()
					}
				}
			}
		}

		fmt.Fprintf(w, "Success")
	})

	log.Infof("starting grafana webhook service at %s", Config.grafana_addr)

	log.Fatal(http.ListenAndServe(Config.grafana_addr, nil))
}
