package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func bulkWrite() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Define the bulk write operations
	operations := []mongo.WriteModel{
		// Insert operation
		mongo.NewInsertOneModel().SetDocument(bson.M{
			"name": "John Doe",
			"age":  30,
		}),
		// Update operation
		mongo.NewUpdateOneModel().
			SetFilter(bson.M{"name": "Jane Smith"}).
			SetUpdate(bson.M{
				"$set": bson.M{
					"age": 35,
				},
			}),
		// Delete operation
		mongo.NewDeleteOneModel().SetFilter(bson.M{
			"name": "Bob Johnson",
		}),
	}

	// Execute the bulk write operations
	result, err := collection.BulkWrite(context.Background(), operations)
	if err != nil {
		log.Fatal(err)
	}

	// Print the bulk write result
	fmt.Printf("Bulk Write Result:\n")
	fmt.Printf("  Inserted: %d\n", result.InsertedCount)
	fmt.Printf("  Matched: %d\n", result.MatchedCount)
	fmt.Printf("  Modified: %d\n", result.ModifiedCount)
	fmt.Printf("  Deleted: %d\n", result.DeletedCount)
	fmt.Printf("  Upserted: %d\n", result.UpsertedCount)
}
