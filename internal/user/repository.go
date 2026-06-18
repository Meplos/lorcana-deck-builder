// Package user: repository.go - Describe repository interface to interact with user datastore
package user

import (
	"context"

	"github.com/meplos/locana-deck-builder/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, u domain.User) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
