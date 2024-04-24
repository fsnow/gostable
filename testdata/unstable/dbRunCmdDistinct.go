package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func runCmdDistinct1() {
	db := client.Database("mydatabase")
	collName := "mycollection"

	// Specify the field for distinct values
	field := "category"

	// Create the "distinct" command
	distinctCommand := bson.D{
		{Key: "distinct", Value: collName},
		{Key: "key", Value: field},
	}

	// Execute the "distinct" command
	var result bson.M
	err := db.RunCommand(context.Background(), distinctCommand).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	// Extract the distinct values from the result
	values := result["values"].(bson.A)
	fmt.Printf("Distinct values for field '%s' in collection '%s':\n", field, collName)
	for _, value := range values {
		fmt.Println(value)
	}
}

func runCmdDistinct2() {
	db := client.Database("mydatabase")
	collName := "mycollection"

	// Specify the field for distinct values
	field := "category"

	// Create the "distinct" command
	distinctCommand := bson.D{
		bson.E{Key: "distinct", Value: "value1"},
		bson.E{Key: "field2", Value: 42},
		bson.E{Key: "field3", Value: true},
	}

	// Execute the "distinct" command
	var result bson.M
	err := db.RunCommand(context.Background(), distinctCommand).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	// Extract the distinct values from the result
	values := result["values"].(bson.A)
	fmt.Printf("Distinct values for field '%s' in collection '%s':\n", field, collName)
	for _, value := range values {
		fmt.Println(value)
	}
}
