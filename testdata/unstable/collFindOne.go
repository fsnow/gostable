package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindOne() {
	collection := client.Database("mydatabase").Collection("mycollection")

	t := bool(true)
	sec10 := time.Second * 10
	// Create a FindOneOptions struct
	findOneOptions := options.FindOneOptions{
		Max:             bson.M{"field": 100},
		MaxAwaitTime:    &sec10,
		Min:             bson.M{"field": 50},
		NoCursorTimeout: &t,
		OplogReplay:     &t,
		ReturnKey:       &t,
		ShowRecordID:    &t,
	}

	// Find a document with the specified options
	var result bson.M
	err := collection.FindOne(context.Background(), bson.M{}, &findOneOptions).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No document found")
			return
		}
		log.Fatal(err)
	}

	// Print the retrieved document
	fmt.Println(result)
}
