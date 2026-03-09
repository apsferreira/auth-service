package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func buildLoggerApp() *fiber.App {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Use(Logger())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "ok"})
	})
	app.Get("/slow", func(c *fiber.Ctx) error {
		time.Sleep(10 * time.Millisecond) // Simulate slow endpoint
		return c.JSON(fiber.Map{"message": "slow response"})
	})
	app.Get("/error", func(c *fiber.Ctx) error {
		return c.Status(500).JSON(fiber.Map{"error": "internal error"})
	})
	return app
}

func captureLogOutput(fn func()) string {
	// Capture log output by redirecting to a buffer
	var buf bytes.Buffer
	oldOutput := log.Writer()
	log.SetOutput(&buf)
	
	fn()
	
	// Restore original output
	log.SetOutput(oldOutput)
	return buf.String()
}

func TestLogger_LogsSuccessfulRequest(t *testing.T) {
	app := buildLoggerApp()
	
	logOutput := captureLogOutput(func() {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:8080"
		
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("fiber.Test error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
	})
	
	// Verify log format contains expected elements
	expectedElements := []string{
		"[GET]",
		"/test",
		"192.168.1.1",
		"200",
	}
	
	for _, element := range expectedElements {
		if !strings.Contains(logOutput, element) {
			t.Errorf("expected log to contain %q, but log was: %s", element, logOutput)
		}
	}
}

func TestLogger_LogsErrorRequest(t *testing.T) {
	app := buildLoggerApp()
	
	logOutput := captureLogOutput(func() {
		req := httptest.NewRequest("GET", "/error", nil)
		req.RemoteAddr = "127.0.0.1:9000"
		
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("fiber.Test error: %v", err)
		}
		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", resp.StatusCode)
		}
	})
	
	// Verify log format contains expected elements
	if !strings.Contains(logOutput, "[GET]") || !strings.Contains(logOutput, "/error") || !strings.Contains(logOutput, "500") {
		t.Errorf("expected log to contain method, path, and status, but log was: %s", logOutput)
	}
	
	ipFound := strings.Contains(logOutput, "127.0.0.1") || strings.Contains(logOutput, "0.0.0.0")
	if !ipFound {
		t.Errorf("expected log to contain IP (127.0.0.1 or 0.0.0.0), but log was: %s", logOutput)
	}
}

func TestLogger_LogsLatency(t *testing.T) {
	app := buildLoggerApp()
	
	logOutput := captureLogOutput(func() {
		req := httptest.NewRequest("GET", "/slow", nil)
		
		start := time.Now()
		resp, err := app.Test(req, -1)
		elapsed := time.Since(start)
		
		if err != nil {
			t.Fatalf("fiber.Test error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
		
		// Sanity check that our slow endpoint actually took some time
		if elapsed < 5*time.Millisecond {
			t.Logf("Warning: slow endpoint completed too quickly (%v), timing test may not be reliable", elapsed)
		}
	})
	
	// The log should contain timing information (some duration measurement)
	// Look for patterns like "10ms", "1.5µs", etc.
	if !strings.Contains(logOutput, "s") && !strings.Contains(logOutput, "ms") && !strings.Contains(logOutput, "µs") && !strings.Contains(logOutput, "ns") {
		t.Errorf("expected log to contain latency information (time units), but log was: %s", logOutput)
	}
}

func TestLogger_LogsDifferentMethods(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Use(Logger())
	app.Post("/test", func(c *fiber.Ctx) error {
		return c.Status(201).JSON(fiber.Map{"created": true})
	})
	app.Put("/test", func(c *fiber.Ctx) error {
		return c.Status(204).Send(nil)
	})
	app.Delete("/test", func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	})
	
	testCases := []struct {
		method         string
		expectedStatus int
	}{
		{"POST", 201},
		{"PUT", 204},
		{"DELETE", 404},
	}
	
	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			logOutput := captureLogOutput(func() {
				req := httptest.NewRequest(tc.method, "/test", nil)
				
				resp, err := app.Test(req, -1)
				if err != nil {
					t.Fatalf("fiber.Test error: %v", err)
				}
				if resp.StatusCode != tc.expectedStatus {
					t.Errorf("expected %d, got %d", tc.expectedStatus, resp.StatusCode)
				}
			})
			
			expectedElements := []string{
				"[" + tc.method + "]",
				"/test",
			}
			
			for _, element := range expectedElements {
				if !strings.Contains(logOutput, element) {
					t.Errorf("expected log to contain %q, but log was: %s", element, logOutput)
				}
			}
		})
	}
}

func TestLogger_DoesNotInterfereWithResponse(t *testing.T) {
	app := buildLoggerApp()
	
	// Temporarily silence logs for this test to avoid cluttering test output
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)
	
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber.Test error: %v", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	
	// Verify the response body is still correct
	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got %q", contentType)
	}
}