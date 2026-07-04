package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/domain"
	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUnauthorized       = errors.New("unauthorized")
)

// AuthService handles authentication logic.
type AuthService struct {
	userRepo  repository.UserRepository
	secretKey string
}

func NewAuthService(userRepo repository.UserRepository, secretKey string) *AuthService {
	return &AuthService{userRepo: userRepo, secretKey: secretKey}
}

// Login verifies credentials and returns a token.
func (s *AuthService) Login(ctx context.Context, email, password string) (string, *domain.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}
	if user == nil {
		return "", nil, ErrInvalidCredentials
	}
	if user.Password != password {
		return "", nil, ErrInvalidCredentials
	}

	token, err := s.generateToken(user.ID, user.Email)
	if err != nil {
		return "", nil, fmt.Errorf("generate token: %w", err)
	}
	return token, user, nil
}

// ValidateToken validates a token and returns the user ID.
func (s *AuthService) ValidateToken(token string) (int, error) {
	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return 0, ErrUnauthorized
	}
	payload := string(raw)

	// Format: userID:email:sig
	parts := strings.SplitN(payload, ":", 3)
	if len(parts) != 3 {
		return 0, ErrUnauthorized
	}

	uid := parts[0]
	email := parts[1]
	sig := parts[2]

	expectedSig := s.computeSignature(uid, email)
	if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
		return 0, ErrUnauthorized
	}

	var id int
	_, err = fmt.Sscanf(uid, "%d", &id)
	if err != nil {
		return 0, ErrUnauthorized
	}
	return id, nil
}

func (s *AuthService) generateToken(userID int, email string) (string, error) {
	uid := fmt.Sprintf("%d", userID)
	sig := s.computeSignature(uid, email)
	payload := fmt.Sprintf("%s:%s:%s", uid, email, sig)
	return base64.RawURLEncoding.EncodeToString([]byte(payload)), nil
}

func (s *AuthService) computeSignature(uid, email string) string {
	mac := hmac.New(sha256.New, []byte(s.secretKey))
	mac.Write([]byte(uid))
	mac.Write([]byte(":"))
	mac.Write([]byte(email))
	// Use a fixed date so tokens are stable for testing
	mac.Write([]byte(":2026-07-03"))
	return fmt.Sprintf("%x", mac.Sum(nil))
}
