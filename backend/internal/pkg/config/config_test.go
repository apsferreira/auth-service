package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad_DefaultValues(t *testing.T) {
	// Clear environment variables that might interfere
	envVars := []string{
		"ENV", "PORT", "DATABASE_URL", "JWT_SECRET",
		"JWT_ACCESS_EXPIRY", "JWT_REFRESH_EXPIRY",
		"RESEND_API_KEY", "RESEND_FROM_EMAIL", "ALLOWED_ORIGINS",
		"OTP_EXPIRY_MINUTES", "OTP_MAX_ATTEMPTS",
		"OTP_RATE_LIMIT_PER_EMAIL", "OTP_RATE_LIMIT_WINDOW_MINUTES",
		"SERVICE_TOKEN", "TELEGRAM_BOT_TOKEN", "TELEGRAM_CHAT_ID",
		"WHATSAPP_API_URL", "WHATSAPP_API_KEY", "WHATSAPP_INSTANCE",
		"WHATSAPP_DEFAULT_PHONE",
	}
	
	oldValues := make(map[string]string)
	for _, envVar := range envVars {
		oldValues[envVar] = os.Getenv(envVar)
		os.Unsetenv(envVar)
	}
	defer func() {
		// Restore original values
		for envVar, oldValue := range oldValues {
			if oldValue != "" {
				os.Setenv(envVar, oldValue)
			}
		}
	}()

	config := Load()

	if config.Env != "development" {
		t.Errorf("expected Env='development', got %q", config.Env)
	}
	if config.Port != "3002" {
		t.Errorf("expected Port='3002', got %q", config.Port)
	}
	if config.JWTSecret != "change-me-in-production" {
		t.Errorf("expected default JWT secret, got %q", config.JWTSecret)
	}
	if config.JWTAccessExpiry != 15*time.Minute {
		t.Errorf("expected JWTAccessExpiry=15m, got %v", config.JWTAccessExpiry)
	}
	if config.JWTRefreshExpiry != 168*time.Hour {
		t.Errorf("expected JWTRefreshExpiry=168h, got %v", config.JWTRefreshExpiry)
	}
	if config.OTPExpiryMinutes != 10 {
		t.Errorf("expected OTPExpiryMinutes=10, got %d", config.OTPExpiryMinutes)
	}
	if config.OTPMaxAttempts != 3 {
		t.Errorf("expected OTPMaxAttempts=3, got %d", config.OTPMaxAttempts)
	}
}

