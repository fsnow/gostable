package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func deleteDocuments() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Delete documents
	filter := bson.M{"key": "value"}
	result, err := collection.DeleteMany(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Deleted %d documents\n", result.DeletedCount)
}
