package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	prometheus.MustRegister(dcpScrapeError)

	// dcpScrapeError is true since I presume that it will working at the start
	dcpScrapeError.With(prometheus.Labels{"error": "0"}).Set(1)

	start := time.Now()

	var flagOnline bool = true

	url := os.Getenv("URL")

	var authHeader HeaderObject
	authHeader.Security.MustUnderstand = "1"
	authHeader.Security.Wsse = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd"
	authHeader.Security.Wsu = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd"

	err := GetAuthValues(&authHeader.Security.UsernameToken.Username, &authHeader.Security.UsernameToken.Password)
	if err != nil {
		log.Printf("error: %v\n", err)
		flagOnline = false
		dcpScrapeError.With(prometheus.Labels{"error": "auth error"}).Set(0)
	}

	if flagOnline {
		SetApiStatusMetric(&authHeader, url, &flagOnline)
		SetSimVolumeMetric(&authHeader, url, &flagOnline)
	}

	elapsed := time.Since(start)
	dcpScrapeDurationSeconds.Set(elapsed.Seconds())

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9742", nil)
}
