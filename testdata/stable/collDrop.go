package main

import (
	"context"
	"fmt"
	"log"
)

func drop() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Drop the collection
	err := collection.Drop(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Collection dropped successfully")
}
