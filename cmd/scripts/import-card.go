package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	_ "github.com/glebarez/go-sqlite"
	"github.com/meplos/locana-deck-builder/internal/db/mongo/schema"
	db "github.com/meplos/locana-deck-builder/internal/db/sql"
	"github.com/meplos/locana-deck-builder/internal/ink"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	home, _ := os.UserHomeDir()
	filename := filepath.Join(home, "Documents", "perso", "lorcana-deck-builder", "lorcana.card.db")

	conn, err := sql.Open("sqlite", filename)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	queries := db.New(conn)

	sqlCards, err := queries.ListCards(ctx)
	if err != nil {
		log.Fatal(err)
	}

	collections := make([]any, 0)
	for _, c := range sqlCards {
		ab := []schema.Ability{}

		if c.Abilities.Valid {
			if err := json.Unmarshal([]byte(c.Abilities.String), &ab); err != nil {
				log.Fatal(err)
			}
		}
		log.Printf("%v", ab)

		nC := schema.Card{
			ID:              c.ID,
			Name:            c.Name,
			Set:             c.SetId,
			Version:         int(c.VersionId),
			Title:           c.Title.String,
			Cost:            int(c.Cost),
			Type:            c.Type,
			Number:          c.Number,
			Color:           ink.GetInkStrings(ink.InkMask(c.ColorMask)),
			Illustrator:     c.Illustrator.String,
			Lore:            int(c.Lore.Int64),
			Strength:        int(c.Strength.Int64),
			Willpower:       int(c.Willpower.Int64),
			Movement:        int(c.Movement.Int64),
			Ink:             int(c.Ink),
			Characteristics: c.Characteristics,
			Abilities:       ab,
			Variants:        c.Variants,
			Rarity:          c.Rarity.String,
			Language:        c.Language,
			Path:            c.Path.String,
			Franchise:       c.Franchise.String,
			Ordinal:         int(c.Ordinal),
			Formats:         c.Formats.String,
		}
		collections = append(collections, nC)

	}

	log.Printf("Size: %d\n", len(collections))

	mongoURI := "mongodb://root:example@localhost:9999"
	log.Printf("[IMPORT_CARD] MONGO URI: %v", mongoURI)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	coll := client.Database("lorcana").Collection("cards")
	log.Println(coll.Name())

	res, err := coll.InsertMany(ctx, collections, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(res)
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()
}
