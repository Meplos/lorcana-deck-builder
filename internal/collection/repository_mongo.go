package collection

import (
	"context"
	"errors"
	"log"

	"github.com/meplos/locana-deck-builder/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCollectionRepository struct {
	col *mongo.Collection
}

func NewCollectionRepo(DB *mongo.Database) Repository {
	return &MongoCollectionRepository{
		col: DB.Collection("collections"),
	}
}

type ListFilter struct {
	Offset int
	Limit  int
	Name   string
}

type ownedCollectionsRaw struct {
	ID    interface{}                      `bson:"_id"`
	Name  string                           `bson:"name"`
	Cards map[string]domain.CollectionCard `bson:"cards"`
}

func (r *MongoCollectionRepository) List(ctx context.Context, f ListFilter) []domain.Collection {
	filter := bson.D{}
	if f.Name != "" {
		filter = append(filter, bson.E{
			Key:   "name",
			Value: f.Name,
		})
	}
	cursor, err := r.col.Find(ctx, filter, options.Find().SetSkip(int64(f.Offset)).SetLimit(int64(f.Limit)))
	if err != nil {
		return []domain.Collection{}
	}
	defer cursor.Close(ctx)

	owned := make([]domain.Collection, 0)
	for cursor.Next(ctx) {
		raw := domain.Collection{}
		if err := cursor.Decode(&raw); err != nil {
			continue
		}
		owned = append(owned, raw)
	}

	return owned
}

func (r *MongoCollectionRepository) Count(ctx context.Context, f ListFilter) int {
	count, err := r.col.CountDocuments(ctx, bson.D{
		{Key: "name", Value: f.Name},
	})
	if err != nil {
		return 0
	}

	return int(count)
}

func (r *MongoCollectionRepository) FindOrCreate(ctx context.Context, name string) (domain.Collection, error) {
	result := r.col.FindOne(ctx, bson.D{
		{Key: "name", Value: name},
	})

	parsed := domain.Collection{}
	var createErr error
	if err := result.Decode(&parsed); errors.Is(err, mongo.ErrNoDocuments) {
		_, createErr = r.col.InsertOne(ctx, domain.Collection{
			Name:  name,
			Cards: make(map[string]domain.CollectionCard),
		})
	}

	return parsed, createErr
}

func (r *MongoCollectionRepository) FindCollectionByName(ctx context.Context, name string) (domain.Collection, error) {
	res := r.col.FindOne(ctx, bson.M{
		"name": name,
	})

	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		log.Printf("Not found")
		return domain.Collection{}, res.Err()
	}

	raw := domain.Collection{}
	err := res.Decode(&raw)
	if err != nil {
		log.Printf("decode err : %s", err.Error())
	}
	return raw, err
}

func (r *MongoCollectionRepository) AddCardToCollection(ctx context.Context, collID primitive.ObjectID, card domain.CollectionCard) error {
	log.Printf("Col: %s, card: %s\n", collID.Hex(), card.ID)
	path := "cards." + card.ID
	res, err := r.col.UpdateOne(
		ctx,
		bson.M{"_id": collID},
		bson.M{"$set": bson.D{{Key: path, Value: card}}},
	)
	log.Printf("Matched: %d, Modified:%d, Upsert:%d", res.MatchedCount, res.ModifiedCount, res.UpsertedCount)
	return err
}

func (r *MongoCollectionRepository) FindByID(ctx context.Context, ID string) (domain.Card, error) {
	return domain.Card{}, nil
}
