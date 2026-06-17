package domain

import (
	"context"
	"errors"

	"github.com/meplos/locana-deck-builder/internal/collection/infra"
	"github.com/meplos/locana-deck-builder/internal/db/mongo/schema"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddCardUseCase struct {
	ctx      context.Context
	collRepo infra.CollectionRepository
	cardRepo infra.CardRepository
}

func NewAddCardUC(ctx context.Context, collRepo infra.CollectionRepository, cardRepo infra.CardRepository) *AddCardUseCase {
	return &AddCardUseCase{
		ctx:      ctx,
		collRepo: collRepo,
		cardRepo: cardRepo,
	}
}

func (uc *AddCardUseCase) AddCardToCollection(collID primitive.ObjectID, cardID string, quantity int) error {
	card, err := uc.cardRepo.FindById(uc.ctx, cardID)
	if err != nil {
		return err
	}
	if card.ID == "" {
		card.ID = cardID
	}
	if card.ID == "" {
		return errors.New("card not found")
	}

	return uc.collRepo.AddCardToCollection(uc.ctx, collID, schema.OwnedCard{
		Card:     card,
		Quantity: quantity,
	})
}
