package repository

import (
	"context"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/domain"
)

// SavedProposalRepository defines the interface for saved proposals data access.
type SavedProposalRepository interface {
	Save(ctx context.Context, userID, publicCallID string) (*domain.SavedProposal, error)
	FindByUserID(ctx context.Context, userID string) ([]*domain.SavedProposal, error)
	Delete(ctx context.Context, userID, publicCallID string) error
}
