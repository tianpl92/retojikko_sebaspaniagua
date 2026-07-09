package repository

import (
	"context"
	"time"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/domain"
)

// SessionRepository defines the interface for session/token data access.
type SessionRepository interface {
	Create(ctx context.Context, userID, token string, expiresAt time.Time) error
	FindByToken(ctx context.Context, token string) (*domain.UserSession, error)
	DeleteByToken(ctx context.Context, token string) error
	DeleteExpired(ctx context.Context) error
}
