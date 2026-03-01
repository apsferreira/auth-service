package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

const handlerTestServiceToken = "test-service-token-xyz"

// buildAuthApp monta um app Fiber com as rotas do AuthHandler usando nil para os serviços.
// Isso é seguro para todos os caminhos de validação (400/401), que retornam ANTES de chamar qualquer serviço.
func buildAuthApp() *fiber.App {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	h := NewAuthHandler(nil, nil)
	app.Post("/api/v1/auth/request-otp", h.RequestOTP)
	app.Post("/api/v1/auth/verify-otp", h.VerifyOTP)
	app.Post("/api/v1/auth/refresh", h.Refresh)
	app.Post("/api/v1/auth/admin-login", h.AdminLogin)
	app.Post("/api/v1/auth/validate", h.Validate)
	app.Post("/api/v1/auth/provision-user", h.ProvisionUser(handlerTestServiceToken))
	return app
}

// doJSON executa uma request POST JSON no app de teste.
func doJSON(t *testing.T, app *fiber.App, path, body string, headers map[string]string) *http.Response {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	return resp
}

// ─── RequestOTP ──────────────────────────────────────────────────────────────

func TestRequestOTP_EmptyEmail_Returns400(t *testing.T) {
	app := buildAuthApp()
	resp := doJSON(t, app, "/api/v1/auth/request-otp", `{"email":""}`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for empty email, got %d", resp.StatusCode)
	}
}

func TestRequestOTP_MissingEmailField_Returns400(t *testing.T) {
	app := buildAuthApp()
	resp := doJSON(t, app, "/api/v1/auth/request-otp", `{}`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 when email field is missing, got %d", resp.StatusCode)
	}
}

func TestRequestOTP_WhitespaceOnlyEmail_Returns400(t *testing.T) {
	app := buildAuthApp()
	resp := doJSON(t, app, "/api/v1/auth/request-otp", `{"email":"   "}`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for whitespace-only email, got %d", resp.StatusCode)
	}
}

func TestRequestOTP_InvalidJSON_Returns400(t *testing.T) {
	app := buildAuthApp()
	resp := doJSON(t, app, "/api/v1/auth/request-otp", `{bad json`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid JSON, got %d", resp.StatusCode)
	}
}

// ─── VerifyOTP ───────────────────────────────────────────────────────────────

func TestVerifyOTP_EmptyEmail_Returns400(t *testing.T) {
	app := buildAuthApp()
	resp := doJSON(t, app, "/api/v1/auth/verify-otp", `{"email":"","code":"123456"}`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for empty email, got %d", resp.StatusCode)
	}
}

func TestVerifyOTP_EmptyCode_Returns400(t *testing.T) {
	app := buildAuthApp()
	resp := doJSON(t, app, "/api/v1/auth/verify-otp", `{"email":"a@b.com","code":""}`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for empty code, got %d", resp.StatusCode)
	}
}

func TestVerifyOTP_CodeTooShort_Returns400(t *testing.T) {
	app := buildAuthApp()
	resp := doJSON(t, app, "/api/v1/auth/verify-otp", `{"email":"a@b.com","code":"12345"}`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for 5-digit code, got %d", resp.StatusCode)
	}
}

func TestVerifyOTP_CodeTooLong_Returns400(t *testing.T) {
	app := buildAuthApp()
	resp := doJSON(t, app, "/api/v1/auth/verify-otp", `{"email":"a@b.com","code":"1234567"}`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for 7-digit code, got %d", resp.StatusCode)
	}
}

func TestVerifyOTP_InvalidJSON_Returns400(t *testing.T) {
	app := buildAuthApp()
	resp := doJSON(t, app, "/api/v1/auth/verify-otp", `{bad json`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid JSON, got %d", resp.StatusCode)
	}
}

// ─── Refresh ─────────────────────────────────────────────────────────────────

func TestRefresh_EmptyToken_Returns400(t *testing.T) {
	app := buildAuthApp()
	resp := doJSON(t, app, "/api/v1/auth/refresh", `{"refresh_token":""}`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for empty refresh_token, got %d", resp.StatusCode)
	}
}

