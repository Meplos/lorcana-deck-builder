package collection

import (
	"context"
)

type CreateUseCase struct {
	repo Repository
}

func NewCreateUC(repo Repository) *CreateUseCase {
	return &CreateUseCase{
		repo: repo,
	}
}

func (uc *CreateUseCase) Create(ctx context.Context, name string) error {
	_, err := uc.repo.FindOrCreate(ctx, name)
	return err
}
