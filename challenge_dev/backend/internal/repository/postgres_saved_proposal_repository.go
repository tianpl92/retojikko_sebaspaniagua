package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/domain"
)

// PostgresSavedProposalRepository is a PostgreSQL-backed implementation of SavedProposalRepository.
type PostgresSavedProposalRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresSavedProposalRepository creates a new PostgresSavedProposalRepository.
func NewPostgresSavedProposalRepository(pool *pgxpool.Pool) *PostgresSavedProposalRepository {
	return &PostgresSavedProposalRepository{pool: pool}
}

// Save creates a new saved proposal association.
func (r *PostgresSavedProposalRepository) Save(ctx context.Context, userID, publicCallID string) (*domain.SavedProposal, error) {
	pcID, err := strconv.Atoi(publicCallID)
	if err != nil {
		return nil, fmt.Errorf("invalid public call id: %w", err)
	}

	now := time.Now()
	query := `
		INSERT INTO public_call_user_associations
			(user_id, public_call_id, association_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var id string
	err = r.pool.QueryRow(ctx, query,
		userID, pcID, now, now, now,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("save saved proposal: %w", err)
	}

	return &domain.SavedProposal{
		ID:              id,
		PublicCallID:    pcID,
		UserID:          userID,
		AssociationDate: now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

// FindByUserID retrieves all saved proposals for a user.
func (r *PostgresSavedProposalRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.SavedProposal, error) {
	query := `
		SELECT id, public_call_id, user_id, association_date,
		       created_at, updated_at, deleted_at
		FROM public_call_user_associations
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("find saved proposals by user: %w", err)
	}
	defer rows.Close()

	var results []*domain.SavedProposal
	for rows.Next() {
		s := &domain.SavedProposal{}
		err := rows.Scan(
			&s.ID, &s.PublicCallID, &s.UserID,
			&s.AssociationDate, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan saved proposal: %w", err)
		}
		results = append(results, s)
	}
	return results, rows.Err()
}

// Delete removes a saved proposal association.
func (r *PostgresSavedProposalRepository) Delete(ctx context.Context, userID, publicCallID string) error {
	pcID, err := strconv.Atoi(publicCallID)
	if err != nil {
		return fmt.Errorf("invalid public call id: %w", err)
	}

	query := `
		DELETE FROM public_call_user_associations
		WHERE user_id = $1 AND public_call_id = $2
	`
	ct, err := r.pool.Exec(ctx, query, userID, pcID)
	if err != nil {
		return fmt.Errorf("delete saved proposal: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
