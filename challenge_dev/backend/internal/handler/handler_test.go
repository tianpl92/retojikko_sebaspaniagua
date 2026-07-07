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
	savedRepo := repository.NewMockSavedProposalRepository()

	sessionRepo := repository.NewMockSessionRepository()
	authService := service.NewAuthService(userRepo, sessionRepo, "test-secret")
	datosGovService := service.NewDatosGovService("https://datos.gov.co/resource/p6dx-8zbt.json")
	savedProposalService := service.NewSavedProposalService(savedRepo)

	mux := http.NewServeMux()
	mux.Handle("/health", NewHealthHandler())
	mux.Handle("/login", NewAuthHandler(authService))
	mux.Handle("/user-create", NewCreateUserHandler(authService))

	userHandler := NewUserHandler(authService)
	mux.Handle("/user-info", AuthMiddleware(authService, userHandler))

	userModifyHandler := NewUserModifyHandler(authService)
	mux.Handle("/user-modify", AuthMiddleware(authService, userModifyHandler))

	mux.Handle("/public-proposals", AuthMiddleware(authService, NewProposalHandler(datosGovService)))

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
}

func TestCreateUserEndpoint_DuplicateDocument(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

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

	resp, err := http.Get(server.URL + "/public-proposals")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

func TestSavedProposalsEndpoint_Unauthorized(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

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

	// Save a proposal
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

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	// Should return success message
	if result["message"] != "Guardado satisfactoriamente" {
		t.Errorf("expected message 'Guardado satisfactoriamente', got %v", result["message"])
	}

	// List saved proposals
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

func TestUserInfoEndpoint_RequiresToken(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Post(server.URL+"/user-info", "application/json", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

func TestUserInfoEndpoint_Success(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	token := login(t, server.URL, "test@example.com", "password123")

	req, err := http.NewRequest("POST", server.URL+"/user-info", nil)
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

	var user map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if user["email"] != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %v", user["email"])
	}

	// Password must NOT be in response
	if _, hasPassword := user["password"]; hasPassword {
		t.Error("password should NOT be included in user-info response")
	}
}

func TestUserInfoEndpoint_InvalidToken(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	req, _ := http.NewRequest("POST", server.URL+"/user-info", nil)
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

func TestUserModifyEndpoint_RequiresToken(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Post(server.URL+"/user-modify", "application/json", strings.NewReader(`{"first_name":"New"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

func TestUserModifyEndpoint_Success(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	token := login(t, server.URL, "test@example.com", "password123")

	body := `{"first_name":"NewName","last_name":"NewLast","phone_number":"3001112233"}`
	req, err := http.NewRequest("POST", server.URL+"/user-modify", strings.NewReader(body))
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

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if result["message"] != "Informacion actualizada correctamente" {
		t.Errorf("expected success message, got %v", result["message"])
	}
}

func TestUserModifyEndpoint_NoFields(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	token := login(t, server.URL, "test@example.com", "password123")

	body := `{}`
	req, err := http.NewRequest("POST", server.URL+"/user-modify", strings.NewReader(body))
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

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestUserModifyEndpoint_InvalidToken(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	req, _ := http.NewRequest("POST", server.URL+"/user-modify", strings.NewReader(`{"first_name":"X"}`))
	req.Header.Set("Content-Type", "application/json")
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

func TestUserModifyEndpoint_RejectsEmail(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	token := login(t, server.URL, "test@example.com", "password123")

	// Only sending email (blocked field) — should fail with no allowed fields
	body := `{"email":"new@new.com"}`
	req, err := http.NewRequest("POST", server.URL+"/user-modify", strings.NewReader(body))
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

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}
