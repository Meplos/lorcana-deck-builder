package cards

import (
	"context"

	"github.com/labstack/echo/v5"
	"github.com/meplos/locana-deck-builder/internal/cards/domain"
	"github.com/meplos/locana-deck-builder/internal/cards/infra"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetUp(ctx context.Context, e *echo.Group, DB *mongo.Database) {
	repo := infra.SetupMongoCardRepository(DB)
	uc := domain.CreateListCardUC(ctx, repo)
	Init(e, uc)
}
