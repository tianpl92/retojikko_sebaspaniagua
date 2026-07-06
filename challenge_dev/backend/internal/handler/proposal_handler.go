package handler

import (
	"net/http"
	"strconv"

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
	fase := r.URL.Query().Get("fase")
	entidad := r.URL.Query().Get("entidad")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 100
	offset := 0
	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil && v > 0 {
			limit = v
		}
	}
	if offsetStr != "" {
		if v, err := strconv.Atoi(offsetStr); err == nil && v >= 0 {
			offset = v
		}
	}

	var proposals interface{}
	var err error

	if query != "" || fase != "" || entidad != "" {
		proposals, err = h.proposalService.FilterProposals(r.Context(), query, fase, entidad, limit, offset)
	} else {
		proposals, err = h.proposalService.ListProposals(r.Context(), limit, offset)
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, proposals)
}
