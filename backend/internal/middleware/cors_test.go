package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func buildCORSApp(allowedOrigins string) *fiber.App {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Use(CORS(allowedOrigins))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "ok"})
	})
	app.Post("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "created"})
	})
	return app
}

func TestCORS_AllowsAllOrigins_WhenWildcard(t *testing.T) {
	app := buildCORSApp("*")
	
	testCases := []struct {
		origin   string
		expected string
	}{
		{"http://localhost:3000", "http://localhost:3000"},
		{"https://example.com", "https://example.com"},
		{"http://malicious.com", "http://malicious.com"},
		{"", ""},
	}
	
	for _, tc := range testCases {
		t.Run(tc.origin, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			if tc.origin != "" {
				req.Header.Set("Origin", tc.origin)
			}
			
			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatalf("fiber.Test error: %v", err)
			}
			
			allowOrigin := resp.Header.Get("Access-Control-Allow-Origin")
			if allowOrigin != tc.expected {
				t.Errorf("expected Access-Control-Allow-Origin %q, got %q", tc.expected, allowOrigin)
			}
		})
	}
}

func TestCORS_AllowsSpecificOrigins_WhenListed(t *testing.T) {
	app := buildCORSApp("http://localhost:3000,https://example.com")
	
	testCases := []struct {
		origin   string
		expected string
	}{
		{"http://localhost:3000", "http://localhost:3000"},
		{"https://example.com", "https://example.com"},
		{"http://malicious.com", ""},
		{"http://localhost:3001", ""},
		{"", ""},
	}
	
	for _, tc := range testCases {
		t.Run(tc.origin, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			if tc.origin != "" {
				req.Header.Set("Origin", tc.origin)
			}
			
			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatalf("fiber.Test error: %v", err)
			}
			
			allowOrigin := resp.Header.Get("Access-Control-Allow-Origin")
			if allowOrigin != tc.expected {
				t.Errorf("expected Access-Control-Allow-Origin %q, got %q", tc.expected, allowOrigin)
			}
		})
	}
}

func TestCORS_SetsRequiredHeaders(t *testing.T) {
	app := buildCORSApp("*")
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	
	expectedHeaders := map[string]string{
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Headers":     "Content-Type, Authorization, Accept",
		"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, OPTIONS",
	}
	
	for header, expected := range expectedHeaders {
		actual := resp.Header.Get(header)
		if actual != expected {
			t.Errorf("expected %s %q, got %q", header, expected, actual)
		}
	}
}

func TestCORS_HandlesOPTIONS_Returns204(t *testing.T) {
	app := buildCORSApp("*")
	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected 204 for OPTIONS request, got %d", resp.StatusCode)
	}
}

func TestCORS_ContinuesForNormalRequests(t *testing.T) {
	app := buildCORSApp("*")
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 for normal request, got %d", resp.StatusCode)
	}
}

func TestCORS_TrimsWhitespaceFromOrigins(t *testing.T) {
	app := buildCORSApp(" http://localhost:3000 , https://example.com ")
	
	testCases := []struct {
		origin   string
		expected string
	}{
		{"http://localhost:3000", "http://localhost:3000"},
		{"https://example.com", "https://example.com"},
		{"http://other.com", ""},
	}
	
	for _, tc := range testCases {
		t.Run(tc.origin, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Origin", tc.origin)
			
			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatalf("fiber.Test error: %v", err)
			}
			
			allowOrigin := resp.Header.Get("Access-Control-Allow-Origin")
			if allowOrigin != tc.expected {
				t.Errorf("expected Access-Control-Allow-Origin %q, got %q", tc.expected, allowOrigin)
			}
		})
	}
}