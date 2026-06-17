package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	_ "github.com/glebarez/go-sqlite"
	localMongo "github.com/meplos/locana-deck-builder/internal/database/mongo"
	database "github.com/meplos/locana-deck-builder/internal/database/sql"
	"github.com/meplos/locana-deck-builder/internal/domain"
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

	queries := database.New(conn)

	sqlCards, err := queries.ListCards(ctx)
	if err != nil {
		log.Fatal(err)
	}

	collections := make([]any, 0)
	for _, c := range sqlCards {
		ab := []domain.Ability{}

		if c.Abilities.Valid {
			if err := json.Unmarshal([]byte(c.Abilities.String), &ab); err != nil {
				log.Fatal(err)
			}
		}
		log.Printf("%v", ab)

		nC := domain.Card{
			ID:              c.ID,
			Name:            c.Name,
			Set:             c.SetId,
			Version:         int(c.VersionId),
			Title:           c.Title.String,
			Cost:            int(c.Cost),
			Type:            c.Type,
			Number:          c.Number,
			Color:           domain.GetInkStrings(domain.InkMask(c.ColorMask)),
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

	db, err := localMongo.Connect(ctx, "lorcana")
	if err != nil {
		panic(err)
	}

	coll := db.Collection("cards")
	log.Println(coll.Name())

	res, err := coll.InsertMany(ctx, collections, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(res)
}
