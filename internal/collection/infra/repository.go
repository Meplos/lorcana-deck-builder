package infra

import (
	"context"
	"errors"
	"log"

	"github.com/meplos/locana-deck-builder/internal/db/mongo/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CollectionRepository interface {
	List(ctx context.Context, f ListFilter) []OwnedCollections
	Count(ctx context.Context, f ListFilter) int
	FindOrCreate(ctx context.Context, name string) (OwnedCollections, error)
	FindCollectionByName(ctx context.Context, name string) (OwnedCollections, error)
	AddCardToCollection(ctx context.Context, collID primitive.ObjectID, card schema.OwnedCard) error
}

type MongoCollectionRepository struct {
	col *mongo.Collection
}

func NewCollectionRepo(DB *mongo.Database) CollectionRepository {
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
	ID    interface{}          `bson:"_id"`
	Name  string               `bson:"name"`
	Cards map[string]OwnedCard `bson:"cards"`
}

func (r *MongoCollectionRepository) List(ctx context.Context, f ListFilter) []OwnedCollections {
	filter := bson.D{}
	if f.Name != "" {
		filter = append(filter, bson.E{
			Key:   "name",
			Value: f.Name,
		})
	}
	cursor, err := r.col.Find(ctx, filter, options.Find().SetSkip(int64(f.Offset)).SetLimit(int64(f.Limit)))
	if err != nil {
		return []OwnedCollections{}
	}
	defer cursor.Close(ctx)

	owned := make([]OwnedCollections, 0)
	for cursor.Next(ctx) {
		raw := ownedCollectionsRaw{}
		if err := cursor.Decode(&raw); err != nil {
			continue
		}
		normalized, err := r.normalizeCollectionID(ctx, raw)
		if err != nil {
			continue
		}
		owned = append(owned, normalized)
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

func (r *MongoCollectionRepository) FindOrCreate(ctx context.Context, name string) (OwnedCollections, error) {
	result := r.col.FindOne(ctx, bson.D{
		{Key: "name", Value: name},
	})

	parsed := OwnedCollections{}
	var createErr error
	if err := result.Decode(&parsed); errors.Is(err, mongo.ErrNoDocuments) {
		_, createErr = r.col.InsertOne(ctx, OwnedCollections{
			Name:  name,
			Cards: make(map[string]OwnedCard),
		})
	}

	return parsed, createErr
}

func (r *MongoCollectionRepository) FindCollectionByName(ctx context.Context, name string) (OwnedCollections, error) {
	res := r.col.FindOne(ctx, bson.M{
		"name": name,
	})

	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		log.Printf("Not found")
		return OwnedCollections{}, res.Err()
	}

	raw := ownedCollectionsRaw{}
	err := res.Decode(&raw)
	if err != nil {
		log.Printf("decode err : %s", err.Error())
	}
	data, normalizeErr := r.normalizeCollectionID(ctx, raw)
	if normalizeErr != nil {
		return OwnedCollections{}, normalizeErr
	}
	log.Printf("data ID: %s", data.ID.Hex())
	return data, err
}

func (r *MongoCollectionRepository) AddCardToCollection(ctx context.Context, collID primitive.ObjectID, card schema.OwnedCard) error {
	log.Printf("Col: %s, card: %s\n", collID.Hex(), card.ID)
	path := "cards." + card.ID
	res, err := r.col.UpdateOne(
		ctx,
		bson.M{"_id": collID},
		bson.M{"$set": bson.M{path: card}},
	)
	log.Printf("Matched: %d, Modified:%d, Upsert:%d", res.MatchedCount, res.ModifiedCount, res.UpsertedCount)
	return err
}

func (r *MongoCollectionRepository) normalizeCollectionID(ctx context.Context, raw ownedCollectionsRaw) (OwnedCollections, error) {
	switch v := raw.ID.(type) {
	case primitive.ObjectID:
		return OwnedCollections{ID: v, Name: raw.Name, Cards: raw.Cards}, nil
	case string:
		if v != "" {
			if oid, err := primitive.ObjectIDFromHex(v); err == nil {
				return OwnedCollections{ID: oid, Name: raw.Name, Cards: raw.Cards}, nil
			}
		}
	}

	newID := primitive.NewObjectID()
	_, err := r.col.InsertOne(ctx, OwnedCollections{
		ID:    newID,
		Name:  raw.Name,
		Cards: raw.Cards,
	})
	if err != nil {
		return OwnedCollections{}, err
	}
	if raw.ID != nil {
		_, _ = r.col.DeleteOne(ctx, bson.M{"_id": raw.ID})
	}

	return OwnedCollections{ID: newID, Name: raw.Name, Cards: raw.Cards}, nil
}

type CardRepository interface {
	FindById(ctx context.Context, ID string) (schema.Card, error)
}

type MongoCardRepository struct {
	col *mongo.Collection
}

func NewCardRepo(DB *mongo.Database) CardRepository {
	return &MongoCardRepository{
		col: DB.Collection("cards"),
	}
}

func (r *MongoCardRepository) FindById(ctx context.Context, ID string) (schema.Card, error) {
	result := r.col.FindOne(ctx, bson.D{{Key: "_id", Value: ID}})
	if result.Err() != nil {
		return schema.Card{}, result.Err()
	}
	card := schema.Card{}
	err := result.Decode(&card)
	return card, err
}
