package main

import (
	"log"

	"github.com/dgraph-io/badger/v3"
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
			Sleeping5s()
		}
	}()
}

func SetSimVolumeMetric(authHeader *HeaderObject, url string, subPackage string, maxSims string, badger *badger.DB) {
	prometheus.MustRegister(simVolume)

	go func() {
		for {
			volSim, err := GetSimVolume(authHeader, url, subPackage, maxSims, badger)
			if err != nil {
				log.Printf("error: %v\n", err)
			}

			simVolume.Set(float64(volSim))

			Sleeping1m()
		}
	}()
}