func TestRefresh_MissingToken_Returns400(t *testing.T) {
	app := buildAuthApp()
	resp := doJSON(t, app, "/api/v1/auth/refresh", `{}`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 when refresh_token is missing, got %d", resp.StatusCode)
	}
}

func TestRefresh_InvalidJSON_Returns400(t *testing.T) {
	app := buildAuthApp()
	resp := doJSON(t, app, "/api/v1/auth/refresh", `{bad json`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid JSON, got %d", resp.StatusCode)
	}
}

// ─── AdminLogin ──────────────────────────────────────────────────────────────

func TestAdminLogin_EmptyIdentifier_Returns400(t *testing.T) {
	app := buildAuthApp()
	resp := doJSON(t, app, "/api/v1/auth/admin-login", `{"identifier":"","password":"secret"}`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for empty identifier, got %d", resp.StatusCode)
	}
}

func TestAdminLogin_EmptyPassword_Returns400(t *testing.T) {
	app := buildAuthApp()
	resp := doJSON(t, app, "/api/v1/auth/admin-login", `{"identifier":"admin","password":""}`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for empty password, got %d", resp.StatusCode)
	}
}

func TestAdminLogin_BothEmpty_Returns400(t *testing.T) {
	app := buildAuthApp()
	resp := doJSON(t, app, "/api/v1/auth/admin-login", `{"identifier":"","password":""}`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 when both identifier and password are empty, got %d", resp.StatusCode)
	}
}

func TestAdminLogin_InvalidJSON_Returns400(t *testing.T) {
	app := buildAuthApp()
	resp := doJSON(t, app, "/api/v1/auth/admin-login", `{bad json`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid JSON, got %d", resp.StatusCode)
	}
}

// ─── Validate ────────────────────────────────────────────────────────────────

func TestValidate_EmptyToken_Returns400(t *testing.T) {
	app := buildAuthApp()
	resp := doJSON(t, app, "/api/v1/auth/validate", `{"token":""}`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for empty token, got %d", resp.StatusCode)
	}
}

func TestValidate_MissingToken_Returns400(t *testing.T) {
	app := buildAuthApp()
	resp := doJSON(t, app, "/api/v1/auth/validate", `{}`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 when token field is missing, got %d", resp.StatusCode)
	}
}

func TestValidate_InvalidJSON_Returns400(t *testing.T) {
	app := buildAuthApp()
	resp := doJSON(t, app, "/api/v1/auth/validate", `{bad json`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid JSON, got %d", resp.StatusCode)
	}
}

// ─── ProvisionUser ───────────────────────────────────────────────────────────

func TestProvisionUser_WrongServiceToken_Returns401(t *testing.T) {
	app := buildAuthApp()
	resp := doJSON(t, app, "/api/v1/auth/provision-user", `{"email":"a@b.com"}`,
		map[string]string{"X-Service-Token": "wrong-token"})
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 for wrong service token, got %d", resp.StatusCode)
	}
}

func TestProvisionUser_MissingServiceToken_Returns401(t *testing.T) {
	app := buildAuthApp()
	resp := doJSON(t, app, "/api/v1/auth/provision-user", `{"email":"a@b.com"}`, nil)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 when X-Service-Token is absent, got %d", resp.StatusCode)
	}
}

// ─── sourceFromOrigin ────────────────────────────────────────────────────────

func TestSourceFromOrigin(t *testing.T) {
	cases := []struct {
		origin   string
		expected string
	}{
		{"", "direct"},
		{"http://localhost:3003", "auth-panel"},
		{"http://localhost:3001", "my-library"},
		{"http://localhost:3002", "focus-hub"},
		{"http://localhost:3005", "iit-agents"},
		{"http://localhost:5173", "dev-server"},
		{"https://app.example.com", "https://app.example.com"},
	}
	for _, tc := range cases {
		got := sourceFromOrigin(tc.origin)
		if got != tc.expected {
			t.Errorf("sourceFromOrigin(%q) = %q, want %q", tc.origin, got, tc.expected)
		}
	}
}
