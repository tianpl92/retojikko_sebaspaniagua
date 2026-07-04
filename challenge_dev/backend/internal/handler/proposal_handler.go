package handler

import (
	"encoding/json"
	"net/http"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/service"
)

// ProposalHandler handles public proposals endpoints.
type ProposalHandler struct {
	proposalService *service.ProposalService
}

func NewProposalHandler(proposalService *service.ProposalService) *ProposalHandler {
	return &ProposalHandler{proposalService: proposalService}
}

func (h *ProposalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("query")
	category := r.URL.Query().Get("category")
	fase := r.URL.Query().Get("fase")

	var proposals interface{}
	var err error

	if query != "" || category != "" || fase != "" {
		proposals, err = h.proposalService.FilterProposals(r.Context(), query, category, fase)
	} else {
		proposals, err = h.proposalService.ListProposals(r.Context())
	}

	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(proposals)
}
