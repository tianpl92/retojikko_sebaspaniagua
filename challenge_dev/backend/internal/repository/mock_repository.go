package repository

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/domain"
)

var (
	ErrDuplicateDocument = errors.New("document already exists")
)

// mockUserRepository is an in-memory implementation of UserRepository for development/testing.
type mockUserRepository struct {
	mu      sync.RWMutex
	users   map[string]*domain.User // keyed by email
	byDoc   map[string]*domain.User // keyed by document (ID)
	counter int
}

// NewMockUserRepository creates a new mock user repository with a default test user.
func NewMockUserRepository() UserRepository {
	repo := &mockUserRepository{
		users: make(map[string]*domain.User),
		byDoc: make(map[string]*domain.User),
	}
	// Pre-seed with a default test user for development
	repo.users["test@example.com"] = &domain.User{
		ID:        "1",
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
		Password:  "$2a$10$IqHCg3FyYUqys/g/InYzouZ5sAwUG/3REhJ5NmWLQVocIsQaVRaQe",
		Status:    "active",
	}
	repo.byDoc["1"] = repo.users["test@example.com"]
	return repo
}

func (r *mockUserRepository) FindByID(_ context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.byDoc[id]
	if !ok {
		return nil, ErrNotFound
	}
	return user, nil
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
	if user.ID != "" {
		if _, exists := r.byDoc[user.ID]; exists {
			return ErrDuplicateDocument
		}
	}
	r.counter++
	if user.ID == "" {
		user.ID = fmt.Sprintf("%d", r.counter)
	}
	r.users[user.Email] = user
	r.byDoc[user.ID] = user
	return nil
}

func (r *mockUserRepository) ExistsByDocument(_ context.Context, document string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.byDoc[document]
	return ok, nil
}

func (r *mockUserRepository) ExistsByEmail(_ context.Context, email string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.users[email]
	return ok, nil
}

func (r *mockUserRepository) Update(_ context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.byDoc[user.ID]
	if !ok {
		return ErrNotFound
	}
	// If email changed, remove old email key
	if existing.Email != user.Email {
		delete(r.users, existing.Email)
	}
	r.users[user.Email] = user
	r.byDoc[user.ID] = user
	return nil
}

// mockProposalRepository is an in-memory implementation of ProposalRepository for development/testing.
type mockProposalRepository struct {
	mu        sync.RWMutex
	proposals []*domain.PublicCallProposal
}

// NewMockProposalRepository creates a new mock proposal repository.
func NewMockProposalRepository() ProposalRepository {
	return &mockProposalRepository{
		proposals: []*domain.PublicCallProposal{},
	}
}

func (r *mockProposalRepository) List(_ context.Context, limit, offset int) ([]*domain.PublicCallProposal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if offset >= len(r.proposals) {
		return []*domain.PublicCallProposal{}, nil
	}
	end := offset + limit
	if end > len(r.proposals) || limit <= 0 {
		end = len(r.proposals)
	}
	if limit <= 0 {
		result := make([]*domain.PublicCallProposal, len(r.proposals))
		copy(result, r.proposals)
		return result, nil
	}
	result := make([]*domain.PublicCallProposal, end-offset)
	copy(result, r.proposals[offset:end])
	return result, nil
}

func (r *mockProposalRepository) ListWithFilters(_ context.Context, query, fase, entidad string, limit, offset int) ([]*domain.PublicCallProposal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	// Simplified mock: return all proposals (clientside filtering would happen in real impl)
	return r.List(context.TODO(), limit, offset)
}

func (r *mockProposalRepository) FindByID(_ context.Context, id string) (*domain.PublicCallProposal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.proposals {
		if fmt.Sprintf("%d", p.ID) == id {
			return p, nil
		}
	}
	return nil, ErrNotFound
}

// mockSavedProposalRepository is an in-memory implementation of SavedProposalRepository.
type mockSavedProposalRepository struct {
	mu      sync.RWMutex
	saved   map[string]*domain.SavedProposal // keyed by ID
	byUser  map[string][]string              // userID -> list of saved IDs
	counter int
}

// NewMockSavedProposalRepository creates a new mock saved proposal repository.
func NewMockSavedProposalRepository() SavedProposalRepository {
	return &mockSavedProposalRepository{
		saved:  make(map[string]*domain.SavedProposal),
		byUser: make(map[string][]string),
	}
}

func (r *mockSavedProposalRepository) Save(_ context.Context, userID, publicCallID string) (*domain.SavedProposal, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counter++
	id := fmt.Sprintf("mock-uuid-%d", r.counter)
	saved := &domain.SavedProposal{
		ID:           id,
		PublicCallID: parseInt(publicCallID),
		UserID:       userID,
	}
	r.saved[id] = saved
	r.byUser[userID] = append(r.byUser[userID], id)
	return saved, nil
}

func (r *mockSavedProposalRepository) FindByUserID(_ context.Context, userID string) ([]*domain.SavedProposal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := r.byUser[userID]
	result := make([]*domain.SavedProposal, 0, len(ids))
	for _, id := range ids {
		if s, ok := r.saved[id]; ok {
			result = append(result, s)
		}
	}
	return result, nil
}

func (r *mockSavedProposalRepository) Delete(_ context.Context, userID, publicCallID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, s := range r.saved {
		if s.UserID == userID && fmt.Sprintf("%d", s.PublicCallID) == publicCallID {
			delete(r.saved, id)
			ids := r.byUser[userID]
			for i, sid := range ids {
				if sid == id {
					r.byUser[userID] = append(ids[:i], ids[i+1:]...)
					break
				}
			}
			return nil
		}
	}
	return ErrNotFound
}

func parseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

// mockSessionRepository is an in-memory implementation of SessionRepository for testing.
type mockSessionRepository struct {
	mu       sync.RWMutex
	sessions map[string]*domain.UserSession // keyed by token
	counter  int
}

// NewMockSessionRepository creates a new mock session repository.
func NewMockSessionRepository() SessionRepository {
	return &mockSessionRepository{
		sessions: make(map[string]*domain.UserSession),
	}
}

func (r *mockSessionRepository) Create(_ context.Context, userID, token string, expiresAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counter++
	session := &domain.UserSession{
		ID:        fmt.Sprintf("session-%d", r.counter),
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
	r.sessions[token] = session
	return nil
}

func (r *mockSessionRepository) FindByToken(_ context.Context, token string) (*domain.UserSession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	session, ok := r.sessions[token]
	if !ok {
		return nil, ErrNotFound
	}
	return session, nil
}

func (r *mockSessionRepository) DeleteExpired(_ context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	for token, session := range r.sessions {
		if session.ExpiresAt.Before(now) {
			delete(r.sessions, token)
		}
	}
	return nil
}
