package main

import (
	"context"
	"fmt"
	"log"
)

func estimatedDocumentCount() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Get the estimated document count
	estimatedCount, err := collection.EstimatedDocumentCount(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Estimated number of documents in the collection: %d\n", estimatedCount)
}
