package schema

import "go.mongodb.org/mongo-driver/bson/primitive"

type OwnedCard struct {
	Card
	Quantity int `bson:"quantity"`
}
type OwnedCollections struct {
	ID    primitive.ObjectID   `bson:"_id,omitempty"`
	Name  string               `bson:"name"`
	Cards map[string]OwnedCard `bson:"cards"`
}
