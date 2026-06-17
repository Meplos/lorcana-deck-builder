// Package mongo: Connection script
package mongo

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func Connect(ctx context.Context, databaseName string) (*mongo.Database, error) {
	mongoclient, err := mongo.Connect(ctx, options.Client().ApplyURI(
		BuildConnURI(),
	))
	if err != nil {
		return nil, err
	}
	client = mongoclient
	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, err
	}
	return client.Database(databaseName), nil
}

func BuildConnURI() string {
	return fmt.Sprintf(
		"%s://%s:%s@%s:%s",
		os.Getenv("MONGO_SCHEME"),
		os.Getenv("MONGO_USER"),
		os.Getenv("MONGO_PASS"),
		os.Getenv("MONGO_HOST"),
		os.Getenv("MONGO_PORT"),
	)
}

func Disconnect(ctx context.Context) {
	client.Disconnect(ctx)
}
