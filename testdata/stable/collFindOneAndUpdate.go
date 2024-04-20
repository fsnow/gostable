package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func findOneAndUpdate() {
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

	// Define the options for the operation
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	// Find and update the document
	var result bson.M
	err := collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&result)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			fmt.Println("No document found matching the filter")
		} else {
			log.Fatal(err)
		}
		return
	}

	fmt.Println("Updated document:")
	fmt.Println(result)
}
