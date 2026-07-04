package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/domain"
)

// mockUserRepository is an in-memory implementation of UserRepository for development/testing.
type mockUserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.User // keyed by email
}

// NewMockUserRepository creates a new mock user repository with a default test user.
func NewMockUserRepository() UserRepository {
	repo := &mockUserRepository{
		users: make(map[string]*domain.User),
	}
	// Pre-seed with a default test user for development
	repo.users["test@example.com"] = &domain.User{
		ID:        1,
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
		Password:  "password123",
		Status:    "active",
	}
	return repo
}

func (r *mockUserRepository) FindByEmail(_ context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[email]
	if !ok {
		return nil, ErrNotFound
	}
	return user, nil
}

func (r *mockUserRepository) Create(_ context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.users[user.Email]; exists {
		return ErrDuplicateEmail
	}
	r.users[user.Email] = user
	return nil
}

// mockProposalRepository is an in-memory implementation of ProposalRepository for development/testing.
type mockProposalRepository struct {
	mu        sync.RWMutex
	proposals []domain.PublicCallProposal
}

// NewMockProposalRepository creates a new mock proposal repository.
func NewMockProposalRepository() ProposalRepository {
	return &mockProposalRepository{
		proposals: []domain.PublicCallProposal{},
	}
}

func (r *mockProposalRepository) List(_ context.Context) ([]domain.PublicCallProposal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]domain.PublicCallProposal, len(r.proposals))
	copy(result, r.proposals)
	return result, nil
}

func (r *mockProposalRepository) Filter(_ context.Context, query, category, fase string) ([]domain.PublicCallProposal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	// Simplified mock: return all proposals (clientside filtering would happen in real impl)
	result := make([]domain.PublicCallProposal, len(r.proposals))
	copy(result, r.proposals)
	return result, nil
}

func (r *mockProposalRepository) FindByID(_ context.Context, id int) (*domain.PublicCallProposal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.proposals {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, ErrNotFound
}

// mockSavedProposalRepository is an in-memory implementation of SavedProposalRepository.
type mockSavedProposalRepository struct {
	mu      sync.RWMutex
	saved   map[string]*domain.SavedProposal // keyed by ID
	byUser  map[int][]string                 // userID -> list of saved IDs
	counter int
}

// NewMockSavedProposalRepository creates a new mock saved proposal repository.
func NewMockSavedProposalRepository() SavedProposalRepository {
	return &mockSavedProposalRepository{
		saved:  make(map[string]*domain.SavedProposal),
		byUser: make(map[int][]string),
	}
}

func (r *mockSavedProposalRepository) Save(_ context.Context, saved *domain.SavedProposal) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counter++
	id := fmt.Sprintf("mock-uuid-%d", r.counter)
	saved.ID = id
	r.saved[id] = saved
	r.byUser[saved.UserID] = append(r.byUser[saved.UserID], id)
	return nil
}

func (r *mockSavedProposalRepository) FindByUserID(_ context.Context, userID int) ([]domain.SavedProposal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := r.byUser[userID]
	result := make([]domain.SavedProposal, 0, len(ids))
	for _, id := range ids {
		if s, ok := r.saved[id]; ok {
			result = append(result, *s)
		}
	}
	return result, nil
}

func (r *mockSavedProposalRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.saved[id]; !ok {
		return ErrNotFound
	}
	delete(r.saved, id)
	for uid, ids := range r.byUser {
		for i, sid := range ids {
			if sid == id {
				r.byUser[uid] = append(ids[:i], ids[i+1:]...)
				break
			}
		}
	}
	return nil
}
