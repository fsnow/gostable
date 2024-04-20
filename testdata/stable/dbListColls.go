package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func listCollections() {
	database := client.Database("mydatabase")

	// List collections
	cursor, err := database.ListCollections(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	// Iterate over the collections
	for cursor.Next(context.Background()) {
		var collection bson.M
		err := cursor.Decode(&collection)
		if err != nil {
			log.Fatal(err)
		}

		// Print the collection information
		collectionBytes, err := bson.MarshalExtJSON(collection, true, false)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(collectionBytes))
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}
