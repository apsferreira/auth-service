package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// buildUserApp mounts a Fiber app with UserHandler routes using nil service.
// Safe for validation paths that return before calling any service.
func buildUserApp() *fiber.App {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	h := NewUserHandler(nil)

	users := app.Group("/api/v1/users")
	users.Get("/", h.List)
	users.Post("/", h.Create)
	users.Get("/:id", h.GetByID)
	users.Put("/:id", h.Update)
	users.Delete("/:id", h.Delete)

	return app
}

// doUserRequest executes a request with proper locals set for user endpoints
func doUserRequest(t *testing.T, app *fiber.App, method, path, body string, headers map[string]string) *http.Response {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	
	// Mock middleware context by setting a test handler that adds locals
	testApp := fiber.New(fiber.Config{DisableStartupMessage: true})
	testApp.Use(func(c *fiber.Ctx) error {
		c.Locals("tenantID", "550e8400-e29b-41d4-a716-446655440000")
		c.Locals("userID", "660e8400-e29b-41d4-a716-446655440000")
		return c.Next()
	})
	
	// Copy routes from original app
	for _, route := range app.GetRoutes() {
		switch route.Method {
		case "GET":
			testApp.Get(route.Path, route.Handlers...)
		case "POST":
			testApp.Post(route.Path, route.Handlers...)
		case "PUT":
			testApp.Put(route.Path, route.Handlers...)
		case "DELETE":
			testApp.Delete(route.Path, route.Handlers...)
		}
	}
	
	resp, err := testApp.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	return resp
}

// ─── List ────────────────────────────────────────────────────────────────────

func TestListUsers_MissingTenantID_Returns401(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	h := NewUserHandler(nil)
	app.Get("/users", h.List)

	req := httptest.NewRequest("GET", "/users", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 when tenantID is missing, got %d", resp.StatusCode)
	}
}

func TestListUsers_InvalidTenantID_Returns400(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	h := NewUserHandler(nil)
	app.Get("/users", func(c *fiber.Ctx) error {
		c.Locals("tenantID", "invalid-uuid")
		return h.List(c)
	})

	req := httptest.NewRequest("GET", "/users", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid tenant ID, got %d", resp.StatusCode)
	}
}

// ─── Create ──────────────────────────────────────────────────────────────────

func TestCreateUser_InvalidJSON_Returns400(t *testing.T) {
	app := buildUserApp()
	resp := doUserRequest(t, app, "POST", "/api/v1/users", `{bad json`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid JSON, got %d", resp.StatusCode)
	}
}

func TestCreateUser_EmptyEmail_Returns400(t *testing.T) {
	app := buildUserApp()
	resp := doUserRequest(t, app, "POST", "/api/v1/users", `{"email":""}`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for empty email, got %d", resp.StatusCode)
	}
}

func TestCreateUser_MissingEmail_Returns400(t *testing.T) {
	app := buildUserApp()
	resp := doUserRequest(t, app, "POST", "/api/v1/users", `{"full_name":"John Doe"}`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 when email is missing, got %d", resp.StatusCode)
	}
}

func TestCreateUser_ValidData_CallsService(t *testing.T) {
	app := buildUserApp()
	// This will return 500 since service is nil, but confirms validation passes
	resp := doUserRequest(t, app, "POST", "/api/v1/users", `{"email":"test@example.com","full_name":"John Doe"}`, nil)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 when service is nil (validation passed), got %d", resp.StatusCode)
	}
}

// ─── GetByID ─────────────────────────────────────────────────────────────────

func TestGetUserByID_InvalidID_Returns400(t *testing.T) {
	app := buildUserApp()
	resp := doUserRequest(t, app, "GET", "/api/v1/users/invalid-uuid", "", nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid user ID, got %d", resp.StatusCode)
	}
}

func TestGetUserByID_ValidID_CallsService(t *testing.T) {
	userID := uuid.New().String()
	app := buildUserApp()
	// This will return 500 since service is nil, but confirms validation passes
	resp := doUserRequest(t, app, "GET", "/api/v1/users/"+userID, "", nil)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 when service is nil (validation passed), got %d", resp.StatusCode)
	}
}

// ─── Update ──────────────────────────────────────────────────────────────────

func TestUpdateUser_InvalidID_Returns400(t *testing.T) {
	app := buildUserApp()
	resp := doUserRequest(t, app, "PUT", "/api/v1/users/invalid-uuid", `{"full_name":"Updated Name"}`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid user ID, got %d", resp.StatusCode)
	}
}

func TestUpdateUser_InvalidJSON_Returns400(t *testing.T) {
	userID := uuid.New().String()
	app := buildUserApp()
	resp := doUserRequest(t, app, "PUT", "/api/v1/users/"+userID, `{bad json`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid JSON, got %d", resp.StatusCode)
	}
}

func TestUpdateUser_ValidData_CallsService(t *testing.T) {
	userID := uuid.New().String()
	app := buildUserApp()
	// This will return 500 since service is nil, but confirms validation passes
	resp := doUserRequest(t, app, "PUT", "/api/v1/users/"+userID, `{"full_name":"Updated Name"}`, nil)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 when service is nil (validation passed), got %d", resp.StatusCode)
	}
}

// ─── Delete ──────────────────────────────────────────────────────────────────

func TestDeleteUser_InvalidID_Returns400(t *testing.T) {
	app := buildUserApp()
	resp := doUserRequest(t, app, "DELETE", "/api/v1/users/invalid-uuid", "", nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid user ID, got %d", resp.StatusCode)
	}
}

func TestDeleteUser_ValidID_CallsService(t *testing.T) {
	userID := uuid.New().String()
	app := buildUserApp()
	// This will return 500 since service is nil, but confirms validation passes
	resp := doUserRequest(t, app, "DELETE", "/api/v1/users/"+userID, "", nil)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 when service is nil (validation passed), got %d", resp.StatusCode)
	}
}