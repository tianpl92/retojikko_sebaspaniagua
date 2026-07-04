package repository

import (
	"context"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/domain"
)

// SavedProposalRepository defines the interface for saved proposals data access.
type SavedProposalRepository interface {
	Save(ctx context.Context, saved *domain.SavedProposal) error
	FindByUserID(ctx context.Context, userID int) ([]domain.SavedProposal, error)
	Delete(ctx context.Context, id string) error
}
