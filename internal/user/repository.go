// Package user: repository.go - Describe repository interface to interact with user datastore
package user

import (
	"context"
	"time"

	"github.com/meplos/locana-deck-builder/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, u domain.User) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	LogUserAt(ctx context.Context, ID string, iat time.Time) error
	FindByID(ctx context.Context, ID string) (domain.User, error)
}
