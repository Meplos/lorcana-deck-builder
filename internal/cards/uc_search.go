package cards

import (
	"context"

	"github.com/meplos/locana-deck-builder/internal/domain"
)

type ListCardUC struct {
	cardRepository Repository
}

func CreateListCardUC(repo Repository) *ListCardUC {
	return &ListCardUC{
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
	Docs  []domain.Card
	Total int64
	Size  int64
}

func (uc *ListCardUC) ListCards(ctx context.Context, params ListCardParams) ListCardReponse {
	count := uc.cardRepository.CountCards(ctx, ListFilter{
		Limit:  params.Limit,
		Offset: int64(params.Page-1) * params.Limit,
		Search: params.Search,
		Color:  domain.InkColor(params.Color),
	})

	hits := uc.cardRepository.ListCards(ctx, ListFilter{
		Limit:  params.Limit,
		Offset: int64(params.Page-1) * params.Limit,
		Search: params.Search,
		Color:  domain.InkColor(params.Color),
	})

	return ListCardReponse{
		Page:  params.Page,
		Limit: params.Limit,
		Total: count,
		Docs:  hits,
		Size:  int64(len(hits)),
	}
}
