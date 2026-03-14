package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// buildAdminApp mounts a Fiber app with AdminHandler routes using nil services.
// Safe for validation paths that return before calling any service.
func buildAdminApp() *fiber.App {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	h := NewAdminHandler(nil, nil)

	// Set up routes
	admin := app.Group("/api/v1/admin")
	
	// Services
	admin.Get("/services", h.ListServices)
	admin.Post("/services", h.CreateService)
	admin.Get("/services/:id", h.GetService)
	admin.Put("/services/:id", h.UpdateService)
	admin.Delete("/services/:id", h.DeleteService)
	
	// Service permissions
	admin.Get("/services/:id/permissions", h.ListServicePermissions)
	admin.Post("/services/:id/permissions", h.CreateServicePermission)
	
	// Permissions
	admin.Get("/permissions", h.ListAllPermissions)
	admin.Delete("/permissions/:id", h.DeletePermission)
	
	// Roles
	admin.Get("/roles", h.ListRoles)
	admin.Post("/roles", h.CreateRole)
	admin.Put("/roles/:id", h.UpdateRole)
	admin.Get("/roles/:id/permissions", h.GetRolePermissions)
	admin.Post("/roles/:id/permissions", h.SetRolePermissions)
	
	// Events
	admin.Get("/events", h.ListEvents)

	return app
}

// doAdminRequest executes a request with proper locals set for admin endpoints
func doAdminRequest(t *testing.T, app *fiber.App, method, path, body string, headers map[string]string) *http.Response {
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

// ─── Services ────────────────────────────────────────────────────────────────

func TestListServices_InvalidTenant_Returns400(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	h := NewAdminHandler(nil, nil)
	app.Get("/services", func(c *fiber.Ctx) error {
		c.Locals("tenantID", "invalid-uuid")
		return h.ListServices(c)
	})

	req := httptest.NewRequest("GET", "/services", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid tenant, got %d", resp.StatusCode)
	}
}

func TestCreateService_InvalidJSON_Returns400(t *testing.T) {
	app := buildAdminApp()
	resp := doAdminRequest(t, app, "POST", "/api/v1/admin/services", `{bad json`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid JSON, got %d", resp.StatusCode)
	}
}

func TestGetService_InvalidID_Returns400(t *testing.T) {
	app := buildAdminApp()
	resp := doAdminRequest(t, app, "GET", "/api/v1/admin/services/invalid-uuid", "", nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid service ID, got %d", resp.StatusCode)
	}
}

func TestUpdateService_InvalidID_Returns400(t *testing.T) {
	app := buildAdminApp()
	resp := doAdminRequest(t, app, "PUT", "/api/v1/admin/services/invalid-uuid", `{"name":"test"}`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid service ID, got %d", resp.StatusCode)
	}
}

func TestDeleteService_InvalidID_Returns400(t *testing.T) {
	app := buildAdminApp()
	resp := doAdminRequest(t, app, "DELETE", "/api/v1/admin/services/invalid-uuid", "", nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid service ID, got %d", resp.StatusCode)
	}
}

// ─── Service Permissions ─────────────────────────────────────────────────────

func TestListServicePermissions_InvalidServiceID_Returns400(t *testing.T) {
	app := buildAdminApp()
	resp := doAdminRequest(t, app, "GET", "/api/v1/admin/services/invalid-uuid/permissions", "", nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid service ID, got %d", resp.StatusCode)
	}
}

func TestCreateServicePermission_InvalidServiceID_Returns400(t *testing.T) {
	app := buildAdminApp()
	resp := doAdminRequest(t, app, "POST", "/api/v1/admin/services/invalid-uuid/permissions", `{"name":"test"}`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid service ID, got %d", resp.StatusCode)
	}
}

func TestCreateServicePermission_InvalidJSON_Returns400(t *testing.T) {
	serviceID := uuid.New().String()
	app := buildAdminApp()
	resp := doAdminRequest(t, app, "POST", "/api/v1/admin/services/"+serviceID+"/permissions", `{bad json`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid JSON, got %d", resp.StatusCode)
	}
}

// ─── Permissions ─────────────────────────────────────────────────────────────

func TestDeletePermission_InvalidID_Returns400(t *testing.T) {
	app := buildAdminApp()
	resp := doAdminRequest(t, app, "DELETE", "/api/v1/admin/permissions/invalid-uuid", "", nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid permission ID, got %d", resp.StatusCode)
	}
}

// ─── Roles ───────────────────────────────────────────────────────────────────

func TestListRoles_InvalidTenant_Returns400(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	h := NewAdminHandler(nil, nil)
	app.Get("/roles", func(c *fiber.Ctx) error {
		c.Locals("tenantID", "invalid-uuid")
		return h.ListRoles(c)
	})

	req := httptest.NewRequest("GET", "/roles", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid tenant, got %d", resp.StatusCode)
	}
}

func TestCreateRole_InvalidJSON_Returns400(t *testing.T) {
	app := buildAdminApp()
	resp := doAdminRequest(t, app, "POST", "/api/v1/admin/roles", `{bad json`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid JSON, got %d", resp.StatusCode)
	}
}

func TestUpdateRole_InvalidID_Returns400(t *testing.T) {
	app := buildAdminApp()
	resp := doAdminRequest(t, app, "PUT", "/api/v1/admin/roles/invalid-uuid", `{"name":"test"}`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid role ID, got %d", resp.StatusCode)
	}
}

func TestGetRolePermissions_InvalidID_Returns400(t *testing.T) {
	app := buildAdminApp()
	resp := doAdminRequest(t, app, "GET", "/api/v1/admin/roles/invalid-uuid/permissions", "", nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid role ID, got %d", resp.StatusCode)
	}
}

func TestSetRolePermissions_InvalidRoleID_Returns400(t *testing.T) {
	app := buildAdminApp()
	resp := doAdminRequest(t, app, "POST", "/api/v1/admin/roles/invalid-uuid/permissions", `{"permission_ids":[]}`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid role ID, got %d", resp.StatusCode)
	}
}

func TestSetRolePermissions_InvalidJSON_Returns400(t *testing.T) {
	roleID := uuid.New().String()
	app := buildAdminApp()
	resp := doAdminRequest(t, app, "POST", "/api/v1/admin/roles/"+roleID+"/permissions", `{bad json`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid JSON, got %d", resp.StatusCode)
	}
}

func TestSetRolePermissions_InvalidPermissionID_Returns400(t *testing.T) {
	roleID := uuid.New().String()
	app := buildAdminApp()
	resp := doAdminRequest(t, app, "POST", "/api/v1/admin/roles/"+roleID+"/permissions", `{"permission_ids":["invalid-uuid"]}`, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid permission ID, got %d", resp.StatusCode)
	}
}

// ─── Events ──────────────────────────────────────────────────────────────────

func TestListEvents_ValidQueryParams_CallsService(t *testing.T) {
	app := buildAdminApp()
	// This will return 500 since service is nil, but confirms validation passes
	resp := doAdminRequest(t, app, "GET", "/api/v1/admin/events?limit=10&offset=20&event_type=login&email=test@example.com", "", nil)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 when service is nil (validation passed), got %d", resp.StatusCode)
	}
}