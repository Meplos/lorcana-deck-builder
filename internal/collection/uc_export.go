package collection

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"os"
	"sort"
	"strconv"

	"github.com/meplos/locana-deck-builder/internal/domain"
)

type ExportUseCase struct {
	repo Repository
}

func NewExportUC(repo Repository) *ExportUseCase {
	return &ExportUseCase{
		repo: repo,
	}
}

func (uc *ExportUseCase) ExportAll(ctx context.Context) (*os.File, error) {
	collections := uc.repo.List(ctx, ListFilter{
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

func buildRowsForCollection(collection domain.Collection) ([][]string, error) {
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
		colorsSerialized, err := json.Marshal(card.Color)
		if err != nil {
			return nil, err
		}
		ab, err := json.Marshal(card.Abilities)
		if err != nil {
			return nil, err
		}
		if card.ID == "" {
			return nil, errors.New("card id is empty")
		}
		rows = append(rows, []string{
			collection.Name,
			card.ID,
			card.Name,
			card.Title,
			card.Set,
			strconv.Itoa(card.Version),
			strconv.Itoa(card.Cost),
			card.Type,
			card.Number,
			string(colorsSerialized),
			card.Illustrator,
			strconv.Itoa(card.Lore),
			strconv.Itoa(card.Strength),
			strconv.Itoa(card.Willpower),
			strconv.Itoa(card.Movement),
			strconv.Itoa(card.Ink),
			card.Characteristics,
			string(ab),
			card.Variants,
			card.Rarity,
			card.Language,
			card.Path,
			card.Franchise,
			strconv.Itoa(card.Ordinal),
			card.Formats,
			strconv.Itoa(card.Quantity),
		})
	}
	return rows, nil
}
