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
		aggregateStable,
		findDocuments,
		deleteDocuments,
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
