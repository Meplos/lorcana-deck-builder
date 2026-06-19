// Package auth: uc_register.go handle register validation logic
package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/meplos/locana-deck-builder/internal/domain"
	"github.com/meplos/locana-deck-builder/internal/security/jwt"
	"github.com/meplos/locana-deck-builder/internal/security/password"
	"github.com/meplos/locana-deck-builder/internal/user"
)

type (
	RegisterUseCase struct {
		repo   user.Repository
		hasher password.Hasher
		jwt    jwt.JWTManager
	}
	RegisterInput struct {
		Name            string
		Email           string
		Password        string
		ConfirmPassword string
	}
	UserCreateRequest struct {
		Name  string
		Email string
		Hash  []byte
	}
	RegisterOutput struct {
		Token string
	}
	UserID string
)

func (uc *RegisterUseCase) Execute(ctx context.Context, input RegisterInput) (RegisterOutput, error) {
	if err := validateEmail(input.Email); err != nil {
		return RegisterOutput{}, err
	}

	if err := validatePassword(input.Password, input.ConfirmPassword); err != nil {
		return RegisterOutput{}, err
	}

	if exists, err := uc.repo.ExistsByEmail(ctx, input.Email); err != nil {
		return RegisterOutput{}, err
	} else if exists {
		return RegisterOutput{}, ErrAlreadyExist
	}

	hash, err := uc.hasher.Hash([]byte(input.Password))
	if err != nil {
		return RegisterOutput{}, err
	}

	user := domain.User{
		ID:        uuid.NewString(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      input.Name,
		Email:     input.Email,
		Hash:      hash,
		LastLogin: time.Now(),
	}

	if err := uc.repo.Create(ctx, user); err != nil {
		return RegisterOutput{}, err
	}

	token, err := uc.jwt.Create(user.ID)
	if err != nil {
		return RegisterOutput{}, err
	}

	return RegisterOutput{
		Token: token,
	}, nil
}

func NewRegisterUC(repo user.Repository, hasher password.Hasher, jwt jwt.JWTManager) *RegisterUseCase {
	return &RegisterUseCase{
		repo:   repo,
		jwt:    jwt,
		hasher: hasher,
	}
}
