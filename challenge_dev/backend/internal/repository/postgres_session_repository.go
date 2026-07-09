//go:build postgres

package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/domain"
)

// PostgresSessionRepository is a PostgreSQL-backed implementation of SessionRepository.
type PostgresSessionRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresSessionRepository creates a new PostgresSessionRepository.
func NewPostgresSessionRepository(pool *pgxpool.Pool) *PostgresSessionRepository {
	return &PostgresSessionRepository{pool: pool}
}

// Create inserts a new session into the database.
func (r *PostgresSessionRepository) Create(ctx context.Context, userID, token string, expiresAt time.Time) error {
	now := time.Now()
	query := `
		INSERT INTO user_sessions (user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.pool.Exec(ctx, query, userID, token, expiresAt, now)
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil
}

// FindByToken retrieves a session by token.
func (r *PostgresSessionRepository) FindByToken(ctx context.Context, token string) (*domain.UserSession, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM user_sessions
		WHERE token = $1 AND expires_at > NOW()
	`
	row := r.pool.QueryRow(ctx, query, token)
	session := &domain.UserSession{}
	err := row.Scan(
		&session.ID, &session.UserID, &session.Token,
		&session.ExpiresAt, &session.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find session by token: %w", err)
	}
	return session, nil
}

// DeleteExpired removes all expired sessions.
func (r *PostgresSessionRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM user_sessions WHERE expires_at <= NOW()`
	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("delete expired sessions: %w", err)
	}
	return nil
}
