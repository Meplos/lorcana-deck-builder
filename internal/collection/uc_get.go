package collection

import (
	"context"

	"github.com/meplos/locana-deck-builder/internal/domain"
)

type ListUseCase struct {
	repo Repository
}

func NewListUC(repo Repository) *ListUseCase {
	return &ListUseCase{
		repo: repo,
	}
}

type ListResponse struct {
	Page  int
	Docs  []domain.Collection
	Total int
	Size  int
}

type PaginateParams struct {
	Page  int
	Limit int
	Name  string
}

func (uc *ListUseCase) List(ctx context.Context, p PaginateParams) ListResponse {
	count := uc.repo.Count(ctx, ListFilter{
		Name: p.Name,
	})
	list := uc.repo.List(ctx, ListFilter{
		Offset: (p.Page - 1) * p.Limit,
		Limit:  p.Limit,
		Name:   p.Name,
	})

	return ListResponse{
		Page:  p.Page,
		Total: count,
		Size:  len(list),
		Docs:  list,
	}
}
