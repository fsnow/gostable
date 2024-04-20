package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func insertOne() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Define the document to be inserted
	document := bson.M{
		"name":    "John Doe",
		"age":     30,
		"city":    "New York",
		"country": "USA",
	}

	// Insert the document
	result, err := collection.InsertOne(context.Background(), document)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted document ID:", result.InsertedID)
}
