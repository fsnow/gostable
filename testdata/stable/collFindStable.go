package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func findStable() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Find documents
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	// Iterate over the documents
	for cursor.Next(context.Background()) {
		var document bson.M
		err := cursor.Decode(&document)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(document)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}
