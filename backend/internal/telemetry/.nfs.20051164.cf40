package telemetry

import (
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestInit_NoEndpoint_ReturnsNoopShutdown(t *testing.T) {
	// Ensure OTEL_EXPORTER_OTLP_ENDPOINT is not set
	oldEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	os.Unsetenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	defer func() {
		if oldEndpoint != "" {
			os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", oldEndpoint)
		}
	}()

	shutdown := Init("test-service")
	
	if shutdown == nil {
		t.Fatal("expected non-nil shutdown function")
	}
	
	// Should not panic when called
	shutdown()
}

func TestInit_WithServiceName_UsesEnvVariable(t *testing.T) {
	// Set service name env var
	oldServiceName := os.Getenv("OTEL_SERVICE_NAME")
	os.Setenv("OTEL_SERVICE_NAME", "custom-service-name")
	defer func() {
		if oldServiceName != "" {
			os.Setenv("OTEL_SERVICE_NAME", oldServiceName)
		} else {
			os.Unsetenv("OTEL_SERVICE_NAME")
		}
	}()

	// Ensure no endpoint (to avoid actual OTLP setup)
	oldEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	os.Unsetenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	defer func() {
		if oldEndpoint != "" {
			os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", oldEndpoint)
		}
	}()

	shutdown := Init("default-service")
	defer shutdown()

	// This test mainly verifies the function doesn't panic and returns a shutdown function
	// The actual service name usage would require more complex testing with mocked exporters
}

func TestInit_DefaultServiceName_WhenEnvNotSet(t *testing.T) {
	// Ensure OTEL_SERVICE_NAME is not set
	oldServiceName := os.Getenv("OTEL_SERVICE_NAME")
	os.Unsetenv("OTEL_SERVICE_NAME")
	defer func() {
		if oldServiceName != "" {
			os.Setenv("OTEL_SERVICE_NAME", oldServiceName)
		}
	}()

	// Ensure no endpoint (to avoid actual OTLP setup)
	oldEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	os.Unsetenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	defer func() {
		if oldEndpoint != "" {
			os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", oldEndpoint)
		}
	}()

	shutdown := Init("my-default-service")
	defer shutdown()

	// This test verifies the function works with default service name
}

func TestInit_EndpointUrlCleaning(t *testing.T) {
	// This test would ideally verify that URLs are cleaned properly,
	// but since the function doesn't expose the cleaned endpoint,
	// we can only test that it doesn't panic with various endpoint formats
	
	testCases := []string{
		"localhost:4317",
		"http://localhost:4317",
		"https://localhost:4317",
		"otel-collector:4317",
	}

	for _, endpoint := range testCases {
		t.Run(endpoint, func(t *testing.T) {
			oldEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
			os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", endpoint)
			defer func() {
				if oldEndpoint != "" {
					os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", oldEndpoint)
				} else {
					os.Unsetenv("OTEL_EXPORTER_OTLP_ENDPOINT")
				}
			}()

			shutdown := Init("test-service")
			defer shutdown()

			// The function should handle various URL formats without panicking
			// Note: This may fail to connect, but shouldn't panic during initialization
		})
	}
}

func TestRegisterFiber_AddsMiddleware(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	
	// Before registering telemetry, count existing handlers
	initialHandlerCount := len(app.GetRoutes())
	
	RegisterFiber(app, "test-service")
	
	// After registering telemetry, there should be at least one new route (/metrics)
	finalHandlerCount := len(app.GetRoutes())
	if finalHandlerCount <= initialHandlerCount {
		t.Error("expected RegisterFiber to add routes (at minimum /metrics endpoint)")
	}
	
	// Check if /metrics route was added
	routes := app.GetRoutes()
	hasMetricsRoute := false
	for _, route := range routes {
		if strings.Contains(route.Path, "metrics") {
			hasMetricsRoute = true
			break
		}
	}
	
	if !hasMetricsRoute {
		t.Error("expected /metrics route to be registered")
	}
}

func TestRegisterFiber_DoesNotPanic(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	
	// This should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("RegisterFiber panicked: %v", r)
		}
	}()
	
	RegisterFiber(app, "test-service")
}

func TestRegisterFiber_WorksWithEmptyServiceName(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	
	// Should handle empty service name gracefully
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("RegisterFiber with empty service name panicked: %v", r)
		}
	}()
	
	RegisterFiber(app, "")
}

func TestRegisterFiber_MultipleCallsDoNotPanic(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	
	// Multiple calls should not cause issues
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Multiple RegisterFiber calls panicked: %v", r)
		}
	}()
	
	RegisterFiber(app, "service1")
	RegisterFiber(app, "service2")
}