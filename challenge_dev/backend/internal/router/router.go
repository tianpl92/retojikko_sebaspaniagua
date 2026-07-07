package router

import (
	"net/http"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/handler"
	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/service"
)

// NewRouter creates and returns a configured http.ServeMux.
func NewRouter(
	authService *service.AuthService,
	datosGovService *service.DatosGovService,
	savedProposalService *service.SavedProposalService,
) *http.ServeMux {
	mux := http.NewServeMux()

	// Health - no auth required
	mux.Handle("/health", handler.NewHealthHandler())

	// Auth - no auth required
	mux.Handle("/login", handler.NewAuthHandler(authService))
	mux.Handle("/user-create", handler.NewCreateUserHandler(authService))

	// User info - auth required (POST)
	userHandler := handler.NewUserHandler(authService)
	mux.Handle("/user-info", handler.AuthMiddleware(authService, userHandler))

	// User modify - auth required (POST)
	userModifyHandler := handler.NewUserModifyHandler(authService)
	mux.Handle("/user-modify", handler.AuthMiddleware(authService, userModifyHandler))

	// Proposals - auth required, fetches from datos.gov.co
	proposalHandler := handler.NewProposalHandler(datosGovService)
	mux.Handle("/public-proposals", handler.AuthMiddleware(authService, proposalHandler))

	// Saved proposals - auth required
	savedProposalHandler := handler.NewSavedProposalHandler(savedProposalService)
	mux.Handle("/saved_proposals", handler.AuthMiddleware(authService, savedProposalHandler))
	mux.Handle("/saved-proposals", handler.AuthMiddleware(authService, savedProposalHandler))

	return mux
}
