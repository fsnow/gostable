package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func findOneAndReplace() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Define the filter
	filter := bson.M{
		"name": "John Doe",
	}

	// Define the replacement document
	replacement := bson.M{
		"name":    "John Doe",
		"age":     35,
		"country": "USA",
	}

	// Define the options for the operation
	opts := options.FindOneAndReplace().SetReturnDocument(options.After)

	// Find and replace the document
	var result bson.M
	err := collection.FindOneAndReplace(context.Background(), filter, replacement, opts).Decode(&result)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			fmt.Println("No document found matching the filter")
		} else {
			log.Fatal(err)
		}
		return
	}

	fmt.Println("Replaced document:")
	fmt.Println(result)
}
