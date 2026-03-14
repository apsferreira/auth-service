package handler

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/apsferreira/auth-service/backend/internal/domain"
	"github.com/apsferreira/auth-service/backend/internal/service"
)

type AuthHandler struct {
	authService   *service.AuthService
	eventService  *service.EventService
	googleService *service.GoogleOAuthService
}

func NewAuthHandler(
	authService *service.AuthService,
	eventService *service.EventService,
	googleService *service.GoogleOAuthService,
) *AuthHandler {
	return &AuthHandler{
		authService:   authService,
		eventService:  eventService,
		googleService: googleService,
	}
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

// UpdateMe handles PATCH /api/v1/auth/me (protected)
func (h *AuthHandler) UpdateMe(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(401).JSON(domain.ErrorResponse{Error: "unauthorized"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(400).JSON(domain.ErrorResponse{Error: "invalid user ID"})
	}

	var req struct {
		FullName string `json:"full_name"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(domain.ErrorResponse{Error: "invalid request body"})
	}
	if req.FullName == "" {
		return c.Status(400).JSON(domain.ErrorResponse{Error: "full_name is required"})
	}

	user, err := h.authService.UpdateCurrentUser(userID, req.FullName)
	if err != nil {
		return c.Status(500).JSON(domain.ErrorResponse{Error: "failed to update profile"})
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

// ProvisionUser handles POST /api/v1/auth/provision-user (service-to-service)
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

// GoogleLogin handles GET /api/v1/auth/google
// Returns the Google OAuth2 authorization URL.
// Query param: redirect_uri — where to send the user after successful auth (frontend URL).
func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	if h.googleService == nil || !h.googleService.IsConfigured() {
		return c.Status(503).JSON(domain.ErrorResponse{Error: "Google OAuth is not configured"})
	}

	redirectURI := c.Query("redirect_uri", "")
	if redirectURI == "" {
		return c.Status(400).JSON(domain.ErrorResponse{Error: "redirect_uri query param is required"})
	}

	// Encode the frontend redirect URI in state so we can use it in the callback
	state := base64.URLEncoding.EncodeToString([]byte(redirectURI))
	authURL := h.googleService.GetAuthURL(state)

	return c.Status(200).JSON(domain.GoogleLoginURLResponse{URL: authURL})
}

// GoogleCallback handles GET /api/v1/auth/google/callback
// Google redirects here after the user authorizes the app.
// Exchanges the code, finds/creates the user, and redirects to the frontend with tokens.
func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")
	errParam := c.Query("error")

	if errParam != "" {
		return c.Status(400).JSON(domain.ErrorResponse{Error: "Google OAuth error: " + errParam})
	}
	if code == "" {
		return c.Status(400).JSON(domain.ErrorResponse{Error: "missing code parameter"})
	}

	// Decode the frontend redirect URI from state
	redirectURIBytes, err := base64.URLEncoding.DecodeString(state)
	if err != nil || len(redirectURIBytes) == 0 {
		return c.Status(400).JSON(domain.ErrorResponse{Error: "invalid state parameter"})
	}
	frontendRedirectURI := string(redirectURIBytes)

	// Exchange code for Google user info
	googleUser, err := h.googleService.ExchangeCode(code)
	if err != nil {
		return c.Status(401).JSON(domain.ErrorResponse{Error: "failed to authenticate with Google: " + err.Error()})
	}

	// Find or create user, generate JWT
	authResp, err := h.authService.LoginWithGoogle(googleUser)
	if err != nil {
		return c.Status(500).JSON(domain.ErrorResponse{Error: "login failed: " + err.Error()})
	}

	h.eventService.Log(domain.EventLoginSuccess, authResp.User.Email, c.IP(), c.Get("User-Agent"), &authResp.User.ID, map[string]interface{}{
		"method":            "google_oauth",
		"access_expires_in": authResp.ExpiresIn,
	})

	// Redirect to frontend with tokens in query string
	redirectURL := fmt.Sprintf("%s?access_token=%s&refresh_token=%s&expires_in=%d",
		frontendRedirectURI,
		authResp.AccessToken,
		authResp.RefreshToken,
		authResp.ExpiresIn,
	)
	return c.Redirect(redirectURL, fiber.StatusTemporaryRedirect)
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
