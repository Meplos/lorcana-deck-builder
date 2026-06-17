package cards

import (
	"context"

	"github.com/meplos/locana-deck-builder/internal/domain"
)

type ListFilter struct {
	Offset int64
	Limit  int64
	Search string
	Color  domain.InkColor
}

type Repository interface {
	ListCards(ctx context.Context, f ListFilter) []domain.Card
	CountCards(ctx context.Context, f ListFilter) int64
	FindByID(ctx context.Context, ID string) (domain.Card, error)
}
