package deck

import (
	"context"

	"github.com/meplos/locana-deck-builder/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	col *mongo.Collection
}

func NewRepository(DB *mongo.Database) Repository {
	return &MongoRepository{
		col: DB.Collection("decks"),
	}
}

func (r *MongoRepository) Save(ctx context.Context, d domain.Deck) error {
	if d.ID == "" {
		d.ID = primitive.NewObjectID().Hex()
	}
	_, err := r.col.InsertOne(ctx, d, options.InsertOne())
	return err
}

func (r *MongoRepository) CountBy(ctx context.Context) (int, error) {
	count, err := r.col.CountDocuments(ctx, bson.D{})
	return int(count), err
}

func (r *MongoRepository) FindBy(ctx context.Context) ([]domain.Deck, error) {
	deck := make([]domain.Deck, 0)
	cursor, err := r.col.Find(ctx, bson.D{})
	if err != nil {
		return deck, err
	}

	if err := cursor.All(ctx, &deck); err != nil {
		return deck, err
	}

	return deck, nil
}
