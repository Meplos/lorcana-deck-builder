package cards

import (
	"context"
	"log"

	"github.com/meplos/locana-deck-builder/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCardRepository struct {
	col *mongo.Collection
}

func SetupMongoCardRepository(db *mongo.Database) *MongoCardRepository {
	return &MongoCardRepository{
		col: db.Collection("cards"),
	}
}

func (r *MongoCardRepository) ListCards(ctx context.Context, f ListFilter) []domain.Card {
	filter := buildMongoFilter(f)
	cursor, err := r.col.Find(ctx, filter, options.Find().SetLimit(int64(f.Limit)).SetSkip(f.Offset))
	if err != nil {
		return make([]domain.Card, 0)
	}
	defer cursor.Close(ctx)

	var docs []domain.Card
	if err := cursor.All(ctx, &docs); err != nil {
		return make([]domain.Card, 0)
	}

	return docs
}

func (r *MongoCardRepository) CountCards(ctx context.Context, f ListFilter) int64 {
	filter := buildMongoFilter(f)
	log.Print(filter)
	size, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		log.Print(err)
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

func (r *MongoCardRepository) FindByID(ctx context.Context, ID string) (domain.Card, error) {
	card := new(domain.Card)
	err := r.col.FindOne(ctx, bson.D{{Key: "_id", Value: ID}}).Decode(card)

	return *card, err
}

func (r *MongoCardRepository) FindByIDs(ctx context.Context, IDs []string) (map[string]domain.Card, error) {
	cursor, err := r.col.Find(ctx, bson.D{{
		Key: "_id",
		Value: bson.D{
			{
				Key:   "$in",
				Value: IDs,
			},
		},
	}})
	if err != nil {
		return nil, err
	}

	tmpArr := make([]domain.Card, 0)
	if err := cursor.All(ctx, &tmpArr); err != nil {
		return nil, err
	}

	result := make(map[string]domain.Card)
	for _, c := range tmpArr {
		result[c.ID] = c
	}

	return result, nil
}
