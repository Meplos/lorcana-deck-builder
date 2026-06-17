package collection

import (
	"fmt"

	"github.com/meplos/locana-deck-builder/internal/collection/infra"
	"github.com/meplos/locana-deck-builder/internal/ink"
)

type OwnedCollectionDTO struct {
	ID    string                    `json:"id"`
	Name  string                    `json:"name"`
	Cards map[string]OwnedCardDTO   `json:"cards"`
}
type OwnedCardDTO struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Title    string         `json:"title"`
	Colors   []ink.InkColor `json:"colors"`
	Number   string         `json:"number"`
	Set      string         `json:"set"`
	Rarity   string         `json:"rarity"`
	FilePath string         `json:"filepath"`
	Quantity int            `json:"quantity"`
}

func BuildOwnedCardDTO(c infra.OwnedCard) OwnedCardDTO {
	return OwnedCardDTO{
		ID:       c.Card.ID,
		Name:     c.Card.Name,
		Title:    c.Card.Title,
		Colors:   c.Card.Color,
		Number:   c.Card.Number,
		Set:      c.Card.Set,
		Rarity:   c.Card.Rarity,
		FilePath: fmt.Sprintf("https://cdn.dreamborn.ink/images/fr/cards/%s", c.Card.ID),
		Quantity: c.Quantity,
	}
}

func BuildOwnedCollectionDTO(col infra.OwnedCollections) OwnedCollectionDTO {
	cards := make(map[string]OwnedCardDTO, len(col.Cards))
	for id, c := range col.Cards {
		cards[id] = BuildOwnedCardDTO(c)
	}

	return OwnedCollectionDTO{
		ID:    col.ID.Hex(),
		Name:  col.Name,
		Cards: cards,
	}
}
