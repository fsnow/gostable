package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func aggregateUnstable() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Define the aggregation pipelines
	pipelines := []mongo.Pipeline{
		// Pipeline 1: Single $match stage
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
