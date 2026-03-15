package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env                       string
	Port                      string
	DatabaseURL               string
	JWTSecret                 string
	JWTAccessExpiry           time.Duration
	JWTRefreshExpiry          time.Duration
	ResendAPIKey              string
	ResendFromEmail           string
	AllowedOrigins            string
	OTPExpiryMinutes          int
	OTPMaxAttempts            int
	OTPRateLimitPerEmail      int
	OTPRateLimitWindowMinutes int
	ServiceToken              string // shared secret for service-to-service calls
	ServiceName               string // human-readable name used in notification payloads
	NotificationServiceURL    string // URL for notification service API
	RabbitMQURL               string // amqp://user:pass@host:5672/
	// Google OAuth2
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURI  string
	TelegramBotToken          string
	TelegramChatID            string
	WhatsAppAPIURL            string
	WhatsAppAPIKey            string
	WhatsAppInstance          string
	WhatsAppDefaultPhone      string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Env:                       getEnv("ENV", "development"),
		Port:                      getEnv("PORT", "3002"),
		DatabaseURL:               getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5433/auth_db?sslmode=disable"),
		JWTSecret:                 getEnv("JWT_SECRET", "change-me-in-production"),
		JWTAccessExpiry:           getDurationEnv("JWT_ACCESS_EXPIRY", 15*time.Minute),
		JWTRefreshExpiry:          getDurationEnv("JWT_REFRESH_EXPIRY", 168*time.Hour),
		ResendAPIKey:              getEnv("RESEND_API_KEY", ""),
		ResendFromEmail:           getEnv("RESEND_FROM_EMAIL", "auth@yourdomain.com"),
		AllowedOrigins:            getEnv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173"),
		OTPExpiryMinutes:          getIntEnv("OTP_EXPIRY_MINUTES", 10),
		OTPMaxAttempts:            getIntEnv("OTP_MAX_ATTEMPTS", 3),
		OTPRateLimitPerEmail:      getIntEnv("OTP_RATE_LIMIT_PER_EMAIL", 3),
		OTPRateLimitWindowMinutes: getIntEnv("OTP_RATE_LIMIT_WINDOW_MINUTES", 15),
		ServiceToken:              getEnv("SERVICE_TOKEN", "service-secret-change-in-production"),
		ServiceName:               getEnv("SERVICE_NAME", "Instituto Itinerante"),
		NotificationServiceURL:    getEnv("NOTIFICATION_SERVICE_URL", "http://localhost:3030"),
		RabbitMQURL:               getEnv("RABBITMQ_URL", ""),
		GoogleClientID:            getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret:        getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURI:         getEnv("GOOGLE_REDIRECT_URI", ""),
		TelegramBotToken:          getEnv("TELEGRAM_BOT_TOKEN", ""),
		TelegramChatID:            getEnv("TELEGRAM_CHAT_ID", ""),
		WhatsAppAPIURL:            getEnv("WHATSAPP_API_URL", "http://shared-evolution-api:8080"),
		WhatsAppAPIKey:            getEnv("WHATSAPP_API_KEY", ""),
		WhatsAppInstance:          getEnv("WHATSAPP_INSTANCE", "auth-service"),
		WhatsAppDefaultPhone:      getEnv("WHATSAPP_DEFAULT_PHONE", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if v, err := strconv.Atoi(value); err == nil {
			return v
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}
