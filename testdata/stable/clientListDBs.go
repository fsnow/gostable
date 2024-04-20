package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func listDatabases() {
	// List databases
	databases, err := client.ListDatabases(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	// Print the database information
	fmt.Println("Databases:")
	for _, db := range databases.Databases {
		dbBytes, err := bson.MarshalExtJSON(db, true, false)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(dbBytes))
	}

	// Print the total size of databases
	fmt.Printf("Total size: %d bytes\n", databases.TotalSize)
}
