package auth

import (
	"context"

	"github.com/meplos/locana-deck-builder/internal/domain"
	"github.com/meplos/locana-deck-builder/internal/security/jwt"
	"github.com/meplos/locana-deck-builder/internal/user"
)

type (
	VerifyUseCase struct {
		repo user.Repository
		jwt  jwt.JWTManager
	}
	VerifyInput struct {
		Token string
	}
	VerifyOutput struct {
		User domain.User
	}
)

func NewVerifyUC(repo user.Repository, jwt jwt.JWTManager) *VerifyUseCase {
	return &VerifyUseCase{
		repo: repo,
		jwt:  jwt,
	}
}

func (uc *VerifyUseCase) Execute(ctx context.Context, input VerifyInput) (VerifyOutput, error) {
	ID, err := uc.jwt.Parse(input.Token)
	if err != nil {
		return VerifyOutput{}, err
	}

	u, err := uc.repo.FindByID(ctx, ID)
	if err != nil {
		return VerifyOutput{}, err
	}

	return VerifyOutput{
		User: u,
	}, nil
}
