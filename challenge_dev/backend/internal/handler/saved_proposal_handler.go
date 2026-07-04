package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/service"
)

type saveProposalRequest struct {
	PublicCallID int `json:"public_call_id"`
}

// SavedProposalHandler handles saved proposals endpoints.
type SavedProposalHandler struct {
	savedProposalService *service.SavedProposalService
	authService          *service.AuthService
}

func NewSavedProposalHandler(savedProposalService *service.SavedProposalService, authService *service.AuthService) *SavedProposalHandler {
	return &SavedProposalHandler{savedProposalService: savedProposalService, authService: authService}
}

func (h *SavedProposalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID, err := h.authenticate(r)
	if err != nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodPost:
		h.save(w, r, userID)
	case http.MethodGet:
		h.list(w, r, userID)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *SavedProposalHandler) authenticate(r *http.Request) (int, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return 0, service.ErrUnauthorized
	}
	// Strip "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}
	return h.authService.ValidateToken(token)
}

func (h *SavedProposalHandler) save(w http.ResponseWriter, r *http.Request, userID int) {
	var req saveProposalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	saved, err := h.savedProposalService.SaveProposal(r.Context(), userID, req.PublicCallID)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(saved)
}

func (h *SavedProposalHandler) list(w http.ResponseWriter, r *http.Request, userID int) {
	saved, err := h.savedProposalService.GetSavedProposals(r.Context(), userID)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(saved)
}

func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
