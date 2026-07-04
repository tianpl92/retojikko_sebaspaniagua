package repository

import (
	"context"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/domain"
)

// ProposalRepository defines the interface for proposal data access.
type ProposalRepository interface {
	List(ctx context.Context) ([]domain.PublicCallProposal, error)
	Filter(ctx context.Context, query, category, fase string) ([]domain.PublicCallProposal, error)
	FindByID(ctx context.Context, id int) (*domain.PublicCallProposal, error)
}
