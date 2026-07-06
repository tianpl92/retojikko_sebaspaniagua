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

	sessionRepo := repository.NewMockSessionRepository()
	authService := service.NewAuthService(userRepo, sessionRepo, "test-secret")
	proposalService := service.NewProposalService(proposalRepo)
	savedProposalService := service.NewSavedProposalService(savedRepo)

	mux := http.NewServeMux()
	mux.Handle("/health", NewHealthHandler())
	mux.Handle("/login", NewAuthHandler(authService))
	mux.Handle("/user-create", NewCreateUserHandler(authService))

	userHandler := NewUserHandler(authService)
	mux.Handle("/user-info", AuthMiddleware(authService, userHandler))
	mux.Handle("/user-modify", AuthMiddleware(authService, userHandler))

	mux.Handle("/public-proposals", AuthMiddleware(authService, NewProposalHandler(proposalService)))

	savedHandler := NewSavedProposalHandler(savedProposalService)
	mux.Handle("/saved_proposals", AuthMiddleware(authService, savedHandler))
	mux.Handle("/saved-proposals", AuthMiddleware(authService, savedHandler))

	return httptest.NewServer(mux)
}

func login(t *testing.T, serverURL, email, password string) string {
	t.Helper()
	body := `{"email":"` + email + `","password":"` + password + `"}`
	resp, err := http.Post(serverURL+"/login", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("login request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login failed with status %d", resp.StatusCode)
	}
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode login response: %v", err)
	}
	token, ok := result["token"].(string)
	if !ok || token == "" {
		t.Fatal("expected a valid token")
	}
	return token
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

func TestCreateUserEndpoint_Success(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	body := `{"id":"doc123","first_name":"New","last_name":"User","email":"new@example.com","password":"pass123","gender":"M","phone_number":"555-1234","status":"active"}`
	resp, err := http.Post(server.URL+"/user-create", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result["id"] != "doc123" {
		t.Errorf("expected id 'doc123', got %v", result["id"])
	}
	if result["email"] != "new@example.com" {
		t.Errorf("expected email 'new@example.com', got %v", result["email"])
	}
	if result["first_name"] != "New" {
		t.Errorf("expected first_name 'New', got %v", result["first_name"])
	}
}

func TestCreateUserEndpoint_DuplicateDocument(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// The default test user has ID "1"
	body := `{"id":"1","first_name":"Dup","last_name":"User","email":"dup@example.com","password":"pass123"}`
	resp, err := http.Post(server.URL+"/user-create", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusConflict {
		t.Errorf("expected 409, got %d", resp.StatusCode)
	}
}

func TestCreateUserEndpoint_DuplicateEmail(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	body := `{"id":"doc999","first_name":"Dup","last_name":"Email","email":"test@example.com","password":"pass123"}`
	resp, err := http.Post(server.URL+"/user-create", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusConflict {
		t.Errorf("expected 409, got %d", resp.StatusCode)
	}
}

func TestCreateUserEndpoint_MissingFields(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	tests := []struct {
		name string
		body string
	}{
		{"missing id", `{"first_name":"N","last_name":"U","email":"n@e.com","password":"p"}`},
		{"missing first_name", `{"id":"x","last_name":"U","email":"n@e.com","password":"p"}`},
		{"missing last_name", `{"id":"x","first_name":"N","email":"n@e.com","password":"p"}`},
		{"missing email", `{"id":"x","first_name":"N","last_name":"U","password":"p"}`},
		{"missing password", `{"id":"x","first_name":"N","last_name":"U","email":"n@e.com"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Post(server.URL+"/user-create", "application/json", strings.NewReader(tt.body))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusBadRequest {
				t.Errorf("expected 400, got %d", resp.StatusCode)
			}
		})
	}
}

func TestPublicProposalsEndpoint_RequiresToken(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Without token should return 401
	resp, err := http.Get(server.URL + "/public-proposals")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

func TestPublicProposalsEndpoint_WithValidToken(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	token := login(t, server.URL, "test@example.com", "password123")

	req, err := http.NewRequest("GET", server.URL+"/public-proposals", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
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
	_ = result
}

func TestPublicProposalsEndpoint_WithFilters(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	token := login(t, server.URL, "test@example.com", "password123")

	req, err := http.NewRequest("GET", server.URL+"/public-proposals?query=test&fase=fase1&entidad=ent1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
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
	resp, err := http.Get(server.URL + "/saved_proposals")
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

	token := login(t, server.URL, "test@example.com", "password123")

	// Save a proposal - use /saved-proposals for POST
	saveBody := `{"public_call_id":"42"}`
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

	// List saved proposals - use /saved_proposals for GET
	req, err = http.NewRequest("GET", server.URL+"/saved_proposals", nil)
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

	req, _ := http.NewRequest("GET", server.URL+"/saved_proposals", nil)
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
