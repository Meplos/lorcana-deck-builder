package deck

import (
	"context"
	"fmt"
	"strings"
)

type (
	ExportUseCase struct {
		repo Repository
	}
	ExportInput struct {
		ID string
	}
	ExportOutput struct {
		Content string
	}
)

func NewExportUC(repo Repository) *ExportUseCase {
	return &ExportUseCase{
		repo: repo,
	}
}

func (uc *ExportUseCase) Execute(ctx context.Context, input ExportInput) (ExportOutput, error) {
	d, err := uc.repo.FindByID(ctx, input.ID)
	if err != nil {
		return ExportOutput{}, err
	}

	lines := make([]string, 0)

	for _, c := range d.Deck {
		if c.Type == "character" {
			lines = append(lines, fmt.Sprintf("%d %s - %s", c.Quantity, c.Name, c.Title))
		} else {
			lines = append(lines, fmt.Sprintf("%d %s", c.Quantity, c.Name))
		}
	}
	content := strings.Join(lines, "\n")

	return ExportOutput{
		Content: content,
	}, nil
}
