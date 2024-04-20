package main

import (
	"context"
	"fmt"
	"log"
)

func listDatabaseNames() {
	// List database names
	names, err := client.ListDatabaseNames(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Print the database names
	fmt.Println("Database names:")
	for _, name := range names {
		fmt.Println(name)
	}
}
