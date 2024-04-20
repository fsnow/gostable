package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func replaceOne() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Define the filter
	filter := bson.M{
		"name": "John Doe",
	}

	// Define the replacement document
	replacement := bson.M{
		"name":    "John Doe",
		"age":     35,
		"city":    "Los Angeles",
		"country": "USA",
	}

	// Replace the document
	result, err := collection.ReplaceOne(context.Background(), filter, replacement)
	if err != nil {
		log.Fatal(err)
	}

	if result.MatchedCount == 0 {
		fmt.Println("No document found matching the filter")
	} else {
		fmt.Printf("Replaced %d document(s)\n", result.ModifiedCount)
	}
}
