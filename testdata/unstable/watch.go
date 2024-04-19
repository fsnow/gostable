package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func watchCollection() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Create a change stream options
	opts := options.ChangeStream().SetFullDocument(options.UpdateLookup)

	// Create a change stream
	changeStream, err := collection.Watch(context.Background(), mongo.Pipeline{}, opts)
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

		fmt.Printf("Operation Type: %v\n", operationType)
		fmt.Printf("Full Document: %v\n", fullDocument)
	}

	if err := changeStream.Err(); err != nil {
		log.Fatal(err)
	}
}
