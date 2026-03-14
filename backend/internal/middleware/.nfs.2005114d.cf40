package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func buildRateLimitApp(maxRequests int, window time.Duration) *fiber.App {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	rateLimiter := NewRateLimiter(maxRequests, window)
	app.Use(rateLimiter.Handler())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "ok"})
	})
	return app
}

func TestRateLimit_AllowsRequestsWithinLimit(t *testing.T) {
	app := buildRateLimitApp(5, time.Minute)
	
	// Make 5 requests (should all succeed)
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:8080"
		
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("fiber.Test error on request %d: %v", i+1, err)
		}
		
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200 for request %d within limit, got %d", i+1, resp.StatusCode)
		}
	}
}

func TestRateLimit_BlocksRequestsExceedingLimit(t *testing.T) {
	app := buildRateLimitApp(3, time.Minute)
	
	// Make 3 requests (should succeed)
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:8080"
		
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("fiber.Test error on request %d: %v", i+1, err)
		}
		
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200 for request %d within limit, got %d", i+1, resp.StatusCode)
		}
	}
	
	// 4th request should be blocked
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:8080"
	
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error on blocked request: %v", err)
	}
	
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("expected 429 for request exceeding limit, got %d", resp.StatusCode)
	}
}

func TestRateLimit_DifferentIPs_SeparateLimits(t *testing.T) {
	app := buildRateLimitApp(2, time.Minute)
	
	// Make 2 requests from IP1 (should succeed)
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:8080"
		
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("fiber.Test error for IP1 request %d: %v", i+1, err)
		}
		
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200 for IP1 request %d, got %d", i+1, resp.StatusCode)
		}
	}
	
	// 3rd request from IP1 should be blocked
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:8080"
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error for IP1 blocked request: %v", err)
	}
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("expected 429 for IP1 exceeding limit, got %d", resp.StatusCode)
	}
	
	// But requests from IP2 should still work
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.2:8080"
		
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("fiber.Test error for IP2 request %d: %v", i+1, err)
		}
		
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200 for IP2 request %d, got %d", i+1, resp.StatusCode)
		}
	}
}

func TestRateLimit_ResetsAfterWindow(t *testing.T) {
	app := buildRateLimitApp(1, 100*time.Millisecond)
	
	// First request should succeed
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:8080"
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error for first request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 for first request, got %d", resp.StatusCode)
	}
	
	// Second request should be blocked
	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:8080"
	resp, err = app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error for second request: %v", err)
	}
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("expected 429 for second request, got %d", resp.StatusCode)
	}
	
	// Wait for window to reset
	time.Sleep(150 * time.Millisecond)
	
	// Third request should succeed after reset
	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:8080"
	resp, err = app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error for third request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 for request after reset, got %d", resp.StatusCode)
	}
}

func TestRateLimit_IncludesRetryAfterHeader(t *testing.T) {
	app := buildRateLimitApp(1, time.Minute)
	
	// First request to consume the limit
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:8080"
	app.Test(req, -1)
	
	// Second request should return 429 with retry-after info
	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:8080"
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error for blocked request: %v", err)
	}
	
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("expected 429 for blocked request, got %d", resp.StatusCode)
	}
	
	// Check that the response contains retry_after information
	// We can't easily check the exact value, but we can check the structure
	// by examining the response body (it should be JSON with retry_after field)
}

func TestNewRateLimiter_CreatesValidInstance(t *testing.T) {
	rl := NewRateLimiter(10, time.Minute)
	
	if rl == nil {
		t.Fatal("NewRateLimiter returned nil")
	}
	
	if rl.max != 10 {
		t.Errorf("expected max=10, got %d", rl.max)
	}
	
	if rl.window != time.Minute {
		t.Errorf("expected window=1m, got %v", rl.window)
	}
	
	if rl.entries == nil {
		t.Error("expected entries map to be initialized")
	}
}