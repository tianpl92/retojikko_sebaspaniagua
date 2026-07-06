package repository

import (
	"context"
	"errors"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/domain"
)

var (
	ErrNotFound       = errors.New("record not found")
	ErrDuplicateEmail = errors.New("email already exists")
)

// UserRepository defines the interface for user data access.
type UserRepository interface {
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	ExistsByDocument(ctx context.Context, document string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Update(ctx context.Context, user *domain.User) error
}
