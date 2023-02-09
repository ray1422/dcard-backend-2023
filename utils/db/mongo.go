package db

import (
	"context"
	"time"

	"github.com/ray1422/dcard-backend-2023/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_client, err := mongo.Connect(ctx, options.Client().ApplyURI(utils.Getenv("MONGODB_URL", "mongodb://localhost:27017")))
	if err != nil {
		panic(err)
	}
	client = _client
}

// Mongo returns a MongoDB client instance
func Mongo() *mongo.Client {
	return client
}

// MongoDB returns the ptr to the database
func MongoDB() *mongo.Database {
	return client.Database(utils.Getenv("MONGODB_NAME", "asdf"))
}
