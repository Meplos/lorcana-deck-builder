package deck

import (
	"context"

	"github.com/meplos/locana-deck-builder/internal/domain"
)

type Repository interface {
	Save(ctx context.Context, d domain.Deck) error
	FindBy(ctx context.Context) ([]domain.Deck, error)
	FindByID(ctx context.Context, ID string) (domain.Deck, error)
	CountBy(ctx context.Context) (int, error)
}
