package schema

import "github.com/meplos/locana-deck-builder/internal/ink"

type Ability struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Card struct {
	ID              string         `bson:"_id"`
	Name            string         `bson:"name"`
	Set             string         `bson:"set"`
	Version         int            `bson:"ver"`
	Title           string         `bson:"title"`
	Cost            int            `bson:"cost"`
	Type            string         `bson:"type"`
	Number          string         `bson:"number"`
	Color           []ink.InkColor `bson:"color"`
	Illustrator     string         `bson:"illustrator"`
	Lore            int            `bson:"lore"`
	Strength        int            `bson:"strength"`
	Willpower       int            `bson:"willpower"`
	Movement        int            `bson:"movement"`
	Ink             int            `bson:"ink"`
	Characteristics string         `bson:"characteristics"`
	Abilities       []Ability      `bson:"abilities"`
	Variants        string         `bson:"variants"`
	Rarity          string         `bson:"rarity"`
	Language        string         `bson:"language"`
	Path            string         `bson:"path"`
	Franchise       string         `bson:"franchise"`
	Ordinal         int            `bson:"ordinal"`
	Formats         string         `bson:"formats"`
}
