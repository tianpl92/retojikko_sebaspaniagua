package repository

import (
	"context"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/domain"
)

// ProposalRepository defines the interface for proposal data access.
type ProposalRepository interface {
	List(ctx context.Context, limit, offset int) ([]*domain.PublicCallProposal, error)
	ListWithFilters(ctx context.Context, query, fase, entidad string, limit, offset int) ([]*domain.PublicCallProposal, error)
	FindByID(ctx context.Context, id string) (*domain.PublicCallProposal, error)
	Upsert(ctx context.Context, proposal *domain.PublicCallProposal) (*domain.PublicCallProposal, error)
}
