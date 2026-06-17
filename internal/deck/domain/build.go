package domain

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"

	collectionInfra "github.com/meplos/locana-deck-builder/internal/collection/infra"
	"github.com/meplos/locana-deck-builder/internal/deck/infra"
	"github.com/meplos/locana-deck-builder/internal/ink"
)

type BuildDeckUseCase struct {
	ctx  context.Context
	repo collectionInfra.CollectionRepository
	ai   infra.AIAgent
}

func CreateBuildDeckUC(ctx context.Context, repo collectionInfra.CollectionRepository, ai infra.AIAgent) *BuildDeckUseCase {
	return &BuildDeckUseCase{
		ctx:  ctx,
		repo: repo,
		ai:   ai,
	}
}

type BuildDeckParams struct {
	CollectionName string
	InkColors      []ink.InkColor
	Level          string
}

type Deck struct {
	Size     int         `json:"size"`
	Strategy string      `json:"strategy"`
	Name     string      `json:"name"`
	Deck     []OwnedCard `json:"deck"`
}

func (uc *BuildDeckUseCase) Build(params BuildDeckParams) (Deck, error) {
	if strings.TrimSpace(params.CollectionName) == "" {
		return Deck{}, errors.New("InvalidColName")
	}

	collection, err := uc.repo.FindCollectionByName(uc.ctx, params.CollectionName)
	if err != nil {
		return Deck{}, err
	}

	cards := filterCardsByColor(collection.Cards, params.InkColors)

	formated := make([]infra.CardSummary, 0)
	for _, c := range cards {
		formated = append(formated, FormatCard(c))
	}

	response, err := uc.ai.Send(infra.DeckRequest{
		Cards: formated,
		Level: params.Level,
	})
	if err != nil {
		log.Fatal(err)
		return Deck{}, err
	}

	deck := make([]OwnedCard, 0)

	for _, c := range response.Deck {
		deck = append(deck, OwnedCard{
			Card:     collection.Cards[c.ID].Card,
			FilePath: fmt.Sprintf("https://cdn.dreamborn.ink/images/fr/cards/%s", c.ID),
			Quantity: c.Quantity,
		})
	}

	return Deck{
		Size:     response.Size,
		Strategy: response.Strategy,
		Name:     response.Name,
		Deck:     deck,
	}, nil
}

func filterCardsByColor(cards map[string]collectionInfra.OwnedCard, colors []ink.InkColor) []collectionInfra.OwnedCard {
	result := make([]collectionInfra.OwnedCard, 0)

	for _, c := range cards {
		added := false
		for _, color := range colors {
			if slices.Contains(c.Card.Color, color) && !added {
				result = append(result, c)
				added = true
			}
		}

	}

	return result
}

func FormatCard(c collectionInfra.OwnedCard) infra.CardSummary {
	abilities := make([]struct {
		Name        string
		Description string
	}, 0)

	for _, a := range c.Card.Abilities {
		abilities = append(abilities, struct {
			Name        string
			Description string
		}{
			Name:        a.Name,
			Description: a.Description,
		})
	}
	return infra.CardSummary{
		ID:              c.Card.ID,
		Name:            c.Card.Name,
		Title:           c.Card.Title,
		Cost:            c.Card.Cost,
		Type:            c.Card.Type,
		Color:           c.Card.Color,
		Lore:            c.Card.Lore,
		Strength:        c.Card.Strength,
		Willpower:       c.Card.Willpower,
		Movement:        c.Card.Movement,
		Ink:             c.Card.Ink,
		Characteristics: c.Card.Characteristics,
		Abilities: []struct {
			Name        string "json:\"name\""
			Description string "json:\"description\""
		}(abilities),
		Quantity: c.Quantity,
	}
}
