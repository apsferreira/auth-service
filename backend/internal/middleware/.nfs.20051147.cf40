package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	jwtpkg "github.com/apsferreira/auth-service/backend/internal/pkg/jwt"
)

func buildAuthApp(jwtService *jwtpkg.JWTService) *fiber.App {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	
	// Protected route
	app.Get("/protected", Auth(jwtService), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"userID":      c.Locals("userID"),
			"tenantID":    c.Locals("tenantID"),
			"email":       c.Locals("email"),
			"roles":       c.Locals("roles"),
			"permissions": c.Locals("permissions"),
		})
	})
	
	// Role-protected route
	app.Get("/admin-only", Auth(jwtService), RequireRole("admin"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "admin access granted"})
	})
	
	// Permission-protected routes
	app.Get("/books-read", Auth(jwtService), RequirePermission("books.read"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "books.read access granted"})
	})
	
	app.Get("/library-books", Auth(jwtService), RequireServicePermission("my-library", "books.create"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "my-library books.create access granted"})
	})
	
	return app
}

func TestAuth_MissingAuthHeader_Returns401(t *testing.T) {
	app := buildAuthApp(nil) // JWT service not needed for this test
	req := httptest.NewRequest("GET", "/protected", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 for missing auth header, got %d", resp.StatusCode)
	}
}

func TestAuth_EmptyAuthHeader_Returns401(t *testing.T) {
	app := buildAuthApp(nil) // JWT service not needed for this test
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 for empty auth header, got %d", resp.StatusCode)
	}
}

func TestAuth_InvalidAuthHeaderFormat_Returns401(t *testing.T) {
	testCases := []string{
		"InvalidToken",
		"Bearer",
		"Basic dGVzdA==",
		"Bearer token with spaces",
	}
	
	app := buildAuthApp(nil) // JWT service not needed for this test
	
	for _, authHeader := range testCases {
		t.Run(authHeader, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set("Authorization", authHeader)
			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatalf("fiber.Test error: %v", err)
			}
			
			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("expected 401 for invalid auth header format %q, got %d", authHeader, resp.StatusCode)
			}
		})
	}
}

func TestAuth_ValidBearerFormat_CallsJWTService(t *testing.T) {
	// This test confirms that with valid "Bearer token" format, 
	// the JWT service validation is called (and fails since service is nil)
	app := buildAuthApp(nil)
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid-looking-token")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 when JWT service is nil (format validation passed), got %d", resp.StatusCode)
	}
}

func TestRequireRole_MissingRoles_Returns403(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/admin", RequireRole("admin"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/admin", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 when roles are missing, got %d", resp.StatusCode)
	}
}

func TestRequireRole_WrongRoles_Returns403(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/admin", func(c *fiber.Ctx) error {
		c.Locals("roles", []string{"user", "moderator"})
		return RequireRole("admin")(c)
	}, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/admin", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 when user lacks required role, got %d", resp.StatusCode)
	}
}

func TestRequireRole_HasRequiredRole_Continues(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/admin", func(c *fiber.Ctx) error {
		c.Locals("roles", []string{"user", "admin", "moderator"})
		return RequireRole("admin")(c)
	}, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/admin", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 when user has required role, got %d", resp.StatusCode)
	}
}

func TestRequirePermission_MissingPermissions_Returns403(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/books", RequirePermission("books.read"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/books", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 when permissions are missing, got %d", resp.StatusCode)
	}
}

func TestRequirePermission_WrongPermissions_Returns403(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/books", func(c *fiber.Ctx) error {
		c.Locals("permissions", map[string][]string{
			"my-library": {"books.create", "books.update"},
			"global":     {"users.read"},
		})
		return RequirePermission("books.delete")(c)
	}, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/books", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 when user lacks required permission, got %d", resp.StatusCode)
	}
}

func TestRequirePermission_HasRequiredPermission_Continues(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/books", func(c *fiber.Ctx) error {
		c.Locals("permissions", map[string][]string{
			"my-library": {"books.read", "books.create"},
			"global":     {"users.read"},
		})
		return RequirePermission("books.read")(c)
	}, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/books", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 when user has required permission, got %d", resp.StatusCode)
	}
}

func TestRequireServicePermission_MissingPermissions_Returns403(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/books", RequireServicePermission("my-library", "books.read"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/books", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 when permissions are missing, got %d", resp.StatusCode)
	}
}

func TestRequireServicePermission_WrongService_Returns403(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/books", func(c *fiber.Ctx) error {
		c.Locals("permissions", map[string][]string{
			"other-service": {"books.read"},
			"global":        {"users.read"},
		})
		return RequireServicePermission("my-library", "books.read")(c)
	}, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/books", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 when service permissions are missing, got %d", resp.StatusCode)
	}
}

func TestRequireServicePermission_WrongPermission_Returns403(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/books", func(c *fiber.Ctx) error {
		c.Locals("permissions", map[string][]string{
			"my-library": {"books.create", "books.update"},
			"global":     {"users.read"},
		})
		return RequireServicePermission("my-library", "books.delete")(c)
	}, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/books", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 when user lacks required service permission, got %d", resp.StatusCode)
	}
}

func TestRequireServicePermission_HasRequiredPermission_Continues(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/books", func(c *fiber.Ctx) error {
		c.Locals("permissions", map[string][]string{
			"my-library": {"books.read", "books.create"},
			"global":     {"users.read"},
		})
		return RequireServicePermission("my-library", "books.read")(c)
	}, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/books", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 when user has required service permission, got %d", resp.StatusCode)
	}
}