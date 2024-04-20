package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func insertMany() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Define the documents to be inserted
	documents := []interface{}{
		bson.M{"name": "Alice", "age": 25, "city": "New York"},
		bson.M{"name": "Bob", "age": 30, "city": "London"},
		bson.M{"name": "Claire", "age": 28, "city": "Paris"},
	}

	// Insert the documents
	result, err := collection.InsertMany(context.Background(), documents)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Inserted %d documents:\n", len(result.InsertedIDs))
	for _, insertedID := range result.InsertedIDs {
		fmt.Println(insertedID)
	}
}
