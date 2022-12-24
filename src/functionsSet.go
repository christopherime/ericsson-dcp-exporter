package main

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

func SetApiStatusMetric(authInfo *HeaderObject, url string, flagOnline *bool) {

	prometheus.MustRegister(apiStatus)

	go func() {
		for {

			valPint, err := GetEchoHost(authInfo, url)
			if err != nil {
				log.Printf("error: %v\n", err)

			}

			flagOnline = &valPint

			if !valPint {
				apiStatus.With(prometheus.Labels{"url": url}).Set(0)

			} else {
				apiStatus.With(prometheus.Labels{"url": url}).Set(1)
			}
			Sleeping5()
		}
	}()
}

func SetSimVolumeMetric(authHeader *HeaderObject, url string, flagOnline *bool) {
	if !*flagOnline {
		return
	}

	prometheus.MustRegister(simVolume)

	go func() {
		for {
			Sleeping5()
		}
	}()
}
