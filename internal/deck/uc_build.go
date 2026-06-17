package deck

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/meplos/locana-deck-builder/internal/collection"
	"github.com/meplos/locana-deck-builder/internal/domain"
)

type BuildDeckUseCase struct {
	repo collection.Repository
	ai   AIAgent
}

func CreateBuildDeckUC(repo collection.Repository, ai AIAgent) *BuildDeckUseCase {
	return &BuildDeckUseCase{
		repo: repo,
		ai:   ai,
	}
}

type BuildDeckParams struct {
	CollectionName string
	InkColors      []domain.InkColor
	Level          string
}

type (
	DeckCard struct {
		ID       string
		Name     string
		Title    string
		Number   string
		Set      string
		Type     string
		Rarity   string
		Quantity int
		Filepath string
	}
	Deck struct {
		Size     int
		Strategy string
		Name     string
		Deck     []DeckCard
	}
)

func (uc *BuildDeckUseCase) Build(ctx context.Context, params BuildDeckParams) (Deck, error) {
	if strings.TrimSpace(params.CollectionName) == "" {
		return Deck{}, errors.New("InvalidColName")
	}

	collection, err := uc.repo.FindCollectionByName(ctx, params.CollectionName)
	if err != nil {
		return Deck{}, err
	}

	cards := filterCardsByColor(collection.Cards, params.InkColors)

	formated := make([]CardSummary, 0)
	for _, c := range cards {
		formated = append(formated, FormatCard(c))
	}

	response, err := uc.ai.Send(DeckRequest{
		Cards: formated,
		Level: params.Level,
	})
	if err != nil {
		log.Fatal(err)
		return Deck{}, err
	}

	deck := make([]DeckCard, 0)

	for _, c := range response.Deck {
		originalCard := collection.Cards[c.ID]
		deck = append(deck, DeckCard{
			ID:       originalCard.ID,
			Name:     originalCard.Name,
			Title:    originalCard.Title,
			Type:     originalCard.Type,
			Number:   originalCard.Number,
			Rarity:   originalCard.Rarity,
			Quantity: c.Quantity,
			Filepath: fmt.Sprintf("%s/%s", os.Getenv("CDN_URL"), originalCard.ID),
		})
	}

	return Deck{
		Size:     response.Size,
		Strategy: response.Strategy,
		Name:     response.Name,
		Deck:     deck,
	}, nil
}

func filterCardsByColor(cards map[string]domain.CollectionCard, colors []domain.InkColor) []domain.CollectionCard {
	result := make([]domain.CollectionCard, 0)

	for _, c := range cards {
		added := false
		for _, color := range colors {
			if slices.Contains(c.Color, color) && !added {
				result = append(result, c)
				added = true
			}
		}

	}

	return result
}

func FormatCard(c domain.CollectionCard) CardSummary {
	abilities := make([]struct {
		Name        string
		Description string
	}, 0)

	for _, a := range c.Abilities {
		abilities = append(abilities, struct {
			Name        string
			Description string
		}{
			Name:        a.Name,
			Description: a.Description,
		})
	}
	return CardSummary{
		ID:              c.ID,
		Name:            c.Name,
		Title:           c.Title,
		Cost:            c.Cost,
		Type:            c.Type,
		Color:           c.Color,
		Lore:            c.Lore,
		Strength:        c.Strength,
		Willpower:       c.Willpower,
		Movement:        c.Movement,
		Ink:             c.Ink,
		Characteristics: c.Characteristics,
		Abilities: []struct {
			Name        string "json:\"name\""
			Description string "json:\"description\""
		}(abilities),
		Quantity: c.Quantity,
	}
}
