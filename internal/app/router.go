// Package app: router.go - Init all HTTP API routes
package app

import (
	"github.com/labstack/echo/v5"
)

func RegisterRoute(e *echo.Echo, container *Container) error {
	group := e.Group("/api/v1")

	// auth
	group.POST("/auth/register", container.AuthHandler.Register)
	group.POST("/auth/login", container.AuthHandler.Login)
	group.POST("/logout", container.AuthHandler.Logout)

	// cards
	group.GET("/cards", container.CardHandler.ListCard)

	// collection
	group.POST("/collections", container.CollectionHandler.Create, container.AuthHandler.IsConnected)
	group.GET("/collections", container.CollectionHandler.List, container.AuthHandler.IsConnected)
	group.GET("/collections/export", container.CollectionHandler.Export, container.AuthHandler.IsConnected)
	group.POST("/collections/add-card", container.CollectionHandler.AddCard, container.AuthHandler.IsConnected)

	// deck
	group.POST("/deck/build", container.DeckHandler.BuildDeck, container.AuthHandler.IsConnected)
	group.GET("/deck", container.DeckHandler.List, container.AuthHandler.IsConnected)
	group.POST("/deck", container.DeckHandler.Save, container.AuthHandler.IsConnected)

	return nil
}
