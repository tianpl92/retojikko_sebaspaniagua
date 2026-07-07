package handler

import (
	"encoding/json"
	"net/http"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/domain"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/service"
)

type saveProposalRequest struct {
	PublicCallID string `json:"public_call_id"`
}

// SavedProposalHandler handles saved proposals endpoints.
type SavedProposalHandler struct {
	savedProposalService *service.SavedProposalService
}

func NewSavedProposalHandler(savedProposalService *service.SavedProposalService) *SavedProposalHandler {
	return &SavedProposalHandler{savedProposalService: savedProposalService}
}

func (h *SavedProposalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok || userID == "" {
		writeError(w, http.StatusUnauthorized, "Credentials invalid")
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

func (h *SavedProposalHandler) save(w http.ResponseWriter, r *http.Request, userID string) {
	var req saveProposalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.PublicCallID == "" {
		writeError(w, http.StatusBadRequest, "public_call_id is required")
		return
	}

	// Idempotent: always return success, even if already saved
	_, _ = h.savedProposalService.SaveProposal(r.Context(), userID, req.PublicCallID)

	writeJSON(w, http.StatusCreated, map[string]string{
		"message": "Guardado satisfactoriamente",
	})
}

func (h *SavedProposalHandler) list(w http.ResponseWriter, r *http.Request, userID string) {
	proposals, err := h.savedProposalService.GetSavedProposals(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Return empty array if no saved proposals
	if proposals == nil {
		proposals = make([]*domain.SavedProposal, 0)
	}

	writeJSON(w, http.StatusOK, proposals)
}
