package main

import (
	"encoding/xml"
	"log"

	"github.com/dgraph-io/badger/v3"
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

// Get the number of SIMs in a subscription package
// also store the imsi value of each SIM in a badger db for later use
func GetSimVolume(authHeader *HeaderObject, url string, subPackage string, maxSims string, badger *badger.DB) (int, error) {

	urlSub := url + "/dcpapi/SubscriptionManagement?WSDL"

	// Build the body
	var envelope EnvelopeObjectSubMgmt

	envelope.Soapenv = "http://schemas.xmlsoap.org/soap/envelope/"
	envelope.Sub = "http://api.dcp.ericsson.net/SubscriptionManagement"

	envelope.Header = *authHeader
	envelope.Body.QuerySubscriptionsRequest.SubscriptionPackage = subPackage
	envelope.Body.QuerySubscriptionsRequest.MaxResults = maxSims

	// marshal the envelope
	body, err := xml.Marshal(envelope)
	if err != nil {
		log.Printf("error: %v", err)
		dcpScrapeError.With(prometheus.Labels{"error": "echo marchalling"}).Set(0)
		return 0, err
	}

	// Convert the byte slice to a string
	bodySend := string(body)

	// Send the request
	bodyResp, err := HTTPCaller(bodySend, urlSub, "POST", "text/xml; charset=utf-8")
	if err != nil {
		log.Printf("error: %v", err)
		dcpScrapeError.With(prometheus.Labels{"error": "echo request"}).Set(0)
		return 0, err
	}

	// Unmarshal the response body into EnvelopeRespEcho
	var envelopeResp EnvelopeRespSubMgmt

	err = xml.Unmarshal([]byte(bodyResp), &envelopeResp)
	if err != nil {
		log.Printf("error: %v", err)
		dcpScrapeError.With(prometheus.Labels{"error": "echo unmarchalling"}).Set(0)
		return 0, err
	}

	// Store the IMSI in a badger db

	go SaveBadgerDB(envelopeResp, badger)

	// Count the number of subscription in the response
	simVolume := len(envelopeResp.Body.QuerySubscriptionsResponse.Subscriptions.Subscription)

	return simVolume, nil
}

// Read the database and return the Imsi
func GetImsiByKey(key string, badger *badger.DB) ([]string, error) {
	var Imsi string
	var ImsiList []string

	// Read the database
	Txn := badger.NewTransaction(false)
	defer Txn.Discard()

	// Get the Imsi from the database
	item, err := Txn.Get([]byte(key))
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	// Convert the Imsi to a string
	err = item.Value(func(val []byte) error {
		Imsi = string(val)
		return nil
	})
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	// Append the Imsi to the ImsiList
	ImsiList = append(ImsiList, Imsi)

	return ImsiList, nil
}
