package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/apsferreira/auth-service/backend/internal/handler"
	"github.com/apsferreira/auth-service/backend/internal/middleware"
	"github.com/apsferreira/auth-service/backend/internal/pkg/config"
	"github.com/apsferreira/auth-service/backend/internal/pkg/database"
	jwtpkg "github.com/apsferreira/auth-service/backend/internal/pkg/jwt"
	"github.com/apsferreira/auth-service/backend/internal/repository"
	"github.com/apsferreira/auth-service/backend/internal/service"
	"github.com/apsferreira/auth-service/backend/internal/telemetry"
)

func main() {
	// 0. Initialize observability
	shutdownTelemetry := telemetry.Init("auth-service")
	defer shutdownTelemetry()

	// 1. Load config
	cfg := config.Load()

	// 2. Connect to database
	if err := database.Connect(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// 3. Initialize infrastructure services
	jwtService := jwtpkg.NewJWTService(cfg.JWTSecret, cfg.JWTAccessExpiry, cfg.JWTRefreshExpiry)
	emailService := service.NewEmailService(cfg.ResendAPIKey, cfg.ResendFromEmail, cfg.Env)
	telegramNotifier := service.NewTelegramNotifier(cfg.TelegramBotToken, cfg.TelegramChatID)
	whatsappService := service.NewWhatsAppService(cfg.WhatsAppAPIURL, cfg.WhatsAppAPIKey, cfg.WhatsAppInstance, cfg.WhatsAppDefaultPhone)

	// 4. Initialize repositories
	userRepo := repository.NewUserRepository()
	otpRepo := repository.NewOTPRepository()
	tokenRepo := repository.NewTokenRepository()
	oauthRepo := repository.NewOAuthRepository()
	serviceRepo := repository.NewServiceRepository()
	permissionRepo := repository.NewPermissionRepository()
	roleRepo := repository.NewRoleRepository()
	eventRepo := repository.NewEventRepository()

	// 5. Initialize services
	otpService := service.NewOTPService(
		otpRepo,
		cfg.OTPExpiryMinutes,
		cfg.OTPMaxAttempts,
		cfg.OTPRateLimitPerEmail,
		cfg.OTPRateLimitWindowMinutes,
	)
	eventService := service.NewEventService(eventRepo, cfg.Env == "development")
	googleService := service.NewGoogleOAuthService(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURI)
	authService := service.NewAuthService(userRepo, otpService, emailService, telegramNotifier, whatsappService, tokenRepo, jwtService, oauthRepo)
	userService := service.NewUserService(userRepo)
	adminService := service.NewAdminService(serviceRepo, permissionRepo, roleRepo)

	// 6. Initialize handlers
	authHandler := handler.NewAuthHandler(authService, eventService, googleService)
	userHandler := handler.NewUserHandler(userService)
	adminHandler := handler.NewAdminHandler(adminService, eventService)
	healthHandler := handler.NewHealthHandler()

	// 7. Initialize rate limiter (REDUZIDO: 5 requests per hour per IP for OTP endpoints)
	otpRateLimiter := middleware.NewRateLimiter(5, time.Hour)

	// 8. Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	// Observabilidade — OTel tracing middleware + /metrics endpoint
	telemetry.RegisterFiber(app, "auth-service")

	// 9. Global middleware
	app.Use(middleware.Logger())
	app.Use(middleware.CORS(cfg.AllowedOrigins))

	// 10. Health check (public)
	app.Get("/health", healthHandler.Health)

	// 11. API routes
	api := app.Group("/api/v1")

	// Public auth routes
	auth := api.Group("/auth")
	auth.Post("/request-otp", otpRateLimiter.Handler(), authHandler.RequestOTP)
	auth.Post("/verify-otp", authHandler.VerifyOTP)
	auth.Post("/admin-login", authHandler.AdminLogin)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Post("/validate", authHandler.Validate)
	auth.Post("/provision-user", authHandler.ProvisionUser(cfg.ServiceToken))
	// Google OAuth2
	auth.Get("/google", authHandler.GoogleLogin)
	auth.Get("/google/callback", authHandler.GoogleCallback)

	// Protected auth routes
	authProtected := api.Group("/auth", middleware.Auth(jwtService))
	authProtected.Get("/me", authHandler.Me)
	authProtected.Patch("/me", authHandler.UpdateMe)
	authProtected.Post("/logout", authHandler.Logout)

	// Protected user management routes (admin only)
	users := api.Group("/users", middleware.Auth(jwtService), middleware.RequireRole("admin", "super_admin"))
	users.Get("/", userHandler.List)
	users.Post("/", userHandler.Create)
	users.Get("/:id", userHandler.GetByID)
	users.Put("/:id", userHandler.Update)
	users.Delete("/:id", userHandler.Delete)

	// Protected admin routes (admin only) - services, permissions, roles
	admin := api.Group("/admin", middleware.Auth(jwtService), middleware.RequireRole("admin", "super_admin"))

	// Services management
	admin.Get("/services", adminHandler.ListServices)
	admin.Post("/services", adminHandler.CreateService)
	admin.Get("/services/:id", adminHandler.GetService)
	admin.Put("/services/:id", adminHandler.UpdateService)
	admin.Delete("/services/:id", adminHandler.DeleteService)

	// Permissions management (per service)
	admin.Get("/services/:id/permissions", adminHandler.ListServicePermissions)
	admin.Post("/services/:id/permissions", adminHandler.CreateServicePermission)
	admin.Get("/permissions", adminHandler.ListAllPermissions)
	admin.Delete("/permissions/:id", adminHandler.DeletePermission)

	// Roles management
	admin.Get("/roles", adminHandler.ListRoles)
	admin.Post("/roles", adminHandler.CreateRole)
	admin.Put("/roles/:id", adminHandler.UpdateRole)
	admin.Get("/roles/:id/permissions", adminHandler.GetRolePermissions)
	admin.Put("/roles/:id/permissions", adminHandler.SetRolePermissions)

	// Audit events
	admin.Get("/events", adminHandler.ListEvents)

	// 12. Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Auth Service starting on %s (env: %s)", addr, cfg.Env)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
