package deck

import (
	"context"

	"github.com/labstack/echo/v5"
	collectionInfra "github.com/meplos/locana-deck-builder/internal/collection/infra"
	"github.com/meplos/locana-deck-builder/internal/deck/domain"
	"github.com/meplos/locana-deck-builder/internal/deck/infra"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetUp(ctx context.Context, e *echo.Group, DB *mongo.Database) {
	repo := collectionInfra.NewCollectionRepo(DB)
	ai := infra.NewAIAgent()

	buildUC := domain.CreateBuildDeckUC(ctx, repo, ai)

	Init(e, buildUC)
}
