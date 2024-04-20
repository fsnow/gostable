package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func updateMany() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Define the filter
	filter := bson.M{
		"age": bson.M{
			"$gte": 30,
		},
	}

	// Define the update operation
	update := bson.M{
		"$set": bson.M{
			"status": "Senior",
		},
	}

	// Update the documents
	result, err := collection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	if result.MatchedCount == 0 {
		fmt.Println("No documents found matching the filter")
	} else {
		fmt.Printf("Updated %d document(s)\n", result.ModifiedCount)
	}
}
