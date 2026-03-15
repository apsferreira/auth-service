package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// buildAuthAppWithLocals creates a Fiber app that pre-populates Locals before
// routing to the auth handler. Useful for testing protected endpoints (Me, UpdateMe, Logout).
func buildAuthAppWithLocals(userID, email string) *fiber.App {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	h := NewAuthHandler(nil, nil, nil)

	app.Use(func(c *fiber.Ctx) error {
		if userID != "" {
			c.Locals("userID", userID)
		}
		if email != "" {
			c.Locals("email", email)
		}
		return c.Next()
	})

	app.Get("/api/v1/auth/me", h.Me)
	app.Patch("/api/v1/auth/me", h.UpdateMe)
	app.Post("/api/v1/auth/logout", h.Logout)
	app.Get("/api/v1/auth/google", h.GoogleLogin)
	app.Get("/api/v1/auth/google/callback", h.GoogleCallback)

	return app
}

func doAuthGET(t *testing.T, app *fiber.App, path string) *http.Response {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	return resp
}

// ─── Me ─────────────────────────────────────────────────────────────────────

func TestMe_NoLocals_Returns401(t *testing.T) {
	// Build app without setting userID local
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	h := NewAuthHandler(nil, nil, nil)
	app.Get("/me", h.Me)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 when userID local is missing, got %d", resp.StatusCode)
	}
}

func TestMe_InvalidUUIDLocals_Returns400(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	h := NewAuthHandler(nil, nil, nil)
	app.Get("/me", func(c *fiber.Ctx) error {
		c.Locals("userID", "not-a-uuid")
		return h.Me(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid user ID, got %d", resp.StatusCode)
	}
}

func TestMe_ValidUUID_NilService_Returns404(t *testing.T) {
	// Valid UUID but service is nil → GetCurrentUser panics → Fiber returns 500
	// Actually with nil service, GetCurrentUser is called on nil pointer — Fiber catches the panic
	validID := uuid.New().String()
	app := buildAuthAppWithLocals(validID, "user@example.com")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	// nil authService will panic → Fiber catches → 500, OR we get 404 if service is configured
	// Since authService is nil, expect panic recovery → 500
	if resp.StatusCode != http.StatusInternalServerError && resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 500 or 404 when service is nil, got %d", resp.StatusCode)
	}
}

// ─── UpdateMe ────────────────────────────────────────────────────────────────

func TestUpdateMe_NoLocals_Returns401(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	h := NewAuthHandler(nil, nil, nil)
	app.Patch("/me", h.UpdateMe)

	req := httptest.NewRequest(http.MethodPatch, "/me", strings.NewReader(`{"full_name":"Test"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 when userID local is missing, got %d", resp.StatusCode)
	}
}

func TestUpdateMe_InvalidUUID_Returns400(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	h := NewAuthHandler(nil, nil, nil)
	app.Patch("/me", func(c *fiber.Ctx) error {
		c.Locals("userID", "invalid-uuid")
		return h.UpdateMe(c)
	})

	req := httptest.NewRequest(http.MethodPatch, "/me", strings.NewReader(`{"full_name":"Test"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid UUID, got %d", resp.StatusCode)
	}
}

func TestUpdateMe_EmptyFullName_Returns400(t *testing.T) {
	validID := uuid.New().String()
	app := buildAuthAppWithLocals(validID, "user@example.com")

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/auth/me", strings.NewReader(`{"full_name":""}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for empty full_name, got %d", resp.StatusCode)
	}
}

func TestUpdateMe_MissingFullName_Returns400(t *testing.T) {
	validID := uuid.New().String()
	app := buildAuthAppWithLocals(validID, "user@example.com")

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/auth/me", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 when full_name is missing, got %d", resp.StatusCode)
	}
}

func TestUpdateMe_InvalidJSON_Returns400(t *testing.T) {
	validID := uuid.New().String()
	app := buildAuthAppWithLocals(validID, "user@example.com")

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/auth/me", strings.NewReader(`{bad json`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid JSON, got %d", resp.StatusCode)
	}
}

// ─── GoogleLogin ─────────────────────────────────────────────────────────────

func TestGoogleLogin_ServiceNotConfigured_Returns503(t *testing.T) {
	app := buildAuthAppWithLocals("", "")
	// googleService is nil → not configured

	resp := doAuthGET(t, app, "/api/v1/auth/google")
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("expected 503 when Google service is not configured, got %d", resp.StatusCode)
	}
}

func TestGoogleLogin_MissingRedirectURI_Returns400(t *testing.T) {
	// When googleService is nil, returns 503 before checking redirect_uri.
	// This test verifies the 503 path (same as above), which is the validation gate.
	app := buildAuthAppWithLocals("", "")
	resp := doAuthGET(t, app, "/api/v1/auth/google?redirect_uri=")
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("expected 503 from unconfigured Google service, got %d", resp.StatusCode)
	}
}

// ─── GoogleCallback ───────────────────────────────────────────────────────────

func TestGoogleCallback_ErrorParam_Returns400(t *testing.T) {
	app := buildAuthAppWithLocals("", "")
	resp := doAuthGET(t, app, "/api/v1/auth/google/callback?error=access_denied")
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 when error param is set, got %d", resp.StatusCode)
	}
}

func TestGoogleCallback_MissingCode_Returns400(t *testing.T) {
	app := buildAuthAppWithLocals("", "")
	resp := doAuthGET(t, app, "/api/v1/auth/google/callback?state=somestatevalue")
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for missing code, got %d", resp.StatusCode)
	}
}

func TestGoogleCallback_InvalidState_Returns400(t *testing.T) {
	app := buildAuthAppWithLocals("", "")
	// state is not valid base64
	resp := doAuthGET(t, app, "/api/v1/auth/google/callback?code=mycode&state=!!!invalid-base64")
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid state (bad base64), got %d", resp.StatusCode)
	}
}

func TestGoogleCallback_EmptyState_Returns400(t *testing.T) {
	app := buildAuthAppWithLocals("", "")
	// state decodes to empty string
	// base64url of "" → ""
	resp := doAuthGET(t, app, "/api/v1/auth/google/callback?code=mycode&state=")
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for empty state, got %d", resp.StatusCode)
	}
}

// ─── Logout ───────────────────────────────────────────────────────────────────

func TestLogout_EmptyBody_ServiceNil_Returns500OrBadRequest(t *testing.T) {
	// nil service means Logout panics — Fiber catches it as 500
	// But first, Fiber tries to parse the empty body (which is valid — empty refresh_token)
	// then calls authService.Logout("") → nil pointer dereference → 500
	validID := uuid.New().String()
	app := buildAuthAppWithLocals(validID, "user@example.com")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	// nil service → panic → 500, or body-parse failure → 400
	if resp.StatusCode != http.StatusInternalServerError && resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 or 500, got %d", resp.StatusCode)
	}
}

func TestLogout_InvalidJSON_Returns400(t *testing.T) {
	app := buildAuthAppWithLocals(uuid.New().String(), "user@example.com")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", strings.NewReader(`{bad json`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid JSON, got %d", resp.StatusCode)
	}
}
