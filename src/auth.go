package main

import (
	"encoding/base64"
	"log"
	"os"
	"strings"
)

func GetAuthValues(username *string, password *string) error {
	// load the USERAUTH env variable
	userAuth := os.Getenv("USERAUTH")

	// Decode the base64-encoded string
	decoded, err := base64.StdEncoding.DecodeString(userAuth)
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}

	// Convert the byte slice to a string
	str := string(decoded)

	// Split the string into the username and password
	parts := strings.Split(str, ":")
	*username = parts[0]
	*password = parts[1]
	return nil
}
