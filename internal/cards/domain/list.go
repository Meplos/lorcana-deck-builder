package domain

import (
	"context"

	"github.com/meplos/locana-deck-builder/internal/cards/infra"
	"github.com/meplos/locana-deck-builder/internal/ink"
)

type ListCardUC struct {
	ctx            context.Context
	cardRepository infra.CardRepository
}

func CreateListCardUC(ctx context.Context, repo infra.CardRepository) *ListCardUC {
	return &ListCardUC{
		ctx:            ctx,
		cardRepository: repo,
	}
}

type ListCardParams struct {
	Page   int64
	Limit  int64
	Search string
	Color  string
}

type ListCardReponse struct {
	Page  int64
	Limit int64
	Docs  []Card
	Total int64
	Size  int64
}

func (uc *ListCardUC) ListCards(params ListCardParams) ListCardReponse {
	count := uc.cardRepository.CountCards(uc.ctx, infra.ListFilter{
		Limit:  params.Limit,
		Offset: int64(params.Page-1) * params.Limit,
		Search: params.Search,
		Color:  ink.InkColor(params.Color),
	})

	hits := uc.cardRepository.ListCards(uc.ctx, infra.ListFilter{
		Limit:  params.Limit,
		Offset: int64(params.Page-1) * params.Limit,
		Search: params.Search,
		Color:  ink.InkColor(params.Color),
	})

	docs := make([]Card, 0)
	for _, h := range hits {
		abilities := make([]Ability, 0)
		for _, a := range h.Abilities {
			abilities = append(abilities, Ability{
				Name:        a.Name,
				Description: a.Description,
			})
		}
		docs = append(docs, Card{
			ID:              h.ID,
			Name:            h.Name,
			Set:             h.Set,
			Version:         h.Version,
			Title:           h.Title,
			Cost:            h.Cost,
			Type:            h.Type,
			Number:          h.Number,
			Color:           h.Color,
			Illustrator:     h.Illustrator,
			Lore:            h.Lore,
			Strength:        h.Strength,
			Willpower:       h.Willpower,
			Movement:        h.Movement,
			Ink:             h.Ink,
			Characteristics: h.Characteristics,
			Abilities:       abilities,
			Variants:        h.Variants,
			Rarity:          h.Rarity,
			Language:        h.Language,
			Path:            h.Path,
			Franchise:       h.Franchise,
			Ordinal:         h.Ordinal,
			Formats:         h.Formats,
		})
	}

	return ListCardReponse{
		Page:  params.Page,
		Limit: params.Limit,
		Total: count,
		Docs:  docs,
		Size:  int64(len(docs)),
	}
}
