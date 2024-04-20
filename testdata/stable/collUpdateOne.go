package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func updateOne() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Define the filter
	filter := bson.M{
		"name": "John Doe",
	}

	// Define the update operation
	update := bson.M{
		"$set": bson.M{
			"age":    35,
			"city":   "New York",
			"status": "Active",
		},
	}

	// Update the document
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	if result.MatchedCount == 0 {
		fmt.Println("No document found matching the filter")
	} else {
		fmt.Printf("Updated %d document(s)\n", result.ModifiedCount)
	}
}
