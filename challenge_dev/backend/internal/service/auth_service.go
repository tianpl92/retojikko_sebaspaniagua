package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/domain"
	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrDuplicateDocument  = errors.New("User already exists")
	ErrEmailAlreadyExists = errors.New("email already in use")
	ErrDuplicateEmail     = ErrEmailAlreadyExists
)

type AuthService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	secretKey   string
}

func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository, secretKey string) *AuthService {
	return &AuthService{userRepo: userRepo, sessionRepo: sessionRepo, secretKey: secretKey}
}

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, *domain.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}
	if user == nil {
		return "", nil, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	token, expiresAt, err := s.generateToken(user.ID, user.Email)
	if err != nil {
		return "", nil, fmt.Errorf("generate token: %w", err)
	}
	if err := s.sessionRepo.Create(ctx, user.ID, token, expiresAt); err != nil {
		return "", nil, fmt.Errorf("store session: %w", err)
	}
	return token, user, nil
}

func (s *AuthService) ValidateToken(token string) (string, error) {
	claims := &JWTClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return "", ErrUnauthorized
	}
	if !parsedToken.Valid {
		return "", ErrUnauthorized
	}

	session, err := s.sessionRepo.FindByToken(context.Background(), token)
	if err != nil {
		return "", ErrUnauthorized
	}
	if session == nil || time.Now().After(session.ExpiresAt) {
		return "", ErrUnauthorized
	}
	return claims.UserID, nil
}

func (s *AuthService) generateToken(userID, email string) (string, time.Time, error) {
	expiresAt := time.Now().Add(1 * time.Hour)
	claims := &JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "public-calls-portal",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", time.Time{}, err
	}
	return tokenStr, expiresAt, nil
}

func (s *AuthService) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	exists, err := s.userRepo.ExistsByDocument(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("check document: %w", err)
	}
	if exists {
		return nil, ErrDuplicateDocument
	}
	exists, err = s.userRepo.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return nil, fmt.Errorf("check email: %w", err)
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}

func (s *AuthService) UpdateUser(ctx context.Context, user *domain.User) error {
	return s.userRepo.Update(ctx, user)
}
