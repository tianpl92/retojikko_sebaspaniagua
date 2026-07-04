package service

import (
	"context"
	"testing"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/repository"
)

func TestAuthService_Login_Success(t *testing.T) {
	userRepo := repository.NewMockUserRepository()
	authService := NewAuthService(userRepo, "test-secret")

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
	authService := NewAuthService(userRepo, "test-secret")

	_, _, err := authService.Login(context.Background(), "test@example.com", "wrongpassword")
	if err != ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	userRepo := repository.NewMockUserRepository()
	authService := NewAuthService(userRepo, "test-secret")

	_, _, err := authService.Login(context.Background(), "nobody@example.com", "password123")
	if err != ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthService_ValidateToken_Valid(t *testing.T) {
	userRepo := repository.NewMockUserRepository()
	authService := NewAuthService(userRepo, "test-secret")

	token, _, err := authService.Login(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	userID, err := authService.ValidateToken(token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if userID != 1 {
		t.Errorf("expected userID 1, got %d", userID)
	}
}

func TestAuthService_ValidateToken_Invalid(t *testing.T) {
	userRepo := repository.NewMockUserRepository()
	authService := NewAuthService(userRepo, "test-secret")

	_, err := authService.ValidateToken("invalid-token")
	if err != ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestProposalService_ListProposals(t *testing.T) {
	proposalRepo := repository.NewMockProposalRepository()
	proposalService := NewProposalService(proposalRepo)

	proposals, err := proposalService.ListProposals(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if proposals == nil {
		t.Error("expected non-nil slice, got nil")
	}
}

func TestProposalService_FilterProposals(t *testing.T) {
	proposalRepo := repository.NewMockProposalRepository()
	proposalService := NewProposalService(proposalRepo)

	proposals, err := proposalService.FilterProposals(context.Background(), "test", "category1", "fase1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if proposals == nil {
		t.Error("expected non-nil slice, got nil")
	}
}

func TestSavedProposalService_SaveProposal(t *testing.T) {
	savedRepo := repository.NewMockSavedProposalRepository()
	savedService := NewSavedProposalService(savedRepo)

	saved, err := savedService.SaveProposal(context.Background(), 1, 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if saved == nil {
		t.Fatal("expected saved proposal, got nil")
	}
	if saved.PublicCallID != 42 {
		t.Errorf("expected PublicCallID 42, got %d", saved.PublicCallID)
	}
	if saved.UserID != 1 {
		t.Errorf("expected UserID 1, got %d", saved.UserID)
	}
}

func TestSavedProposalService_GetSavedProposals(t *testing.T) {
	savedRepo := repository.NewMockSavedProposalRepository()
	savedService := NewSavedProposalService(savedRepo)

	// No saved proposals initially
	saved, err := savedService.GetSavedProposals(context.Background(), 1)
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
	savedService.SaveProposal(context.Background(), 1, 42)
	saved, err = savedService.GetSavedProposals(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(saved) != 1 {
		t.Errorf("expected 1 saved, got %d", len(saved))
	}
}

func TestSavedProposalService_RemoveSavedProposal(t *testing.T) {
	savedRepo := repository.NewMockSavedProposalRepository()
	savedService := NewSavedProposalService(savedRepo)

	saved, _ := savedService.SaveProposal(context.Background(), 1, 42)

	err := savedService.RemoveSavedProposal(context.Background(), saved.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	savedList, _ := savedService.GetSavedProposals(context.Background(), 1)
	if len(savedList) != 0 {
		t.Errorf("expected 0 after delete, got %d", len(savedList))
	}
}
