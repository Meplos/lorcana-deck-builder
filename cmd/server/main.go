package main

import (
	"context"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/meplos/locana-deck-builder/internal/app"
	"github.com/meplos/locana-deck-builder/internal/database/mongo"
)

func corsOrigins() []string {
	if raw := os.Getenv("CORS_ORIGINS"); raw != "" {
		origins := strings.Split(raw, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
		}
		return origins
	}

	return []string{
		"http://localhost:5173",
		"http://127.0.0.1:5173",
		"http://localhost:4173",
		"http://127.0.0.1:4173",
	}
}

func main() {
	_ = godotenv.Load()
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     corsOrigins(),
		AllowCredentials: true,
	}))

	ctx := context.Background()

	DB, err := mongo.Connect(ctx, os.Getenv("DB_NAME"))
	if err != nil {
		e.Logger.Error("error while connecting DB", "error", err)
		panic(err)
	}

	container, _ := app.NewContainer(DB)
	if err := app.RegisterRoute(e, container); err != nil {
		e.Logger.Error("error while registering routes", "error", err)
		panic(err)
	}

	defer mongo.Disconnect(ctx)
	if err := e.Start(":9090"); err != nil {
		e.Logger.Error("error while serving", "error", err)
		panic(err)
	}
}
