package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	localMongo "github.com/meplos/locana-deck-builder/internal/database/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	_ = godotenv.Load()

	ctx := context.Background()

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "lorcana"
	}

	db, err := localMongo.Connect(ctx, dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer localMongo.Disconnect(ctx)

	collections := db.Collection("collections")
	cursor, err := collections.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	updatedCollections := 0
	updatedCards := 0

	for cursor.Next(ctx) {
		doc := bson.M{}
		if err := cursor.Decode(&doc); err != nil {
			log.Printf("decode error: %v", err)
			continue
		}

		cardsRaw, ok := doc["cards"]
		if !ok {
			continue
		}

		cardsMap, ok := toBsonM(cardsRaw)
		if !ok {
			continue
		}

		updates := bson.M{}
		for cardID, cardValue := range cardsMap {
			cardMap, ok := toBsonM(cardValue)
			if !ok {
				continue
			}

			cardDocRaw, hasCard := cardMap["card"]
			if !hasCard {
				continue
			}

			cardDoc, ok := toBsonM(cardDocRaw)
			if !ok {
				continue
			}
			delete(cardDoc, "card")

			quantity, hasQuantity := cardMap["quantity"]
			if !hasQuantity {
				if q, ok := cardMap["Quantity"]; ok {
					quantity = q
					hasQuantity = true
				}
			}

			if hasQuantity {
				cardDoc["quantity"] = quantity
			}

			updates["cards."+cardID] = cardDoc
			updatedCards++
		}

		if len(updates) == 0 {
			continue
		}

		_, err := collections.UpdateOne(ctx, bson.M{"_id": doc["_id"]}, bson.M{"$set": updates})
		if err != nil {
			log.Printf("update error for %v: %v", doc["_id"], err)
			continue
		}
		updatedCollections++
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("done: updated collections=%d, updated cards=%d", updatedCollections, updatedCards)
}

func toBsonM(value any) (bson.M, bool) {
	switch typed := value.(type) {
	case bson.M:
		return typed, true
	case map[string]any:
		return bson.M(typed), true
	case bson.D:
		return typed.Map(), true
	default:
		return nil, false
	}
}
