// CountDocuments executes an aggregation that uses only $match and $group stages
package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func countDocuments() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Define the filter
	filter := bson.M{
		"age": bson.M{
			"$gte": 18,
		},
	}

	// Count the documents that match the filter
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	// Print the count
	fmt.Printf("Number of documents matching the filter: %d\n", count)
}
