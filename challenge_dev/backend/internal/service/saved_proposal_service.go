package service

import (
	"context"
	"time"

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
func (s *SavedProposalService) SaveProposal(ctx context.Context, userID, publicCallID int) (*domain.SavedProposal, error) {
	saved := &domain.SavedProposal{
		PublicCallID:    publicCallID,
		UserID:          userID,
		AssociationDate: time.Now(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if err := s.savedRepo.Save(ctx, saved); err != nil {
		return nil, err
	}
	return saved, nil
}

// GetSavedProposals returns all saved proposals for a user.
func (s *SavedProposalService) GetSavedProposals(ctx context.Context, userID int) ([]domain.SavedProposal, error) {
	return s.savedRepo.FindByUserID(ctx, userID)
}

// RemoveSavedProposal removes a saved proposal.
func (s *SavedProposalService) RemoveSavedProposal(ctx context.Context, id string) error {
	return s.savedRepo.Delete(ctx, id)
}
