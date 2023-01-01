package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// TODO: SetSimDetailsMetric to finish

func Init() {
	// If there is a badger folder, delete it
	if _, err := os.Stat("badger"); !os.IsNotExist(err) {
		os.RemoveAll("badger")
	}

}

func main() {

	log.Println("Starting DCP exporter on port 9742...")

	flagOnline := true
	url := os.Getenv("URL")

	var subWatch string

	// Create a new badger DB to store the SIMs values
	badgerDB, err := ExistBadgerDB("badger")
	if err != nil {
		log.Printf("error: %v", err)
	}

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
		maxSims = "2000"
	}

	// dcpScrapeError is true since I presume that it will working at the start
	dcpScrapeError.With(prometheus.Labels{"error": "0"}).Set(1)

	log.Println("Initialisation finished, starting scraping...")

	start := time.Now()

	var authHeader HeaderObject
	authHeader.Security.MustUnderstand = "1"
	authHeader.Security.Wsse = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd"
	authHeader.Security.Wsu = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd"

	err = GetAuthValues(&authHeader.Security.UsernameToken.Username, &authHeader.Security.UsernameToken.Password)
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
			SetSimVolumeMetric(&authHeader, url, subWatch, maxSims, badgerDB)
		} else {
			log.Printf("Config present, using config values")
			for _, subpackage := range config.configValues.SubToWatch {
				SetSimVolumeMetric(&authHeader, url, subpackage, maxSims, badgerDB)
			}
		}

		// Mesure details for each sim
		// Wait for 2 minutes until the SimVolumeMetric is scraped
		// This is to avoid race conditions as SetSimVolumeMetric will trigger
		// an Async function to register all Imsi value in the badgerDB

		time.Sleep(2 * time.Minute)

	}

	elapsed := time.Since(start)
	dcpScrapeDurationSeconds.Set(elapsed.Seconds())

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9742", nil)
}
