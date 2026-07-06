package repository

import (
	"context"
	"testing"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/domain"
)

func TestMockUserRepository_FindByEmail_Found(t *testing.T) {
	repo := NewMockUserRepository()
	user, err := repo.FindByEmail(context.Background(), "test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %q", user.Email)
	}
	if user.ID != "1" {
		t.Errorf("expected ID \"1\", got %q", user.ID)
	}
}

func TestMockUserRepository_FindByEmail_NotFound(t *testing.T) {
	repo := NewMockUserRepository()
	_, err := repo.FindByEmail(context.Background(), "nonexistent@example.com")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestMockUserRepository_Create(t *testing.T) {
	repo := NewMockUserRepository()
	user := &domain.User{
		ID:       "2",
		Email:    "newuser@example.com",
		Password: "secret",
	}
	err := repo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify it was saved
	found, err := repo.FindByEmail(context.Background(), "newuser@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found.Email != "newuser@example.com" {
		t.Errorf("expected email 'newuser@example.com', got %q", found.Email)
	}
}

func TestMockUserRepository_Create_Duplicate(t *testing.T) {
	repo := NewMockUserRepository()
	user := &domain.User{
		Email:    "test@example.com",
		Password: "another",
	}
	err := repo.Create(context.Background(), user)
	if err != ErrDuplicateEmail {
		t.Errorf("expected ErrDuplicateEmail, got %v", err)
	}
}

func TestMockProposalRepository_List(t *testing.T) {
	repo := NewMockProposalRepository()
	proposals, err := repo.List(context.Background(), 100, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(proposals) != 0 {
		t.Errorf("expected 0 proposals, got %d", len(proposals))
	}
}

func TestMockSavedProposalRepository_SaveAndFind(t *testing.T) {
	repo := NewMockSavedProposalRepository()
	saved, err := repo.Save(context.Background(), "1", "42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if saved.ID == "" {
		t.Error("expected ID to be set after Save")
	}
	if saved.PublicCallID != 42 {
		t.Errorf("expected PublicCallID 42, got %d", saved.PublicCallID)
	}
	if saved.UserID != "1" {
		t.Errorf("expected UserID \"1\", got %q", saved.UserID)
	}

	// Find by user
	results, err := repo.FindByUserID(context.Background(), "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].PublicCallID != 42 {
		t.Errorf("expected PublicCallID 42, got %d", results[0].PublicCallID)
	}
}

func TestMockSavedProposalRepository_Delete(t *testing.T) {
	repo := NewMockSavedProposalRepository()
	saved, err := repo.Save(context.Background(), "2", "99")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Delete
	err = repo.Delete(context.Background(), "2", "99")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	results, err := repo.FindByUserID(context.Background(), "2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 saved proposals after delete, got %d", len(results))
	}
	_ = saved // used to avoid unused variable complaint
}

func TestMockSavedProposalRepository_Delete_NotFound(t *testing.T) {
	repo := NewMockSavedProposalRepository()
	err := repo.Delete(context.Background(), "nonexistent", "999")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestMockUserRepository_FindByID(t *testing.T) {
	repo := NewMockUserRepository()
	user, err := repo.FindByID(context.Background(), "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.ID != "1" {
		t.Errorf("expected ID \"1\", got %q", user.ID)
	}
}

func TestMockUserRepository_ExistsByEmail(t *testing.T) {
	repo := NewMockUserRepository()
	exists, err := repo.ExistsByEmail(context.Background(), "test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Error("expected user to exist")
	}

	exists, err = repo.ExistsByEmail(context.Background(), "nobody@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exists {
		t.Error("expected user not to exist")
	}
}
