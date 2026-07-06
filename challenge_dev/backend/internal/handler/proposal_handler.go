package handler

import (
	"net/http"
	"strconv"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/service"
)

// ProposalHandler handles public proposals endpoints by fetching from datos.gov.co.
type ProposalHandler struct {
	datosGovService *service.DatosGovService
}

func NewProposalHandler(datosGovService *service.DatosGovService) *ProposalHandler {
	return &ProposalHandler{datosGovService: datosGovService}
}

func (h *ProposalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
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

	proposals, err := h.datosGovService.FetchProposals(query, fase, entidad, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch proposals: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, proposals)
}
