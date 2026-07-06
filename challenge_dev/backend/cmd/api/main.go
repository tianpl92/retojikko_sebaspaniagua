package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/config"
	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/repository"
	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/router"
	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/service"
)

func main() {
	cfg := config.Load()

	userRepo := repository.NewMockUserRepository()
	proposalRepo := repository.NewMockProposalRepository()
	savedProposalRepo := repository.NewMockSavedProposalRepository()
	sessionRepo := repository.NewMockSessionRepository()

	authService := service.NewAuthService(userRepo, sessionRepo, cfg.SecretKey)
	proposalService := service.NewProposalService(proposalRepo)
	savedProposalService := service.NewSavedProposalService(savedProposalRepo)

	mux := router.NewRouter(authService, proposalService, savedProposalService)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting server on %s", addr)

	if _, err := os.Stat(".env"); err == nil {
		log.Println("Found .env file — ensure environment variables are set")
	}

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
