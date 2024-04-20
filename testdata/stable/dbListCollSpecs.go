package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func listCollectionSpecifications() {
	database := client.Database("mydatabase")

	// List collection specifications
	specs, err := database.ListCollectionSpecifications(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Print the collection specifications
	fmt.Println("Collection specifications:")
	for _, spec := range specs {
		specBytes, err := bson.MarshalExtJSON(spec, true, false)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(specBytes))
	}
}
