package infra

import "github.com/meplos/locana-deck-builder/internal/ink"

type AICard struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type CardSummary struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	Title           string         `json:"title"`
	Cost            int            `json:"cost"`
	Type            string         `json:"type"`
	Color           []ink.InkColor `json:"colors"`
	Lore            int            `json:"lore"`
	Strength        int            `json:"stregth"`
	Willpower       int            `json:"willpower"`
	Movement        int            `json:"movement"`
	Ink             int            `json:"ink"`
	Characteristics string         `json:"characteristics"`
	Abilities       []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"abilities"`
	Quantity int `json:"quantity"`
}

type DeckRequest struct {
	Cards []CardSummary `json:"cards"`
	Level string        `json:"level"`
}

type DeckResponse struct {
	Name     string   `json:"name"`
	Deck     []AICard `json:"deck"`
	Strategy string   `json:"strategie"`
	Size     int      `json:"nbCard"`
}
