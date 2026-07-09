package service

import (
	"context"
	"testing"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/domain"
	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/repository"
)

func TestAuthService_Login_Success(t *testing.T) {
	userRepo := repository.NewMockUserRepository()
	sessionRepo := repository.NewMockSessionRepository()
	authService := NewAuthService(userRepo, sessionRepo, "test-secret")

	token, user, err := authService.Login(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %q", user.Email)
	}
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	userRepo := repository.NewMockUserRepository()
	sessionRepo := repository.NewMockSessionRepository()
	authService := NewAuthService(userRepo, sessionRepo, "test-secret")

	_, _, err := authService.Login(context.Background(), "test@example.com", "wrongpassword")
	if err != ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	userRepo := repository.NewMockUserRepository()
	sessionRepo := repository.NewMockSessionRepository()
	authService := NewAuthService(userRepo, sessionRepo, "test-secret")

	_, _, err := authService.Login(context.Background(), "nobody@example.com", "password123")
	if err != ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthService_ValidateToken_Valid(t *testing.T) {
	userRepo := repository.NewMockUserRepository()
	sessionRepo := repository.NewMockSessionRepository()
	authService := NewAuthService(userRepo, sessionRepo, "test-secret")

	token, _, err := authService.Login(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	userID, err := authService.ValidateToken(token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if userID != "1" {
		t.Errorf("expected userID '1', got %q", userID)
	}
}

func TestAuthService_ValidateToken_Invalid(t *testing.T) {
	userRepo := repository.NewMockUserRepository()
	sessionRepo := repository.NewMockSessionRepository()
	authService := NewAuthService(userRepo, sessionRepo, "test-secret")

	_, err := authService.ValidateToken("invalid-token")
	if err != ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestAuthService_CreateUser(t *testing.T) {
	userRepo := repository.NewMockUserRepository()
	sessionRepo := repository.NewMockSessionRepository()
	authService := NewAuthService(userRepo, sessionRepo, "test-secret")

	user := &domain.User{
		ID:        "new-doc-123",
		FirstName: "New",
		LastName:  "User",
		Email:     "newuser@example.com",
		Password:  "secret",
		Status:    "active",
	}

	created, err := authService.CreateUser(context.Background(), user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created == nil {
		t.Fatal("expected created user, got nil")
	}
	if created.Email != "newuser@example.com" {
		t.Errorf("expected email 'newuser@example.com', got %q", created.Email)
	}
}

func TestAuthService_GetUserByID(t *testing.T) {
	userRepo := repository.NewMockUserRepository()
	sessionRepo := repository.NewMockSessionRepository()
	authService := NewAuthService(userRepo, sessionRepo, "test-secret")

	user, err := authService.GetUserByID(context.Background(), "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %q", user.Email)
	}
}

func TestProposalService_ListProposals(t *testing.T) {
	proposalRepo := repository.NewMockProposalRepository()
	proposalService := NewProposalService(proposalRepo)

	proposals, err := proposalService.ListProposals(context.Background(), 10, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if proposals == nil {
		t.Error("expected non-nil slice, got nil")
	}
}

func TestProposalService_ListWithFilters(t *testing.T) {
	proposalRepo := repository.NewMockProposalRepository()
	proposalService := NewProposalService(proposalRepo)

	proposals, err := proposalService.ListWithFilters(context.Background(), "test", "fase1", "entidad1", 10, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if proposals == nil {
		t.Error("expected non-nil slice, got nil")
	}
}

func TestSavedProposalService_SaveProposal(t *testing.T) {
	savedRepo := repository.NewMockSavedProposalRepository(nil)
	savedService := NewSavedProposalService(savedRepo)

	saved, err := savedService.SaveProposal(context.Background(), "1", "42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if saved == nil {
		t.Fatal("expected saved proposal, got nil")
	}
	if saved.PublicCallID != "42" {
		t.Errorf("expected PublicCallID 42, got %s", saved.PublicCallID)
	}
	if saved.UserID != "1" {
		t.Errorf("expected UserID '1', got %q", saved.UserID)
	}
}

func TestSavedProposalService_GetSavedProposals(t *testing.T) {
	savedRepo := repository.NewMockSavedProposalRepository(nil)
	savedService := NewSavedProposalService(savedRepo)

	// No saved proposals initially
	saved, err := savedService.GetSavedProposals(context.Background(), "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if saved == nil {
		t.Error("expected non-nil slice, got nil")
	}
	if len(saved) != 0 {
		t.Errorf("expected 0 saved, got %d", len(saved))
	}

	// Save one and check
	savedService.SaveProposal(context.Background(), "1", "42")
	saved, err = savedService.GetSavedProposals(context.Background(), "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(saved) != 1 {
		t.Errorf("expected 1 saved, got %d", len(saved))
	}
}

func TestSavedProposalService_RemoveSavedProposal(t *testing.T) {
	savedRepo := repository.NewMockSavedProposalRepository(nil)
	savedService := NewSavedProposalService(savedRepo)

	saved, _ := savedService.SaveProposal(context.Background(), "1", "42")

	err := savedService.RemoveSavedProposal(context.Background(), saved.UserID, "42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	savedList, _ := savedService.GetSavedProposals(context.Background(), "1")
	if len(savedList) != 0 {
		t.Errorf("expected 0 after delete, got %d", len(savedList))
	}
}
