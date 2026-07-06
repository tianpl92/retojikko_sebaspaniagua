package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestUserInfoEndpoint_RequiresToken(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/user-info")
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

	req, err := http.NewRequest("GET", server.URL+"/user-info", nil)
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

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result["email"] != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %v", result["email"])
	}
	if result["first_name"] != "Test" {
		t.Errorf("expected first_name 'Test', got %v", result["first_name"])
	}
	if result["id"] != "1" {
		t.Errorf("expected id '1', got %v", result["id"])
	}
	// Password should not be in response
	if _, ok := result["password"]; ok {
		t.Error("password should not be in response")
	}
}

func TestUserInfoEndpoint_InvalidToken(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL+"/user-info", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	if result["error"] != "Credentials invalid" {
		t.Errorf("expected error 'Credentials invalid', got %q", result["error"])
	}
}

func TestUserModifyEndpoint_RequiresToken(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Post(server.URL+"/user-modify", "application/json", strings.NewReader(`{"first_name":"Updated"}`))
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

	// Modify first_name and last_name
	body := `{"first_name":"UpdatedFirst","last_name":"UpdatedLast","phone_number":"555-9999"}`
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
		t.Fatalf("failed to decode response: %v", err)
	}

	if result["first_name"] != "UpdatedFirst" {
		t.Errorf("expected first_name 'UpdatedFirst', got %v", result["first_name"])
	}
	if result["last_name"] != "UpdatedLast" {
		t.Errorf("expected last_name 'UpdatedLast', got %v", result["last_name"])
	}
	if result["phone_number"] != "555-9999" {
		t.Errorf("expected phone_number '555-9999', got %v", result["phone_number"])
	}
	// Unchanged fields should still be present
	if result["email"] != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %v", result["email"])
	}
}

func TestUserModifyEndpoint_InvalidToken(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	body := `{"first_name":"Hacker"}`
	req, _ := http.NewRequest("POST", server.URL+"/user-modify", strings.NewReader(body))
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

func TestUserModifyEndpoint_BadRequest(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	token := login(t, server.URL, "test@example.com", "password123")

	req, _ := http.NewRequest("POST", server.URL+"/user-modify", strings.NewReader(`not-json`))
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
