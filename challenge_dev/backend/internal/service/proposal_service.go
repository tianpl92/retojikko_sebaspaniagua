package service

import (
	"context"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/domain"
	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/repository"
)

// ProposalService handles business logic for proposals.
type ProposalService struct {
	proposalRepo repository.ProposalRepository
}

func NewProposalService(proposalRepo repository.ProposalRepository) *ProposalService {
	return &ProposalService{proposalRepo: proposalRepo}
}

// ListProposals returns all proposals with pagination.
func (s *ProposalService) ListProposals(ctx context.Context, limit, offset int) ([]*domain.PublicCallProposal, error) {
	return s.proposalRepo.List(ctx, limit, offset)
}

// ListWithFilters returns proposals matching filter criteria.
func (s *ProposalService) ListWithFilters(ctx context.Context, query, fase, entidad string, limit, offset int) ([]*domain.PublicCallProposal, error) {
	return s.proposalRepo.ListWithFilters(ctx, query, fase, entidad, limit, offset)
}

// FilterProposals is an alias for ListWithFilters.
func (s *ProposalService) FilterProposals(ctx context.Context, query, fase, entidad string, limit, offset int) ([]*domain.PublicCallProposal, error) {
	return s.ListWithFilters(ctx, query, fase, entidad, limit, offset)
}
