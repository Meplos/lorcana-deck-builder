// Package auth: uc_register.go handle register validation logic

package auth

import (
	"context"
	"time"

	"github.com/meplos/locana-deck-builder/internal/security/jwt"
	"github.com/meplos/locana-deck-builder/internal/security/password"
	"github.com/meplos/locana-deck-builder/internal/user"
)

type (
	LoginUseCase struct {
		repo user.Repository
		hash password.Hasher
		jwt  jwt.JWTManager
	}
	LoginInput struct {
		Email    string
		Password string
	}
	LoginOutput struct {
		Token string
	}
)

func NewLoginUC(repo user.Repository, hash password.Hasher, jwt jwt.JWTManager) *LoginUseCase {
	return &LoginUseCase{
		repo: repo,
		hash: hash,
		jwt:  jwt,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (LoginOutput, error) {
	u, err := uc.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		return LoginOutput{}, ErrNotExist
	}

	if !uc.hash.Compare(u.Hash, []byte(input.Password)) {
		return LoginOutput{}, ErrWrongPassword
	}

	if err := uc.repo.LogUserAt(ctx, u.ID, time.Now()); err != nil {
		return LoginOutput{}, err
	}

	token, err := uc.jwt.Create(u.ID)
	if err != nil {
		return LoginOutput{}, err
	}

	return LoginOutput{
		Token: token,
	}, nil
}
