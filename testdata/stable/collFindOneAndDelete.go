package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func findOneAndDelete() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Define the filter
	filter := bson.M{
		"name": "John Doe",
	}

	// Find and delete the document
	var deletedDocument bson.M
	err := collection.FindOneAndDelete(context.Background(), filter).Decode(&deletedDocument)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			fmt.Println("No document found matching the filter")
		} else {
			log.Fatal(err)
		}
		return
	}

	fmt.Println("Deleted document:")
	fmt.Println(deletedDocument)
}
