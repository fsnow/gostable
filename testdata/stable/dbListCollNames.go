package main

import (
	"context"
	"fmt"
	"log"
)

func listCollectionNames() {
	database := client.Database("mydatabase")

	// List collection names
	names, err := database.ListCollectionNames(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Print the collection names
	fmt.Println("Collection names:")
	for _, name := range names {
		fmt.Println(name)
	}
}
