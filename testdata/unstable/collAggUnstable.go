package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Pipeline as slice of bson.D
func aggregateUnstable1() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Define the aggregation pipelines
	pipelines := []bson.D{
		{{"$currentOp", bson.D{}}},
	}

	// Iterate over the pipelines and execute the aggregation
	for _, pipeline := range pipelines {
		cursor, err := collection.Aggregate(context.Background(), pipeline)
		if err != nil {
			log.Fatal(err)
		}
		defer cursor.Close(context.Background())

		// Iterate over the aggregation results
		for cursor.Next(context.Background()) {
			var result bson.M
			err := cursor.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(result)
		}

		if err := cursor.Err(); err != nil {
			log.Fatal(err)
		}
	}
}

// Pipeline as slice of bson.M
func aggregateUnstable2() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Define the aggregation pipelines
	pipelines := []bson.M{
		{"$currentOp": bson.M{}},
	}

	// Iterate over the pipelines and execute the aggregation
	for _, pipeline := range pipelines {
		cursor, err := collection.Aggregate(context.Background(), pipeline)
		if err != nil {
			log.Fatal(err)
		}
		defer cursor.Close(context.Background())

		// Iterate over the aggregation results
		for cursor.Next(context.Background()) {
			var result bson.M
			err := cursor.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(result)
		}

		if err := cursor.Err(); err != nil {
			log.Fatal(err)
		}
	}
}

// Pipeline as slice of interface{}, mixed types
func aggregateUnstable3() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Define the aggregation pipelines
	pipelines := []interface{}{
		bson.D{{"$currentOp", bson.D{}}},
		bson.M{"$indexStats": bson.M{}},
	}

	// Iterate over the pipelines and execute the aggregation
	for _, pipeline := range pipelines {
		cursor, err := collection.Aggregate(context.Background(), pipeline)
		if err != nil {
			log.Fatal(err)
		}
		defer cursor.Close(context.Background())

		// Iterate over the aggregation results
		for cursor.Next(context.Background()) {
			var result bson.M
			err := cursor.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(result)
		}

		if err := cursor.Err(); err != nil {
			log.Fatal(err)
		}
	}
}

// Pipeline as slice of primitive.D
func aggregateUnstable4() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Define the aggregation pipelines
	pipelines := []primitive.D{
		{{"$currentOp", bson.D{}}},
	}

	// Iterate over the pipelines and execute the aggregation
	for _, pipeline := range pipelines {
		cursor, err := collection.Aggregate(context.Background(), pipeline)
		if err != nil {
			log.Fatal(err)
		}
		defer cursor.Close(context.Background())

		// Iterate over the aggregation results
		for cursor.Next(context.Background()) {
			var result bson.M
			err := cursor.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(result)
		}

		if err := cursor.Err(); err != nil {
			log.Fatal(err)
		}
	}
}

// Pipeline as slice of mongo.Pipeline
func aggregateUnstable5() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Define the aggregation pipelines
	pipelines := []mongo.Pipeline{
		{{
			{Key: "$currentOp", Value: bson.D{}},
		}},
	}

	// Iterate over the pipelines and execute the aggregation
	for _, pipeline := range pipelines {
		cursor, err := collection.Aggregate(context.Background(), pipeline)
		if err != nil {
			log.Fatal(err)
		}
		defer cursor.Close(context.Background())

		// Iterate over the aggregation results
		for cursor.Next(context.Background()) {
			var result bson.M
			err := cursor.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(result)
		}

		if err := cursor.Err(); err != nil {
			log.Fatal(err)
		}
	}
}

// Pipeline where key is in a separate variable
func aggregateUnstable6() {
	collection := client.Database("mydatabase").Collection("mycollection")

	key := "$currentOp"

	// Define the aggregation pipelines
	pipelines := []mongo.Pipeline{
		{{
			{Key: key, Value: bson.D{}},
		}},
	}

	// Iterate over the pipelines and execute the aggregation
	for _, pipeline := range pipelines {
		cursor, err := collection.Aggregate(context.Background(), pipeline)
		if err != nil {
			log.Fatal(err)
		}
		defer cursor.Close(context.Background())

		// Iterate over the aggregation results
		for cursor.Next(context.Background()) {
			var result bson.M
			err := cursor.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(result)
		}

		if err := cursor.Err(); err != nil {
			log.Fatal(err)
		}
	}
}
