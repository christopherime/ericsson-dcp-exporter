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

	log.Println("Starting DCP exporter on port 9742...")

	flagOnline := true
	url := os.Getenv("URL")

	var subWatch string

	prometheus.MustRegister(dcpScrapeError, dcpScrapeDurationSeconds)

	// Load value in config.yaml
	config, _ := LoadConfig()
	if !config.configPresent {
		log.Printf("No pre config present, using default values")
		subWatch = os.Getenv("SUBPACKAGE")
	}

	// maxSims the maximum number of sims to be scraped
	maxSims := os.Getenv("MAXSIMS")
	if maxSims == "" {
		maxSims = "5000"
	}

	// dcpScrapeError is true since I presume that it will working at the start
	dcpScrapeError.With(prometheus.Labels{"error": "0"}).Set(1)

	log.Println("Initialisation finished, starting scraping...")

	start := time.Now()

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

	// Set the API status metric to 1 if the API is online
	log.Println("Getting API status...")
	SetApiStatusMetric(&authHeader, url, &flagOnline)

	// If the API is online, scrape the metrics
	// otherwise just don't bother
	if flagOnline {

		// Go through all subpackages and scrape the sims
		// test if config.configValues is empty
		if len(config.configValues.SubToWatch) == 0 {
			log.Printf("No config present, using default values")
			SetSimVolumeMetric(&authHeader, url, subWatch, maxSims)
		} else {
			log.Printf("Config present, using config values")
			for _, subpackage := range config.configValues.SubToWatch {
				SetSimVolumeMetric(&authHeader, url, subpackage, maxSims)
			}
		}

	}

	elapsed := time.Since(start)
	dcpScrapeDurationSeconds.Set(elapsed.Seconds())

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9742", nil)
}
