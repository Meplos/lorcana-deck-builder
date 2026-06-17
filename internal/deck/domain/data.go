package domain

import "github.com/meplos/locana-deck-builder/internal/db/mongo/schema"

type OwnedCard struct {
	Card     schema.Card `bson:"card"`
	FilePath string      `bson:"filepath"`
	Quantity int         `bson:"quantity"`
}
