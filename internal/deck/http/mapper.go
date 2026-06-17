package http

import "github.com/meplos/locana-deck-builder/internal/deck"

func MapDeckCards(cards []deck.DeckCard) []DeckBuildCard {
	result := make([]DeckBuildCard, 0)

	for _, c := range cards {
		result = append(result, MapCard(c))
	}
	return result
}

func MapCard(card deck.DeckCard) DeckBuildCard {
	return DeckBuildCard{
		ID:       card.ID,
		Name:     card.Name,
		Title:    card.Title,
		Number:   card.Number,
		Set:      card.Set,
		Type:     card.Type,
		Rarity:   card.Rarity,
		Quantity: card.Quantity,
		Filepath: card.Filepath,
	}
}
