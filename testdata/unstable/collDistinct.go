package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func distinct() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Specify the field for distinct values
	field := "category"

	// Retrieve the distinct values
	values, err := collection.Distinct(context.Background(), field, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	// Print the distinct values
	fmt.Printf("Distinct values for field '%s':\n", field)
	for _, value := range values {
		fmt.Printf("- %v\n", value)
	}
}
