package domain

import (
	"context"

	"github.com/meplos/locana-deck-builder/internal/collection/infra"
)

type CreateUseCase struct {
	ctx  context.Context
	repo infra.CollectionRepository
}

func NewCreateUC(ctx context.Context, collRepo infra.CollectionRepository) *CreateUseCase {
	return &CreateUseCase{
		ctx:  ctx,
		repo: collRepo,
	}
}

func (uc *CreateUseCase) Create(name string) error {
	_, err := uc.repo.FindOrCreate(uc.ctx, name)
	return err
}
