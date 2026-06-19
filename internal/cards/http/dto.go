package http

import (
	"fmt"

	"github.com/meplos/locana-deck-builder/internal/domain"
)

type CardDTO struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Title     string            `json:"title"`
	Colors    []domain.InkColor `json:"colors"`
	Number    string            `json:"number"`
	Set       string            `json:"set"`
	Rarity    string            `json:"rarity"`
	Franchise string            `json:"franchise"`
	FilePath  string            `json:"filepath"`
}

type PaginateRequest struct {
	PageNumber uint   `query:"page"`
	Limit      uint   `query:"limit"`
	Search     string `query:"search"`
	Color      string `query:"color"`
}

type PaginateResponse struct {
	Page  uint      `json:"page"`
	Docs  []CardDTO `json:"docs"`
	Total uint      `json:"total"`
	Size  uint      `json:"size"`
}

func BuildCardDTO(c domain.Card) CardDTO {
	return CardDTO{
		ID:        c.ID,
		Name:      c.Name,
		Title:     c.Title,
		Colors:    c.Color,
		Number:    c.Number,
		Set:       c.Set,
		Rarity:    c.Rarity,
		Franchise: c.Franchise,
		FilePath:  fmt.Sprintf("https://cdn.dreamborn.ink/images/fr/cards/%s", c.ID),
	}
}
