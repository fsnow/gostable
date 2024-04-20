package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func findOneStable() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Define the filter
	filter := bson.M{
		"name": "John Doe",
	}

	// Find the document
	var result bson.M
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			fmt.Println("No document found matching the filter")
		} else {
			log.Fatal(err)
		}
		return
	}

	fmt.Println("Found document:")
	fmt.Println(result)
}
