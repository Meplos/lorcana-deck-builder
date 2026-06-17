package infra

import (
	"context"

	"github.com/meplos/locana-deck-builder/internal/db/mongo/schema"
	"github.com/meplos/locana-deck-builder/internal/ink"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ListFilter struct {
	Offset int64
	Limit  int64
	Search string
	Color  ink.InkColor
}

type CardRepository interface {
	ListCards(ctx context.Context, f ListFilter) []schema.Card
	CountCards(ctx context.Context, f ListFilter) int64
}

type MongoCardRepository struct {
	col *mongo.Collection
}

func SetupMongoCardRepository(db *mongo.Database) *MongoCardRepository {
	return &MongoCardRepository{
		col: db.Collection("cards"),
	}
}

func (r *MongoCardRepository) ListCards(ctx context.Context, f ListFilter) []schema.Card {
	filter := buildMongoFilter(f)
	cursor, err := r.col.Find(ctx, filter, options.Find().SetLimit(int64(f.Limit)).SetSkip(f.Offset))
	if err != nil {
		return make([]schema.Card, 0)
	}
	defer cursor.Close(ctx)

	var docs []schema.Card
	if err := cursor.All(ctx, &docs); err != nil {
		return make([]schema.Card, 0)
	}

	return docs
}

func (r *MongoCardRepository) CountCards(ctx context.Context, f ListFilter) int64 {
	filter := buildMongoFilter(f)
	size, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return 0
	}
	return size
}

func buildMongoFilter(f ListFilter) bson.D {
	filter := bson.D{}
	if f.Search != "" {
		filter = append(filter, bson.E{
			Key: "name", Value: bson.D{
				{Key: "$regex", Value: "^" + f.Search},
				{Key: "$options", Value: "i"},
			},
		})
	}

	if f.Color != "" {
		filter = append(filter, bson.E{Key: "color", Value: f.Color})
	}
	return filter
}
