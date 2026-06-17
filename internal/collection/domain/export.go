package domain

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"os"
	"sort"
	"strconv"

	"github.com/meplos/locana-deck-builder/internal/collection/infra"
)

type ExportUseCase struct {
	ctx  context.Context
	repo infra.CollectionRepository
}

func NewExportUC(ctx context.Context, repo infra.CollectionRepository) *ExportUseCase {
	return &ExportUseCase{
		ctx:  ctx,
		repo: repo,
	}
}

func (uc *ExportUseCase) ExportAll() (*os.File, error) {
	collections := uc.repo.List(uc.ctx, infra.ListFilter{
		Offset: 0,
		Limit:  500,
	})
	headers := []string{
		"collection",
		"card_id",
		"Name",
		"Title",
		"Set",
		"Version",
		"Cost",
		"Type",
		"Number",
		"Color",
		"Illustrator",
		"Lore",
		"Strength",
		"Willpower",
		"Movement",
		"Ink",
		"Characteristics",
		"Abilities",
		"Variants",
		"Rarity",
		"Language",
		"Path",
		"Franchise",
		"Ordinal",
		"Formats",
		"Quantity",
	}
	f, err := os.CreateTemp(os.TempDir(), "export.csv")
	if err != nil {
		return nil, err
	}

	writer := csv.NewWriter(f)

	if err := writer.Write(headers); err != nil {
		_ = f.Close()
		return nil, err
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		_ = f.Close()
		return nil, err
	}

	for _, c := range collections {
		rows, err := buildRowsForCollection(c)
		if err != nil {
			_ = f.Close()
			return nil, err
		}
		if err := writer.WriteAll(rows); err != nil {
			_ = f.Close()
			return nil, err
		}
		writer.Flush()
		if err := writer.Error(); err != nil {
			_ = f.Close()
			return nil, err
		}
	}

	if _, err := f.Seek(0, 0); err != nil {
		_ = f.Close()
		return nil, err
	}

	return f, nil
}

func buildRowsForCollection(collection infra.OwnedCollections) ([][]string, error) {
	rows := make([][]string, 0)
	keys := make([]string, 0, len(collection.Cards))
	for id := range collection.Cards {
		keys = append(keys, id)
	}
	sort.Strings(keys)
	for _, id := range keys {
		card, ok := collection.Cards[id]
		if !ok {
			continue
		}
		colorsSerialized, err := json.Marshal(card.Card.Color)
		if err != nil {
			return nil, err
		}
		ab, err := json.Marshal(card.Card.Abilities)
		if err != nil {
			return nil, err
		}
		if card.Card.ID == "" {
			return nil, errors.New("card id is empty")
		}
		rows = append(rows, []string{
			collection.Name,
			card.Card.ID,
			card.Card.Name,
			card.Card.Title,
			card.Card.Set,
			strconv.Itoa(card.Card.Version),
			strconv.Itoa(card.Card.Cost),
			card.Card.Type,
			card.Card.Number,
			string(colorsSerialized),
			card.Card.Illustrator,
			strconv.Itoa(card.Card.Lore),
			strconv.Itoa(card.Card.Strength),
			strconv.Itoa(card.Card.Willpower),
			strconv.Itoa(card.Card.Movement),
			strconv.Itoa(card.Card.Ink),
			card.Card.Characteristics,
			string(ab),
			card.Card.Variants,
			card.Card.Rarity,
			card.Card.Language,
			card.Card.Path,
			card.Card.Franchise,
			strconv.Itoa(card.Card.Ordinal),
			card.Card.Formats,
			strconv.Itoa(card.Quantity),
		})
	}
	return rows, nil
}
