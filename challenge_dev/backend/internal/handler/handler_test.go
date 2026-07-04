package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/repository"
	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/service"
)

func setupTestServer() *httptest.Server {
	userRepo := repository.NewMockUserRepository()
	proposalRepo := repository.NewMockProposalRepository()
	savedRepo := repository.NewMockSavedProposalRepository()

	authService := service.NewAuthService(userRepo, "test-secret")
	proposalService := service.NewProposalService(proposalRepo)
	savedProposalService := service.NewSavedProposalService(savedRepo)

	mux := http.NewServeMux()
	mux.Handle("/health", NewHealthHandler())
	mux.Handle("/login", NewAuthHandler(authService))
	mux.Handle("/public-proposals", NewProposalHandler(proposalService))
	mux.Handle("/saved-proposals", NewSavedProposalHandler(savedProposalService, authService))

	return httptest.NewServer(mux)
}

func TestHealthEndpoint(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result["status"] != "ok" {
		t.Errorf("expected status 'ok', got %q", result["status"])
	}
}

func TestLoginEndpoint_Success(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	body := `{"email":"test@example.com","password":"password123"}`
	resp, err := http.Post(server.URL+"/login", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if _, ok := result["token"]; !ok {
		t.Error("expected 'token' in response")
	}
	if _, ok := result["user"]; !ok {
		t.Error("expected 'user' in response")
	}
}

func TestLoginEndpoint_InvalidCredentials(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	body := `{"email":"test@example.com","password":"wrongpassword"}`
	resp, err := http.Post(server.URL+"/login", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

func TestLoginEndpoint_BadRequest(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	body := `not-json`
	resp, err := http.Post(server.URL+"/login", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestPublicProposalsEndpoint(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/public-proposals")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	_ = result // empty array is valid
}

func TestPublicProposalsEndpoint_WithFilters(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/public-proposals?query=test&category=cat1&fase=fase1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestSavedProposalsEndpoint_Unauthorized(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// GET without auth header
	resp, err := http.Get(server.URL + "/saved-proposals")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

func TestSavedProposalsEndpoint_CreateAndList(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// First login to get a token
	loginBody := `{"email":"test@example.com","password":"password123"}`
	loginResp, err := http.Post(server.URL+"/login", "application/json", strings.NewReader(loginBody))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var loginResult map[string]interface{}
	json.NewDecoder(loginResp.Body).Decode(&loginResult)
	loginResp.Body.Close()

	token, ok := loginResult["token"].(string)
	if !ok || token == "" {
		t.Fatal("expected a valid token")
	}

	// Save a proposal
	saveBody := `{"public_call_id":42}`
	req, err := http.NewRequest("POST", server.URL+"/saved-proposals", strings.NewReader(saveBody))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}

	var savedResult map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&savedResult); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if savedResult["public_call_id"] != float64(42) {
		t.Errorf("expected public_call_id 42, got %v", savedResult["public_call_id"])
	}

	// List saved proposals
	req, err = http.NewRequest("GET", server.URL+"/saved-proposals", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp2.StatusCode)
	}

	var listResult []interface{}
	if err := json.NewDecoder(resp2.Body).Decode(&listResult); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if len(listResult) != 1 {
		t.Errorf("expected 1 saved proposal, got %d", len(listResult))
	}
}

func TestSavedProposalsEndpoint_InvalidToken(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL+"/saved-proposals", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}
