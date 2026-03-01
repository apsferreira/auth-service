package handler

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/apsferreira/auth-service/backend/internal/domain"
	"github.com/apsferreira/auth-service/backend/internal/service"
)

type AuthHandler struct {
	authService  *service.AuthService
	eventService *service.EventService
}

func NewAuthHandler(authService *service.AuthService, eventService *service.EventService) *AuthHandler {
	return &AuthHandler{authService: authService, eventService: eventService}
}

// RequestOTP handles POST /api/v1/auth/request-otp
func (h *AuthHandler) RequestOTP(c *fiber.Ctx) error {
	var req domain.OTPRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(domain.ErrorResponse{Error: "invalid request body"})
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" {
		return c.Status(400).JSON(domain.ErrorResponse{Error: "email is required"})
	}

	resp, plainCode, err := h.authService.RequestOTP(req.Email, req.Channel)
	if err != nil {
		if strings.Contains(err.Error(), "rate limit") {
			return c.Status(429).JSON(domain.ErrorResponse{Error: err.Error()})
		}
		if strings.Contains(err.Error(), "não está configurado") {
			return c.Status(400).JSON(domain.ErrorResponse{Error: err.Error()})
		}
		return c.Status(500).JSON(domain.ErrorResponse{Error: err.Error()})
	}

	h.eventService.Log(domain.EventOTPRequested, req.Email, c.IP(), c.Get("User-Agent"), nil, map[string]interface{}{
		"otp_expires_in_minutes": resp.ExpiresIn / 60,
		"otp_code":               plainCode,
		"channel":                resp.Channel,
		"source":                 sourceFromOrigin(c.Get("Origin")),
	})

	return c.Status(200).JSON(resp)
}

// VerifyOTP handles POST /api/v1/auth/verify-otp
func (h *AuthHandler) VerifyOTP(c *fiber.Ctx) error {
	var req domain.OTPVerifyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(domain.ErrorResponse{Error: "invalid request body"})
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Code == "" {
		return c.Status(400).JSON(domain.ErrorResponse{Error: "email and code are required"})
	}
	if len(req.Code) != 6 {
		return c.Status(400).JSON(domain.ErrorResponse{Error: "code must be 6 digits"})
	}

	resp, err := h.authService.VerifyOTP(req.Email, req.Code)
	if err != nil {
		h.eventService.Log(domain.EventLoginFailed, req.Email, c.IP(), c.Get("User-Agent"), nil, map[string]interface{}{
			"method": "otp",
			"error":  err.Error(),
			"source": sourceFromOrigin(c.Get("Origin")),
		})
		return c.Status(401).JSON(domain.ErrorResponse{Error: err.Error()})
	}

	userID := resp.User.ID
	h.eventService.Log(domain.EventLoginSuccess, resp.User.Email, c.IP(), c.Get("User-Agent"), &userID, map[string]interface{}{
		"method":               "otp",
		"access_expires_in":    resp.ExpiresIn,
		"session_active_until": time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second).UTC().Format(time.RFC3339),
		"source":               sourceFromOrigin(c.Get("Origin")),
	})

	return c.Status(200).JSON(resp)
}

// Refresh handles POST /api/v1/auth/refresh
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req domain.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(domain.ErrorResponse{Error: "invalid request body"})
	}

	if req.RefreshToken == "" {
		return c.Status(400).JSON(domain.ErrorResponse{Error: "refresh_token is required"})
	}

	resp, err := h.authService.RefreshTokens(req.RefreshToken)
	if err != nil {
		return c.Status(401).JSON(domain.ErrorResponse{Error: err.Error()})
	}

	userID := resp.User.ID
	h.eventService.Log(domain.EventTokenRefreshed, resp.User.Email, c.IP(), c.Get("User-Agent"), &userID, map[string]interface{}{
		"access_expires_in": resp.ExpiresIn,
		"source":            sourceFromOrigin(c.Get("Origin")),
	})

	return c.Status(200).JSON(resp)
}

