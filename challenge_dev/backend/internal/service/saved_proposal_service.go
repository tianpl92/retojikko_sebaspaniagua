package service

import (
	"context"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/domain"
	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/repository"
)

// SavedProposalService handles business logic for saved proposals.
type SavedProposalService struct {
	savedRepo repository.SavedProposalRepository
}

func NewSavedProposalService(savedRepo repository.SavedProposalRepository) *SavedProposalService {
	return &SavedProposalService{savedRepo: savedRepo}
}

// SaveProposal saves a proposal for a user.
func (s *SavedProposalService) SaveProposal(ctx context.Context, userID, publicCallID string) (*domain.SavedProposal, error) {
	return s.savedRepo.Save(ctx, userID, publicCallID)
}

// GetSavedProposals returns all saved proposals for a user.
func (s *SavedProposalService) GetSavedProposals(ctx context.Context, userID string) ([]*domain.SavedProposal, error) {
	return s.savedRepo.FindByUserID(ctx, userID)
}

// RemoveSavedProposal removes a saved proposal for a user.
func (s *SavedProposalService) RemoveSavedProposal(ctx context.Context, userID, publicCallID string) error {
	return s.savedRepo.Delete(ctx, userID, publicCallID)
}
