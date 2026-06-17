package http

import (
	"fmt"

	"github.com/meplos/locana-deck-builder/internal/domain"
)

type OwnedCollectionDTO struct {
	ID    string                  `json:"id"`
	Name  string                  `json:"name"`
	Cards map[string]OwnedCardDTO `json:"cards"`
}
type OwnedCardDTO struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Title    string            `json:"title"`
	Colors   []domain.InkColor `json:"colors"`
	Number   string            `json:"number"`
	Set      string            `json:"set"`
	Rarity   string            `json:"rarity"`
	FilePath string            `json:"filepath"`
	Quantity int               `json:"quantity"`
}

func BuildOwnedCardDTO(c domain.CollectionCard) OwnedCardDTO {
	return OwnedCardDTO{
		ID:       c.ID,
		Name:     c.Name,
		Title:    c.Title,
		Colors:   c.Color,
		Number:   c.Number,
		Set:      c.Set,
		Rarity:   c.Rarity,
		FilePath: fmt.Sprintf("https://cdn.dreamborn.ink/images/fr/cards/%s", c.ID),
		Quantity: c.Quantity,
	}
}

func BuildOwnedCollectionDTO(col domain.Collection) OwnedCollectionDTO {
	cards := make(map[string]OwnedCardDTO, len(col.Cards))
	for id, c := range col.Cards {
		cards[id] = BuildOwnedCardDTO(c)
	}

	return OwnedCollectionDTO{
		ID:    col.ID,
		Name:  col.Name,
		Cards: cards,
	}
}
