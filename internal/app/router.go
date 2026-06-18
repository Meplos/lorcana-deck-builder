// Package app: router.go - Init all HTTP API routes
package app

import (
	"github.com/labstack/echo/v5"
)

func RegisterRoute(e *echo.Echo, container *Container) error {
	group := e.Group("/api/v1")

	// auth
	group.POST("/auth/register", container.AuthHandler.Register)

	// cards
	group.GET("/cards", container.CardHandler.ListCard)

	// collection
	group.POST("/collections", container.CollectionHandler.Create)
	group.GET("/collections", container.CollectionHandler.List)
	group.GET("/collections/export", container.CollectionHandler.Export)
	group.GET("/collections/add-card", container.CollectionHandler.AddCard)

	// deck
	group.POST("/deck/build", container.DeckHandler.BuildDeck)
	group.GET("/deck", container.DeckHandler.List)
	group.POST("/deck", container.DeckHandler.Save)

	return nil
}