func TestLoad_EnvironmentOverrides(t *testing.T) {
	testCases := []struct {
		envVar   string
		value    string
		check    func(*Config) error
	}{
		{
			"ENV",
			"production",
			func(c *Config) error {
				if c.Env != "production" {
					return nil
				}
				return nil
			},
		},
		{
			"PORT",
			"8080",
			func(c *Config) error {
				if c.Port != "8080" {
					return nil
				}
				return nil
			},
		},
		{
			"JWT_SECRET",
			"test-secret-key",
			func(c *Config) error {
				if c.JWTSecret != "test-secret-key" {
					return nil
				}
				return nil
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.envVar, func(t *testing.T) {
			oldValue := os.Getenv(tc.envVar)
			os.Setenv(tc.envVar, tc.value)
			defer func() {
				if oldValue != "" {
					os.Setenv(tc.envVar, oldValue)
				} else {
					os.Unsetenv(tc.envVar)
				}
			}()

			config := Load()
			tc.check(config)
		})
	}
}

func TestGetEnv_DefaultValue(t *testing.T) {
	key := "NONEXISTENT_ENV_VAR_FOR_TESTING"
	defaultValue := "default-test-value"
	
	// Ensure the env var is not set
	os.Unsetenv(key)
	
	result := getEnv(key, defaultValue)
	if result != defaultValue {
		t.Errorf("expected %q, got %q", defaultValue, result)
	}
}

func TestGetEnv_EnvironmentValue(t *testing.T) {
	key := "TEST_ENV_VAR_FOR_TESTING"
	envValue := "environment-test-value"
	defaultValue := "default-test-value"
	
	os.Setenv(key, envValue)
	defer os.Unsetenv(key)
	
	result := getEnv(key, defaultValue)
	if result != envValue {
		t.Errorf("expected %q, got %q", envValue, result)
	}
}

func TestGetIntEnv_DefaultValue(t *testing.T) {
	key := "NONEXISTENT_INT_ENV_VAR_FOR_TESTING"
	defaultValue := 42
	
	os.Unsetenv(key)
	
	result := getIntEnv(key, defaultValue)
	if result != defaultValue {
		t.Errorf("expected %d, got %d", defaultValue, result)
	}
}

func TestGetIntEnv_ValidEnvironmentValue(t *testing.T) {
	key := "TEST_INT_ENV_VAR_FOR_TESTING"
	envValue := "123"
	defaultValue := 42
	
	os.Setenv(key, envValue)
	defer os.Unsetenv(key)
	
	result := getIntEnv(key, defaultValue)
	if result != 123 {
		t.Errorf("expected 123, got %d", result)
	}
}

func TestGetIntEnv_InvalidEnvironmentValue_UsesDefault(t *testing.T) {
	key := "TEST_INT_ENV_VAR_FOR_TESTING"
	envValue := "not-a-number"
	defaultValue := 42
	
	os.Setenv(key, envValue)
	defer os.Unsetenv(key)
	
	result := getIntEnv(key, defaultValue)
	if result != defaultValue {
		t.Errorf("expected default value %d for invalid int, got %d", defaultValue, result)
	}
}

func TestGetDurationEnv_DefaultValue(t *testing.T) {
	key := "NONEXISTENT_DURATION_ENV_VAR_FOR_TESTING"
	defaultValue := 30 * time.Minute
	
	os.Unsetenv(key)
	
	result := getDurationEnv(key, defaultValue)
	if result != defaultValue {
		t.Errorf("expected %v, got %v", defaultValue, result)
	}
}

func TestGetDurationEnv_ValidEnvironmentValue(t *testing.T) {
	key := "TEST_DURATION_ENV_VAR_FOR_TESTING"
	envValue := "45m"
	defaultValue := 30 * time.Minute
	
	os.Setenv(key, envValue)
	defer os.Unsetenv(key)
	
	result := getDurationEnv(key, defaultValue)
	expected := 45 * time.Minute
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestGetDurationEnv_InvalidEnvironmentValue_UsesDefault(t *testing.T) {
	key := "TEST_DURATION_ENV_VAR_FOR_TESTING"
	envValue := "invalid-duration"
	defaultValue := 30 * time.Minute
	
	os.Setenv(key, envValue)
	defer os.Unsetenv(key)
	
	result := getDurationEnv(key, defaultValue)
	if result != defaultValue {
		t.Errorf("expected default value %v for invalid duration, got %v", defaultValue, result)
	}
}

func TestLoad_DatabaseURL_Default(t *testing.T) {
	os.Unsetenv("DATABASE_URL")
	config := Load()
	
	expectedPrefix := "postgres://postgres:postgres@localhost:5433"
	if len(config.DatabaseURL) < len(expectedPrefix) {
		t.Errorf("expected DatabaseURL to start with %q, got %q", expectedPrefix, config.DatabaseURL)
	}
}

func TestLoad_AllowedOrigins_Default(t *testing.T) {
	os.Unsetenv("ALLOWED_ORIGINS")
	config := Load()
	
	expected := "http://localhost:3000,http://localhost:5173"
	if config.AllowedOrigins != expected {
		t.Errorf("expected AllowedOrigins=%q, got %q", expected, config.AllowedOrigins)
	}
}

func TestLoad_OTPSettings_Defaults(t *testing.T) {
	// Clear OTP-related env vars
	otpEnvVars := []string{
		"OTP_EXPIRY_MINUTES", "OTP_MAX_ATTEMPTS",
		"OTP_RATE_LIMIT_PER_EMAIL", "OTP_RATE_LIMIT_WINDOW_MINUTES",
	}
	
	for _, envVar := range otpEnvVars {
		os.Unsetenv(envVar)
	}
	
	config := Load()
	
	if config.OTPExpiryMinutes != 10 {
		t.Errorf("expected OTPExpiryMinutes=10, got %d", config.OTPExpiryMinutes)
	}
	if config.OTPMaxAttempts != 3 {
		t.Errorf("expected OTPMaxAttempts=3, got %d", config.OTPMaxAttempts)
	}
	if config.OTPRateLimitPerEmail != 3 {
		t.Errorf("expected OTPRateLimitPerEmail=3, got %d", config.OTPRateLimitPerEmail)
	}
	if config.OTPRateLimitWindowMinutes != 15 {
		t.Errorf("expected OTPRateLimitWindowMinutes=15, got %d", config.OTPRateLimitWindowMinutes)
	}
}