package router

import (
	"net/http"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/handler"
	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/service"
)

// NewRouter creates and returns a configured http.ServeMux with CORS support.
func NewRouter(
	authService *service.AuthService,
	datosGovService *service.DatosGovService,
	savedProposalService *service.SavedProposalService,
	proposalService *service.ProposalService,
) http.Handler {
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

	// Logout - auth required (POST)
	logoutHandler := handler.NewLogoutHandler(authService)
	mux.Handle("/logout", handler.AuthMiddleware(authService, logoutHandler))

	// Proposals - auth required, fetches from datos.gov.co
	proposalHandler := handler.NewProposalHandler(datosGovService)
	mux.Handle("/public-proposals", handler.AuthMiddleware(authService, proposalHandler))

	// Saved proposals - auth required
	savedProposalHandler := handler.NewSavedProposalHandler(savedProposalService)
	mux.Handle("/saved_proposals", handler.AuthMiddleware(authService, savedProposalHandler))
	mux.Handle("/saved-proposals", handler.AuthMiddleware(authService, savedProposalHandler))

	// Proposal save - auth required, upserts a full proposal
	proposalSaveHandler := handler.NewProposalSaveHandler(proposalService)
	mux.Handle("/proposal-save", handler.AuthMiddleware(authService, proposalSaveHandler))

	return corsMiddleware(mux)
}

// corsMiddleware adds CORS headers to allow cross-origin requests
// from file:// and HTTP frontend origins during development.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
