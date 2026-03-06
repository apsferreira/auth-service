package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func buildHealthApp() *fiber.App {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	h := NewHealthHandler()
	app.Get("/health", h.Health)
	return app
}

func TestHealth_Returns200WithValidResponse(t *testing.T) {
	app := buildHealthApp()
	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	// Check response body structure
	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	// Verify expected fields
	expectedFields := map[string]string{
		"status":  "ok",
		"service": "auth-service",
		"version": "1.0.0",
	}

	for field, expectedValue := range expectedFields {
		value, exists := body[field]
		if !exists {
			t.Errorf("expected field %q to exist in response", field)
			continue
		}
		
		if value != expectedValue {
			t.Errorf("expected %q to be %q, got %q", field, expectedValue, value)
		}
	}
}

func TestHealth_ContentTypeIsJSON(t *testing.T) {
	app := buildHealthApp()
	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got %q", contentType)
	}
}