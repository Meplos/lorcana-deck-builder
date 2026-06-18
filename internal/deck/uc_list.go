// Package deck: uc_list.go -> List All decks UseCase
package deck

import (
	"context"

	"github.com/meplos/locana-deck-builder/assets/images"
	"github.com/meplos/locana-deck-builder/internal/domain"
)

type ListUseCase struct {
	repo       Repository
	uriBuilder images.ImageURIBuilder
}

func NewListUC(repo Repository, uriBuilder images.ImageURIBuilder) *ListUseCase {
	return &ListUseCase{
		repo:       repo,
		uriBuilder: uriBuilder,
	}
}

type (
	CardItem struct {
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
	DeckItem struct {
		Strategy string
		Name     string
		Size     int
		Deck     []CardItem
	}
	ListOutput struct {
		Total int
		Docs  []DeckItem
	}
)

func (uc *ListUseCase) Execute(ctx context.Context) ListOutput {
	count, err := uc.repo.CountBy(ctx)
	if err != nil {
		return ListOutput{
			Total: 0,
			Docs:  []DeckItem{},
		}
	}
	decks, err := uc.repo.FindBy(ctx)
	if err != nil {
		return ListOutput{
			Total: 0,
			Docs:  []DeckItem{},
		}
	}

	docs := make([]DeckItem, 0)
	for _, d := range decks {
		docs = append(docs, DeckItem{
			Strategy: d.Strategy,
			Name:     d.Name,
			Size:     d.Size,
			Deck:     uc.mapCards(d.Deck),
		})
	}

	return ListOutput{
		Total: count,
		Docs:  docs,
	}
}

func (uc *ListUseCase) mapCards(cards []domain.DeckCard) []CardItem {
	result := make([]CardItem, 0)
	for _, c := range cards {
		result = append(result, CardItem{
			ID:       c.ID,
			Name:     c.Name,
			Title:    c.Title,
			Number:   c.Number,
			Set:      c.Set,
			Type:     c.Type,
			Rarity:   c.Rarity,
			Quantity: c.Quantity,
			Filepath: uc.uriBuilder.Card(c.ID),
		})
	}
	return result
}
