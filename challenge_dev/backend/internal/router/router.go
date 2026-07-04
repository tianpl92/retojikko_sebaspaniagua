package router

import (
	"net/http"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/handler"
	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/service"
)

// NewRouter creates and returns a configured http.ServeMux.
func NewRouter(
	authService *service.AuthService,
	proposalService *service.ProposalService,
	savedProposalService *service.SavedProposalService,
) *http.ServeMux {
	mux := http.NewServeMux()

	// Health
	healthHandler := handler.NewHealthHandler()
	mux.Handle("/health", healthHandler)

	// Auth
	authHandler := handler.NewAuthHandler(authService)
	mux.Handle("/login", authHandler)

	// Proposals
	proposalHandler := handler.NewProposalHandler(proposalService)
	mux.Handle("/public-proposals", proposalHandler)

	// Saved proposals
	savedProposalHandler := handler.NewSavedProposalHandler(savedProposalService, authService)
	mux.Handle("/saved-proposals", savedProposalHandler)

	return mux
}
