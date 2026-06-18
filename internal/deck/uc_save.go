// Package deck : uc_save.go -> Handle Save New deck use case
package deck

import (
	"context"
	"errors"
	"log"

	"github.com/meplos/locana-deck-builder/internal/cards"
	"github.com/meplos/locana-deck-builder/internal/domain"
)

type SaveUseCase struct {
	deckRepo Repository
	cardRepo cards.Repository
}

func NewSaveUC(deckRepo Repository, cardRepo cards.Repository) *SaveUseCase {
	return &SaveUseCase{
		deckRepo: deckRepo,
		cardRepo: cardRepo,
	}
}

type (
	CardInfo struct {
		ID       string
		Quantity int
	}
	SaveDeckInput struct {
		Name     string
		Strategy string
		Size     int
		Deck     []CardInfo
	}
)

func (uc *SaveUseCase) Execute(ctx context.Context, input *SaveDeckInput) error {
	ids := extractIDs(input.Deck)

	cards, err := uc.cardRepo.FindByIDs(ctx, ids)
	if err != nil {
		log.Print(err)
		return err
	}
	deckCards, err := formatCard(cards, input.Deck)
	if err != nil {
		log.Print(err)
		return err
	}
	deck := domain.Deck{
		Size:     input.Size,
		Strategy: input.Strategy,
		Name:     input.Name,
		Deck:     deckCards,
	}

	err = uc.deckRepo.Save(ctx, deck)
	if err != nil {
		log.Print(err)
	}

	return err
}

func extractIDs(cards []CardInfo) []string {
	ids := make([]string, 0)

	for _, c := range cards {
		ids = append(ids, c.ID)
	}

	return ids
}

func formatCard(cards map[string]domain.Card, info []CardInfo) ([]domain.DeckCard, error) {
	result := make([]domain.DeckCard, 0)
	for _, i := range info {
		c, ok := cards[i.ID]
		if !ok {
			log.Printf("Card: %s", i.ID)
			return result, errors.New("InvalidCard")
		}

		result = append(result, buildDeckCard(c, i))
	}

	return result, nil
}

func buildDeckCard(c domain.Card, i CardInfo) domain.DeckCard {
	abilities := make([]domain.CardAbility, 0)
	for _, a := range c.Abilities {
		abilities = append(abilities, domain.CardAbility(a))
	}
	return domain.DeckCard{
		ID:              c.ID,
		Name:            c.Name,
		Set:             c.Set,
		Version:         c.Version,
		Title:           c.Title,
		Cost:            c.Cost,
		Type:            c.Type,
		Number:          c.Number,
		Color:           c.Color,
		Illustrator:     c.Illustrator,
		Lore:            c.Lore,
		Strength:        c.Strength,
		Willpower:       c.Willpower,
		Movement:        c.Movement,
		Ink:             c.Ink,
		Characteristics: c.Characteristics,
		Abilities:       abilities,
		Variants:        c.Variants,
		Rarity:          c.Rarity,
		Language:        c.Language,
		Path:            c.Path,
		Franchise:       c.Franchise,
		Ordinal:         c.Ordinal,
		Formats:         c.Formats,
		Quantity:        i.Quantity,
	}
}
