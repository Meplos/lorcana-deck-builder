package main

import (
	"context"
	"log"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/meplos/locana-deck-builder/internal/cards"
	"github.com/meplos/locana-deck-builder/internal/collection"
	"github.com/meplos/locana-deck-builder/internal/deck"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	_ = godotenv.Load()
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS("*"))
	ctx := context.Background()
	mongoURI := "mongodb://root:example@localhost:9999"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	db := client.Database("lorcana")

	api := e.Group("/api/v1")

	cards.SetUp(ctx, api, db)
	collection.SetUp(ctx, api, db)
	deck.SetUp(ctx, api, db)

	if err := e.Start(":9090"); err != nil {
		log.Fatal(err)
	}
}
