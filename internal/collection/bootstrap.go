package collection

import (
	"context"

	"github.com/labstack/echo/v5"
	"github.com/meplos/locana-deck-builder/internal/collection/domain"
	"github.com/meplos/locana-deck-builder/internal/collection/infra"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetUp(ctx context.Context, g *echo.Group, DB *mongo.Database) {
	collectionRepo := infra.NewCollectionRepo(DB)
	cardRepo := infra.NewCardRepo(DB)

	listUC := domain.NewListUC(ctx, collectionRepo)
	createUC := domain.NewCreateUC(ctx, collectionRepo)
	addCardUC := domain.NewAddCardUC(ctx, collectionRepo, cardRepo)
	exportUC := domain.NewExportUC(ctx, collectionRepo)

	Init(ctx, g, listUC, createUC, addCardUC, exportUC)
}
