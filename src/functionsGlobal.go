package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func LoadConfig() (TypeConfig, error) {

	var config TypeConfig

	// check if config.yaml exists
	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		log.Printf("error: %v", err)
		config.configPresent = false
		return config, err
	}

	return config, nil

}

func HTTPCaller(bodySend string, urlSend string, methodSend string, contentTypeSend string) (string, error) {

	// Build the request
	req, err := http.NewRequest(methodSend, urlSend, strings.NewReader(bodySend))
	if err != nil {
		log.Printf("error: %v", err)
		return "", err
	}

	// Set the content type
	req.Header.Set("Content-Type", contentTypeSend)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error: %v", err)
		return "", err
	}

	// Read the response
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error: %v", err)
		return "", err
	}

	// Convert the byte slice to a string
	str := string(body)

	return str, nil

}

func Sleeping5s() {

	// Sleep for 5 seconds
	time.Sleep(5 * time.Second)
}

func Sleeping1m() {

	// Sleep for 1 minute
	time.Sleep(1 * time.Minute)
}
