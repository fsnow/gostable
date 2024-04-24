package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func runCmdCount() {
	db := client.Database("mydatabase")
	collName := "mycollection"

	// Create the "count" command
	countCommand := bson.D{
		{Key: "count", Value: collName},
		{Key: "query", Value: bson.M{}}, // Empty query to count all documents
	}

	// Execute the "count" command
	var result bson.M
	err := db.RunCommand(context.Background(), countCommand).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	// Extract the count from the result
	count := result["n"].(int64)
	fmt.Printf("Collection %s has %d documents\n", collName, count)
}
