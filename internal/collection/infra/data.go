package infra

import (
	"github.com/meplos/locana-deck-builder/internal/db/mongo/schema"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OwnedCard struct {
	Card     schema.Card `bson:"card"`
	Quantity int         `bson:"quantity"`
}
type OwnedCollections struct {
	ID    primitive.ObjectID   `bson:"_id,omitempty"`
	Name  string               `bson:"name"`
	Cards map[string]OwnedCard `bson:"cards"`
}
