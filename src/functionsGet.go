package main

import (
	"encoding/xml"
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

func GetEchoHost(authHeader *HeaderObject, url string) (bool, error) {

	isEcho := false

	urlEcho := url + "/dcpapi/ApiStatus?WSDL"

	// Build the body
	var envelope EnvelopeObjectEcho
	envelope.Header = *authHeader
	envelope.Body.Echo.Message = "hello"

	envelope.Soapenv = "http://schemas.xmlsoap.org/soap/envelope/"
	envelope.Apis = "http://api.dcp.ericsson.net/ApiStatus"

	// marshal the envelope
	body, err := xml.Marshal(envelope)
	if err != nil {
		log.Printf("error: %v", err)
		dcpScrapeError.With(prometheus.Labels{"error": "echo marchalling"}).Set(0)
		return isEcho, err
	}

	// Convert the byte slice to a string
	bodySend := string(body)

	// Send the request
	bodyResp, err := HTTPCaller(bodySend, urlEcho, "POST", "text/xml; charset=utf-8")
	if err != nil {
		log.Printf("error: %v", err)
		dcpScrapeError.With(prometheus.Labels{"error": "echo request"}).Set(0)
		return isEcho, err
	}

	// Unmarshal the response body into EnvelopeRespEcho
	var envelopeResp EnvelopeRespEcho
	err = xml.Unmarshal([]byte(bodyResp), &envelopeResp)
	if err != nil {
		log.Printf("error: %v", err)
		dcpScrapeError.With(prometheus.Labels{"error": "echo unmarchalling"}).Set(0)
		return isEcho, err
	}

	// Check the response
	if envelopeResp.Body.EchoResponse.Message == "hello" {
		isEcho = true
	}

	return isEcho, nil
}

func GetSimVolume(authHeader *HeaderObject, url string) (int, error) {

	return 0, nil
}