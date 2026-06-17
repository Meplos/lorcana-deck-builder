package collection

import (
	"context"

	"github.com/meplos/locana-deck-builder/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	List(ctx context.Context, f ListFilter) []domain.Collection
	Count(ctx context.Context, f ListFilter) int
	FindOrCreate(ctx context.Context, name string) (domain.Collection, error)
	FindCollectionByName(ctx context.Context, name string) (domain.Collection, error)
	AddCardToCollection(ctx context.Context, collID primitive.ObjectID, card domain.CollectionCard) error
}
