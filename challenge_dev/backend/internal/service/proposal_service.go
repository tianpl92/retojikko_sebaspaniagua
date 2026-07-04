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

// ListProposals returns all proposals.
func (s *ProposalService) ListProposals(ctx context.Context) ([]domain.PublicCallProposal, error) {
	return s.proposalRepo.List(ctx)
}

// FilterProposals filters proposals by query, category, and fase.
func (s *ProposalService) FilterProposals(ctx context.Context, query, category, fase string) ([]domain.PublicCallProposal, error) {
	return s.proposalRepo.Filter(ctx, query, category, fase)
}
