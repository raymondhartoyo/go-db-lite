package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	dblite "github.com/raymondhartoyo/go-db-lite"
)

func main() {
	ctx := context.Background()
	os.Remove("example.db") // make sure deleting sqlite file used in this example

	db, err := dblite.New("example.db")
	if err != nil {
		log.Fatal(err)
	}

	appState := dblite.State{
		Key:   "last_run",
		Value: time.Now().Format(time.RFC3339),
	}
	if err := db.Save(ctx, appState); err != nil {
		log.Fatal(err)
	}

	lastRun, err := db.Get(ctx, "last_run")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("last_run: %s\n", lastRun.Value)

	bulkStates := []dblite.State{
		{Key: "last_open", Value: time.Now().Format(time.RFC3339)},
		{Key: "count", Value: "420"},
	}
	if err := db.SaveBulk(ctx, bulkStates); err != nil {
		log.Fatal(err)
	}

	lastOpen, err := db.Get(ctx, "last_open")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("last_open: %s\n", lastOpen.Value)

	count, err := db.Get(ctx, "count")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("count: %s\n", count.Value)

	if err := db.Delete(ctx, "count"); err != nil {
		log.Fatal(err)
	}

	afterDelete, err := db.Get(ctx, "count")
	fmt.Printf("Expecting nil value: %v, nil error: %v\n", afterDelete == nil, err == nil)

	os.Remove("example.db") // make sure deleting sqlite file used in this example
}
