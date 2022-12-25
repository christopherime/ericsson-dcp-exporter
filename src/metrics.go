package main

import "github.com/prometheus/client_golang/prometheus"

var apiStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "api_status",
	Help: "API status as a boolean value (1 = true, 0 = false)",
}, []string{"url"})

var dcpScrapeDurationSeconds = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "dcp_scrape_duration_seconds",
	Help: "Time taken to scrape the DCP",
})

var dcpScrapeError = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "dcp_scrape_error",
	Help: "Error during a scraping the DCP",
}, []string{"error"})

var simVolume = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "sim_volume",
	Help: "Number of sim cards in the system",
})
