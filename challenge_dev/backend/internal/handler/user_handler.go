package handler

import (
	"encoding/json"
	"net/http"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/service"
)

// UserHandler handles user info and modify endpoints.
type UserHandler struct {
	authService *service.AuthService
}

func NewUserHandler(authService *service.AuthService) *UserHandler {
	return &UserHandler{authService: authService}
}

// ServeHTTP dispatches based on method: GET for /user-info, POST for /user-modify
func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.userInfo(w, r)
	case http.MethodPost:
		h.userModify(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) userInfo(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok || userID == "" {
		writeError(w, http.StatusUnauthorized, "Credentials invalid")
		return
	}

	user, err := h.authService.GetUserByID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if user == nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

type modifyUserRequest struct {
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Password    string `json:"password,omitempty"`
	Status      string `json:"status,omitempty"`
}

func (h *UserHandler) userModify(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok || userID == "" {
		writeError(w, http.StatusUnauthorized, "Credentials invalid")
		return
	}

	var req modifyUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Get existing user
	user, err := h.authService.GetUserByID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if user == nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	// Update only provided fields
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Gender != "" {
		user.Gender = req.Gender
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.PhoneNumber != "" {
		user.Phone = req.PhoneNumber
	}
	if req.Password != "" {
		user.Password = req.Password
	}
	if req.Status != "" {
		user.Status = req.Status
	}

	if err := h.authService.UpdateUser(r.Context(), user); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, user)
}

// contextKey type for context values
type contextKey string

const userIDKey contextKey = "userID"

// AuthMiddleware checks token from Authorization header.
func AuthMiddleware(authService *service.AuthService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			writeError(w, http.StatusUnauthorized, "Credentials invalid")
			return
		}
		// Strip "Bearer " prefix if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		} else {
			writeError(w, http.StatusUnauthorized, "Credentials invalid")
			return
		}

		userID, err := authService.ValidateToken(token)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "Credentials invalid")
			return
		}

		// Store userID in context
		ctx := r.Context()
		ctx = contextWithUserID(ctx, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
