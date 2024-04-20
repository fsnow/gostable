package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func deleteOne() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Define the filter
	filter := bson.M{
		"name": "John Doe",
	}

	// Delete the document that matches the filter
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	// Print the delete result
	fmt.Printf("Deleted %d document(s)\n", result.DeletedCount)
}
