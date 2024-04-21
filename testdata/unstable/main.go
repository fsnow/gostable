package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func init() {
	// Set up MongoDB client
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Array of functions to be called
	functions := []func(){
		// client functions
		watchClient,

		// collection functions
		aggregateUnstable1,
		aggregateUnstable2,
		aggregateUnstable3,
		aggregateUnstable4,
		aggregateUnstable5,
		aggregateUnstable6,
		distinct,
		find,
		searchIndexes,
		watchCollection,

		// database functions
		watchDatabase,
	}

	// Execute each function
	for _, fn := range functions {
		fn()
	}

	// Disconnect from MongoDB
	err := client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
