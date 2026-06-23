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

	abilities := make([]domain.CardAbility, len(card.Abilities))
	for i, a := range card.Abilities {
		abilities[i] = domain.CardAbility{
			Name:        a.Name,
			Description: a.Description,
		}
	}

	return uc.collRepo.AddCardToCollection(ctx, collID, domain.CollectionCard{
		ID:              card.ID,
		Name:            card.Name,
		Set:             card.Set,
		Version:         card.Version,
		Title:           card.Title,
		Cost:            card.Cost,
		Type:            card.Type,
		Number:          card.Number,
		Color:           card.Color,
		Illustrator:     card.Illustrator,
		Lore:            card.Lore,
		Strength:        card.Strength,
		Willpower:       card.Willpower,
		Movement:        card.Movement,
		Ink:             card.Ink,
		Characteristics: card.Characteristics,
		Abilities:       abilities,
		Variants:        card.Variants,
		Rarity:          card.Rarity,
		Language:        card.Language,
		Path:            card.Path,
		Franchise:       card.Franchise,
		Ordinal:         card.Ordinal,
		Formats:         card.Formats,
		Quantity:        quantity,
	})
}
