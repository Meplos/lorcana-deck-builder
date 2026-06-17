package domain

import "github.com/meplos/locana-deck-builder/internal/ink"

type Ability struct {
	Name        string
	Description string
}

type Card struct {
	ID              string
	Name            string
	Set             string
	Version         int
	Title           string
	Cost            int
	Type            string
	Number          string
	Color           []ink.InkColor
	Illustrator     string
	Lore            int
	Strength        int
	Willpower       int
	Movement        int
	Ink             int
	Characteristics string
	Abilities       []Ability
	Variants        string
	Rarity          string
	Language        string
	Path            string
	Franchise       string
	Ordinal         int
	Formats         string
}
