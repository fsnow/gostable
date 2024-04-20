package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func updateByID() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Define the document ID
	documentID, err := primitive.ObjectIDFromHex("60f1c1a1b3f7d64fdcf9b5e7")
	if err != nil {
		log.Fatal(err)
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
	result, err := collection.UpdateByID(context.Background(), documentID, update)
	if err != nil {
		log.Fatal(err)
	}

	if result.MatchedCount == 0 {
		fmt.Println("No document found with the specified ID")
	} else {
		fmt.Printf("Updated %d document(s)\n", result.ModifiedCount)
	}
}
