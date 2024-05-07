package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func find1() {
	collection := client.Database("mydatabase").Collection("mycollection")

	findOptions := options.Find()
	findOptions.SetShowRecordID(true)
	findOptions.SetNoCursorTimeout(true)

	cursor, err := collection.Find(context.Background(), bson.M{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var document bson.M
		err := cursor.Decode(&document)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(document)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}

func find2() {
	collection := client.Database("mydatabase").Collection("mycollection")

	show := bool(true)
	findOptions := options.FindOptions{
		ShowRecordID: &show,
	}

	cursor, err := collection.Find(context.Background(), bson.M{}, &findOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var document bson.M
		err := cursor.Decode(&document)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(document)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}

func find3() {
	collection := client.Database("mydatabase").Collection("mycollection")

	show := bool(true)
	findOptions := &options.FindOptions{
		Sort:         bson.D{{"name", 1}},
		ShowRecordID: &show,
	}

	cursor, err := collection.Find(context.Background(), bson.M{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var document bson.M
		err := cursor.Decode(&document)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(document)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}

func find4() {
	collection := client.Database("mydatabase").Collection("mycollection")

	findOptions := options.Find()
	findOptions.SetShowRecordID(true)
	findOptions.SetSort(bson.D{{"name", 1}})

	cursor, err := collection.Find(context.Background(), bson.M{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var document bson.M
		err := cursor.Decode(&document)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(document)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}

func find5() {
	collection := client.Database("mydatabase").Collection("mycollection")

	show := bool(true)
	findOptions := options.FindOptions{
		ShowRecordID: &show,
		Sort:         bson.D{{"name", 1}},
		Projection: bson.D{
			{"name", 1},
			{"age", 1},
			{"_id", 0},
		},
		Hint: bson.D{{"category", 1}},
	}

	cursor, err := collection.Find(context.Background(), bson.M{}, &findOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var document bson.M
		err := cursor.Decode(&document)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(document)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}

func find6() {
	collection := client.Database("mydatabase").Collection("mycollection")

	findOptions := &options.FindOptions{}
	findOptions.SetShowRecordID(true)
	findOptions.SetSort(bson.D{{"name", 1}})
	findOptions.SetProjection(bson.D{
		{"name", 1},
		{"age", 1},
		{"_id", 0},
	})
	findOptions.SetSkip(20)
	findOptions.SetHint(bson.D{{"category", 1}})
	findOptions.SetBatchSize(5)

	cursor, err := collection.Find(context.Background(), bson.M{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var document bson.M
		err := cursor.Decode(&document)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(document)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}

func find7() {
	collection := client.Database("mydatabase").Collection("mycollection")

	// Create a tailable cursor
	findOptions := options.Find().SetCursorType(options.TailableAwait)
	cur, err := collection.Find(context.Background(), bson.M{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())

	// Iterate over the tailable cursor
	for cur.Next(context.Background()) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
}
