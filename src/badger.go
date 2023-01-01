package main

import (
	"log"
	"os"

	"github.com/dgraph-io/badger/v3"
)

func CreateBadger() (bool, error) {

	// if the database exists, return false
	if _, err := os.Stat("badger"); !os.IsNotExist(err) {
		log.Printf("error: %v", err)
		return false, err
	}

	// Create the Badger database
	opts := badger.DefaultOptions("badger")
	db, err := badger.Open(opts)
	if err != nil {
		log.Printf("error: %v", err)
		return false, err
	}
	defer db.Close()

	return true, nil
}

func InitBadgerDB(dbName string) (*badger.DB, error) {
	opts := badger.DefaultOptions(dbName)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func ExistBadgerDB(dbName string) (*badger.DB, error) {
	opts := badger.DefaultOptions(dbName)
	db, err := badger.Open(opts)
	if err != nil {
		// DB does not exist, create it
		db, err = InitBadgerDB(dbName)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

func SaveBadgerDB(envelope EnvelopeRespSubMgmt, badger *badger.DB) {

	sims := envelope.Body.QuerySubscriptionsResponse.Subscriptions.Subscription

	log.Println("Saving to BadgerDB")

	// Save the sims to the database
	Txn := badger.NewTransaction(true)
	defer Txn.Discard()

	// Loop through the sims
	// and save them to the database
	// with the SubscriptionPackageName as key and the Imsi as value
	for _, sim := range sims {
		// Create a column for the SubscriptionPackageName
		// and save the Imsi as value
		err := Txn.Set([]byte(sim.SubscriptionPackageName), []byte(sim.Imsi))
		if err != nil {
			log.Printf("error: %v", err)
		}
	}

}
