package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func aggregateStable() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Define the aggregation pipelines
	pipelines := []mongo.Pipeline{
		// Pipeline 1: Single $match stage
		{{
			{Key: "$match", Value: bson.D{
				{Key: "field1", Value: "value1"},
			}},
		}},

		// Pipeline 2: $group and $project stages
		{{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$field2"},
				{Key: "count", Value: bson.D{
					{Key: "$sum", Value: 1},
				}},
			}},
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "field2", Value: "$_id"},
				{Key: "count", Value: 1},
			}},
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
