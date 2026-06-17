package collection

import (
	"context"
	"errors"

	"github.com/meplos/locana-deck-builder/internal/cards"
	"github.com/meplos/locana-deck-builder/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddCardUseCase struct {
	collRepo Repository
	cardRepo cards.Repository
}

func NewAddCardUC(collRepo Repository, cardRepo cards.Repository) *AddCardUseCase {
	return &AddCardUseCase{
		collRepo: collRepo,
		cardRepo: cardRepo,
	}
}

func (uc *AddCardUseCase) AddCardToCollection(ctx context.Context, collID primitive.ObjectID, cardID string, quantity int) error {
	card, err := uc.cardRepo.FindByID(ctx, cardID)
	if err != nil {
		return err
	}
	if card.ID == "" {
		card.ID = cardID
	}
	if card.ID == "" {
		return errors.New("card not found")
	}

	return uc.collRepo.AddCardToCollection(ctx, collID, domain.CollectionCard{
		Quantity: quantity,
		// Card:     card,
	})
}
