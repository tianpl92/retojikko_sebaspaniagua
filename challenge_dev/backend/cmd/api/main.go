package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/config"
	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/database"
	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/repository"
	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/router"
	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/service"
)

func main() {
	cfg := config.Load()

	// Attempt database migration on startup (non-fatal if PostgreSQL unavailable)
	if err := runMigration(cfg); err != nil {
		log.Printf("WARNING: Database migration skipped — %v", err)
		log.Println("Continuing with in-memory mocks")
	}

	userRepo := repository.NewMockUserRepository()
	proposalRepo := repository.NewMockProposalRepository()
	savedProposalRepo := repository.NewMockSavedProposalRepository(proposalRepo)
	sessionRepo := repository.NewMockSessionRepository()

	authService := service.NewAuthService(userRepo, sessionRepo, cfg.SecretKey)
	datosGovService := service.NewDatosGovService(cfg.IntegrationURL)
	proposalService := service.NewProposalService(proposalRepo)
	savedProposalService := service.NewSavedProposalService(savedProposalRepo)

	mux := router.NewRouter(authService, datosGovService, savedProposalService, proposalService)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting server on %s", addr)

	if _, err := os.Stat(".env"); err == nil {
		log.Println("Found .env file — ensure environment variables are set")
	}

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// runMigration attempts to connect to PostgreSQL and run schema migration.
// Returns nil on success or if the database is unreachable (graceful fallback).
func runMigration(cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := database.NewPool(ctx, cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	if err != nil {
		return fmt.Errorf("create pool: %w", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	sqlPath := "../database/public_calls_database.sql"
	if err := database.AutoMigrate(ctx, pool, sqlPath); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}

	return nil
}
