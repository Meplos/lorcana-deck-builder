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

func MapCardItems(cards []deck.CardItem) []DeckBuildCard {
	result := make([]DeckBuildCard, 0)

	for _, c := range cards {
		result = append(result, DeckBuildCard(c))
	}
	return result
}

func MapDeck(decks []deck.DeckItem) []DeckBuildResponse {
	result := make([]DeckBuildResponse, 0)
	for _, d := range decks {
		result = append(result, DeckBuildResponse{
			Size:     d.Size,
			Name:     d.Name,
			Strategy: d.Strategy,
			Deck:     MapCardItems(d.Deck),
		})
	}
	return result
}
