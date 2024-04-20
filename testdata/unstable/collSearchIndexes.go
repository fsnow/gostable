package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func searchIndexes() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Get the Atlas Search indexes for the collection
	searchIndexView := collection.SearchIndexes()

	// List all the indexes
	cursor, err := searchIndexView.List(context.Background(), options.SearchIndexes().SetName("mySearchIndex"))
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	// Iterate over the indexes
	for cursor.Next(context.Background()) {
		var searchIndex bson.M
		err := cursor.Decode(&searchIndex)
		if err != nil {
			log.Fatal(err)
		}

		// Print the index information
		fmt.Printf("Index: %v\n", searchIndex)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}
