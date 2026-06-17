package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/meplos/locana-deck-builder/internal/app"
	"github.com/meplos/locana-deck-builder/internal/database/mongo"
)

func main() {
	_ = godotenv.Load()
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS("*"))

	ctx := context.Background()

	DB, err := mongo.Connect(ctx, os.Getenv("DB_NAME"))
	if err != nil {
		e.Logger.Error("error while connecting DB", err)
		panic(err)
	}

	container, _ := app.NewContainer(DB)
	if err := app.RegisterRoute(e, container); err != nil {
		e.Logger.Error("error while registering routes", err)
		panic(err)
	}

	defer mongo.Disconnect(ctx)
	if err := e.Start(":9090"); err != nil {
		e.Logger.Error("error while serving", err)
		panic(err)
	}
}
