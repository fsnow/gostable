package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func watchClient() {
	// Create a change stream options
	opts := options.ChangeStream().SetFullDocument(options.UpdateLookup)

	// Create a change stream
	changeStream, err := client.Watch(context.Background(), mongo.Pipeline{}, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer changeStream.Close(context.Background())

	// Iterate over the change stream
	for changeStream.Next(context.Background()) {
		var change bson.M
		err := changeStream.Decode(&change)
		if err != nil {
			log.Fatal(err)
		}

		// Process the change event
		fmt.Printf("Change event: %v\n", change)

		// Access specific fields from the change event
		operationType := change["operationType"]
		fullDocument := change["fullDocument"]
		documentKey := change["documentKey"]

		fmt.Printf("Operation Type: %v\n", operationType)
		fmt.Printf("Full Document: %v\n", fullDocument)
		fmt.Printf("Document Key: %v\n", documentKey)
	}

	if err := changeStream.Err(); err != nil {
		log.Fatal(err)
	}
}
