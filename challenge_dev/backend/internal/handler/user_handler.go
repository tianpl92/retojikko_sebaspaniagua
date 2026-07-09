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

// ServeHTTP dispatches based on method: POST for /user-info, POST for /user-modify
func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// Use path to determine which handler: /user-info or /user-modify
		// Since both share the same handler, we check the request path
		h.userInfo(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// UserModifyHandler handles POST /user-modify separately
type UserModifyHandler struct {
	authService *service.AuthService
}

func NewUserModifyHandler(authService *service.AuthService) *UserModifyHandler {
	return &UserModifyHandler{authService: authService}
}

func (h *UserModifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.userModify(w, r)
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
		writeError(w, http.StatusNotFound, "Usuario no registrado")
		return
	}

	// Create response without password field
	resp := map[string]interface{}{
		"id":           user.ID,
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"gender":       user.Gender,
		"email":        user.Email,
		"phone_number": user.Phone,
		"status":       user.Status,
		"created_at":   user.CreatedAt,
		"updated_at":   user.UpdatedAt,
		"deleted_at":   user.DeletedAt,
	}

	writeJSON(w, http.StatusOK, resp)
}

type modifyUserRequest struct {
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Gender      string `json:"gender,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Status      string `json:"status,omitempty"`
}

func (h *UserModifyHandler) userModify(w http.ResponseWriter, r *http.Request) {
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
		writeError(w, http.StatusNotFound, "Usuario no registrado")
		return
	}

	hasAllowedField := false

	if req.FirstName != "" {
		hasAllowedField = true
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		hasAllowedField = true
		user.LastName = req.LastName
	}
	if req.Gender != "" {
		hasAllowedField = true
		user.Gender = req.Gender
	}
	if req.PhoneNumber != "" {
		hasAllowedField = true
		user.Phone = req.PhoneNumber
	}
	if req.Status != "" {
		hasAllowedField = true
		user.Status = req.Status
	}

	if !hasAllowedField {
		writeError(w, http.StatusBadRequest, "Debe modificarse un campo por lo menos para actualizar")
		return
	}

	if err := h.authService.UpdateUser(r.Context(), user); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Informacion actualizada correctamente",
	})
}
