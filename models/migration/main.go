package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ray1422/dcard-backend-2023/utils/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

// Migrate Migrate
func Migrate() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = db.MongoDB().CreateCollection(ctx, "list")
	if err != nil {
		return
	}
	err = db.MongoDB().CreateCollection(ctx, "list_node")
	if err != nil {
		return
	}
	indexName, err := db.MongoDB().Collection("list").Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{
		{Key: "user_id", Value: 1},
		{Key: "key", Value: 1},
	}})
	if err != nil {
		return
	}
	fmt.Printf("index `%s` for `list` has been created.\n", indexName)
	indexName, err = db.MongoDB().Collection("list_node").Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{
		{Key: "expire_on", Value: 1},
	}, Options: options.Index().SetExpireAfterSeconds(0),
	})
	fmt.Printf("index `%s` for `list_node` has been created.\n", indexName)
	if err != nil {
		return
	}
	return nil
}

// Rollback Rollback
func Rollback() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = db.MongoDB().Collection("list").Drop(ctx)
	if err != nil {
		return
	}
	err = db.MongoDB().Collection("list_node").Drop(ctx)
	if err != nil {
		return
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <m|r>\n", os.Args[0])
		os.Exit(2)
		return
	}
	var err error = nil
	if os.Args[1] != "r" {
		err = Migrate()
	} else {
		err = Rollback()
	}
	if err != nil {
		fmt.Println("error occurred:\n", err)
	} else {
		fmt.Println("operation successfully completed.")
	}
}
