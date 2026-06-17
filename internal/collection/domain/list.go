package domain

import (
	"context"

	"github.com/meplos/locana-deck-builder/internal/collection/infra"
)

type ListUseCase struct {
	ctx  context.Context
	repo infra.CollectionRepository
}

func NewListUC(ctx context.Context, repo infra.CollectionRepository) *ListUseCase {
	return &ListUseCase{
		ctx:  ctx,
		repo: repo,
	}
}

type ListResponse struct {
	Page  int
	Docs  []infra.OwnedCollections
	Total int
	Size  int
}

type PaginateParams struct {
	Page  int
	Limit int
	Name  string
}

func (uc *ListUseCase) List(p PaginateParams) ListResponse {
	count := uc.repo.Count(uc.ctx, infra.ListFilter{
		Name: p.Name,
	})
	list := uc.repo.List(uc.ctx, infra.ListFilter{
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
