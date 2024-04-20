package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func indexes() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Get the index view of the collection
	indexView := collection.Indexes()

	// List all the indexes
	cursor, err := indexView.List(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	// Iterate over the indexes
	for cursor.Next(context.Background()) {
		var index bson.M
		err := cursor.Decode(&index)
		if err != nil {
			log.Fatal(err)
		}

		// Print the index information
		fmt.Printf("Index: %v\n", index)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}
