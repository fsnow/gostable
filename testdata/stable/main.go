package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func init() {
	// Set up MongoDB client
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
}

/*
Notes:
The following functions are client-side only, server no-ops:
Collection functions Clone(), Database(), Name()
Database functions: Collection(), Client(), Name(), ReadConcern(), ReadPreference()
*/

/*
TODO:
Database.RunCommand
Database.RunCommandCursor
*/

/*
These need "Stable" and "Unstable" test cases:

Collection.Aggregate (TODO)
https://www.mongodb.com/docs/manual/reference/command/aggregate/#stable-api
Aggregate has limitations on stages

Database.CreateCollection
https://www.mongodb.com/docs/manual/reference/command/create/#stable-api
Collection create does not allow 6 fields
Hard for analysis of direct runCommand, think of all the ways that this could be constructed, maybe in a different module

Database.CreateIndex
https://www.mongodb.com/docs/manual/reference/command/createIndexes/#stable-api
createIndexes does not allow 4 fields

Collection.Find() and Collection.FindOne*()
https://www.mongodb.com/docs/manual/reference/command/find/#stable-api
Collection find does not allow 8 fields


*/

// stable
func main() {
	// Array of functions to be called
	functions := []func(){
		// client functions
		listDatabaseNames,
		listDatabases,
		ping,

		// collection functions
		aggregateStable,
		bulkWrite,
		countDocuments,
		deleteMany,
		deleteOne,
		drop,
		estimatedDocumentCount,
		findStable,
		findOneStable,
		findOneAndDelete,
		findOneAndReplace,
		findOneAndUpdate,
		indexes,
		insertMany,
		insertOne,
		replaceOne,
		runCmdCount,
		updateByID,
		updateMany,
		updateOne,

		// database functions
		listCollectionNames,
		listCollections,
		listCollectionSpecifications,
	}

	// Execute each function
	for _, fn := range functions {
		fn()
	}

	// Disconnect from MongoDB
	err := client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