// Logout handles POST /api/v1/auth/logout (protected)
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req domain.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(domain.ErrorResponse{Error: "invalid request body"})
	}

	if err := h.authService.Logout(req.RefreshToken); err != nil {
		return c.Status(400).JSON(domain.ErrorResponse{Error: err.Error()})
	}

	email, _ := c.Locals("email").(string)
	userIDStr, _ := c.Locals("userID").(string)
	var userID *uuid.UUID
	if uid, err := uuid.Parse(userIDStr); err == nil {
		userID = &uid
	}
	h.eventService.Log(domain.EventLogout, email, c.IP(), c.Get("User-Agent"), userID, nil)

	return c.Status(200).JSON(fiber.Map{"message": "logged out successfully"})
}

// Me handles GET /api/v1/auth/me (protected)
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(401).JSON(domain.ErrorResponse{Error: "unauthorized"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(400).JSON(domain.ErrorResponse{Error: "invalid user ID"})
	}

	user, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		return c.Status(404).JSON(domain.ErrorResponse{Error: "user not found"})
	}

	return c.Status(200).JSON(user)
}

// AdminLogin handles POST /api/v1/auth/admin-login (username/password for admin panel)
func (h *AuthHandler) AdminLogin(c *fiber.Ctx) error {
	var req domain.AdminLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(domain.ErrorResponse{Error: "invalid request body"})
	}

	req.Identifier = strings.TrimSpace(req.Identifier)
	if req.Identifier == "" || req.Password == "" {
		return c.Status(400).JSON(domain.ErrorResponse{Error: "identifier and password are required"})
	}

	resp, err := h.authService.AdminLogin(req.Identifier, req.Password)
	if err != nil {
		h.eventService.Log(domain.EventLoginFailed, req.Identifier, c.IP(), c.Get("User-Agent"), nil, map[string]interface{}{
			"method": "password",
			"error":  err.Error(),
			"source": sourceFromOrigin(c.Get("Origin")),
		})
		return c.Status(401).JSON(domain.ErrorResponse{Error: err.Error()})
	}

	userID := resp.User.ID
	h.eventService.Log(domain.EventLoginSuccess, resp.User.Email, c.IP(), c.Get("User-Agent"), &userID, map[string]interface{}{
		"method":               "password",
		"access_expires_in":    resp.ExpiresIn,
		"session_active_until": time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second).UTC().Format(time.RFC3339),
		"source":               sourceFromOrigin(c.Get("Origin")),
	})

	return c.Status(200).JSON(resp)
}

// Validate handles POST /api/v1/auth/validate
func (h *AuthHandler) Validate(c *fiber.Ctx) error {
	var req domain.ValidateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(domain.ErrorResponse{Error: "invalid request body"})
	}

	if req.Token == "" {
		return c.Status(400).JSON(domain.ErrorResponse{Error: "token is required"})
	}

	resp, err := h.authService.ValidateToken(req.Token)
	if err != nil {
		return c.Status(200).JSON(domain.ValidateResponse{Valid: false, Message: err.Error()})
	}

	return c.Status(200).JSON(resp)
}

// ProvisionUser handles POST /api/v1/auth/provision-user (service-to-service, protected by X-Service-Token)
// Called by customer-service to ensure a User record exists when a customer registers.
func (h *AuthHandler) ProvisionUser(serviceToken string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("X-Service-Token")
		if token != serviceToken || token == "" {
			return c.Status(401).JSON(domain.ErrorResponse{Error: "invalid service token"})
		}

		var req struct {
			Email    string `json:"email"`
			FullName string `json:"full_name"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(domain.ErrorResponse{Error: "invalid request body"})
		}
		if req.Email == "" {
			return c.Status(400).JSON(domain.ErrorResponse{Error: "email is required"})
		}

		user, err := h.authService.ProvisionUser(req.Email, req.FullName)
		if err != nil {
			return c.Status(500).JSON(domain.ErrorResponse{Error: err.Error()})
		}

		return c.Status(200).JSON(fiber.Map{
			"user_id": user.ID,
			"email":   user.Email,
			"message": "user provisioned",
		})
	}
}

// sourceFromOrigin maps a browser Origin header to a human-readable service name.
func sourceFromOrigin(origin string) string {
	if origin == "" {
		return "direct"
	}
	switch {
	case strings.Contains(origin, ":3003"):
		return "auth-panel"
	case strings.Contains(origin, ":3001"):
		return "my-library"
	case strings.Contains(origin, ":3002"):
		return "focus-hub"
	case strings.Contains(origin, ":3005"):
		return "iit-agents"
	case strings.Contains(origin, ":5173"):
		return "dev-server"
	default:
		return origin
	}
}
