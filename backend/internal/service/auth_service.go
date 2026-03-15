package service

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/apsferreira/auth-service/backend/internal/client"
	"github.com/apsferreira/auth-service/backend/internal/domain"
	jwtpkg "github.com/apsferreira/auth-service/backend/internal/pkg/jwt"
	"github.com/apsferreira/auth-service/backend/internal/pkg/messaging"
	"github.com/apsferreira/auth-service/backend/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo           repository.UserRepositoryInterface
	notificationClient *client.NotificationClient
	whatsappService    *WhatsAppService
	tokenRepo          repository.TokenRepositoryInterface
	jwtService         *jwtpkg.JWTService
	oauthRepo          repository.OAuthRepositoryInterface
	publisher          *messaging.Publisher // nil = fallback to HTTP
	serviceName        string               // human-readable name used in OTP event payloads
}

func NewAuthService(
	userRepo repository.UserRepositoryInterface,
	notificationClient *client.NotificationClient,
	whatsappService *WhatsAppService,
	tokenRepo repository.TokenRepositoryInterface,
	jwtService *jwtpkg.JWTService,
	oauthRepo repository.OAuthRepositoryInterface,
	publisher *messaging.Publisher,
	serviceName string,
) *AuthService {
	if serviceName == "" {
		serviceName = "Instituto Itinerante"
	}
	return &AuthService{
		userRepo:           userRepo,
		notificationClient: notificationClient,
		whatsappService:    whatsappService,
		tokenRepo:          tokenRepo,
		jwtService:         jwtService,
		oauthRepo:          oauthRepo,
		publisher:          publisher,
		serviceName:        serviceName,
	}
}

// RequestOTP handles step 1: user submits email + preferred channel, we send an OTP.
// channel: "email" | "telegram" | "whatsapp" — defaults to "email".
// Returns (response, plainCode, error) — plainCode is for audit logging only, never sent to client.
func (s *AuthService) RequestOTP(email, channel string) (*domain.OTPResponse, string, error) {
	if channel == "" {
		channel = "email"
	}

	// Find or create user (auto-registration)
	_, err := s.userRepo.FindByEmail(email)
	if err == sql.ErrNoRows {
		if err := s.createDefaultUser(email); err != nil {
			return nil, "", fmt.Errorf("failed to create user: %w", err)
		}
	} else if err != nil {
		return nil, "", fmt.Errorf("failed to find user: %w", err)
	}

	switch channel {
	case "whatsapp":
		if s.whatsappService == nil || !s.whatsappService.IsConfigured() {
			return nil, "", fmt.Errorf("WhatsApp não está configurado neste serviço")
		}
		// For now, generate and send via WhatsApp using local service
		// TODO: Move WhatsApp to notification service later
		code, err := generateOTPCode()
		if err != nil {
			return nil, "", err
		}
		if err := s.whatsappService.SendOTP(email, code); err != nil {
			return nil, "", fmt.Errorf("falha ao enviar OTP via WhatsApp: %w", err)
		}
		return &domain.OTPResponse{
			Message:   "Código enviado via WhatsApp",
			ExpiresIn: 600, // 10 minutes
			Channel:   channel,
		}, code, nil
	default: // "email" or "telegram"
		channelLabel := "email"
		if channel == "telegram" {
			channelLabel = "Telegram"
		}

		if s.publisher != nil {
			// Async path: publish to RabbitMQ — notification-service consumes and sends OTP
			if err := s.publisher.PublishOTPRequested(email, s.serviceName, channel); err != nil {
				log.Printf("[auth] RabbitMQ publish failed, falling back to HTTP: %v", err)
				// Fallthrough to HTTP fallback below
			} else {
				return &domain.OTPResponse{
					Message:   "Código enviado via " + channelLabel,
					ExpiresIn: 600, // 10 minutes (notification-service default)
					Channel:   channel,
				}, "", nil
			}
		}

		// HTTP fallback (used when RABBITMQ_URL is not set or publish failed)
		otpResp, err := s.notificationClient.SendOTP(email)
		if err != nil {
			return nil, "", fmt.Errorf("failed to send OTP: %w", err)
		}

		expiresIn := int(otpResp.ExpiresAt.Sub(time.Now()).Seconds())
		if expiresIn <= 0 {
			expiresIn = 600
		}

		return &domain.OTPResponse{
			Message:   "Código enviado via " + channelLabel,
			ExpiresIn: expiresIn,
			Channel:   channel,
		}, "", nil
	}
}

// VerifyOTP handles step 2: user submits email + code, we return JWT tokens.
func (s *AuthService) VerifyOTP(email, code string) (*domain.AuthResponse, error) {
	// Use notification service to verify OTP
	verifyResp, err := s.notificationClient.VerifyOTP(email, code)
	if err != nil {
		return nil, fmt.Errorf("failed to verify OTP: %w", err)
	}
	
	if !verifyResp.Valid {
		return nil, fmt.Errorf(verifyResp.Message)
	}

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	_ = s.userRepo.UpdateLastLogin(user.ID)

	roles, permissions, err := s.userRepo.GetUserRolesAndPermissions(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.TenantID, user.Email, roles, permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshTokenValue, refreshTokenHash, err := s.jwtService.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	rt := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: refreshTokenHash,
		ExpiresAt: time.Now().Add(s.jwtService.GetRefreshExpiry()),
		CreatedAt: time.Now(),
	}
	if err := s.tokenRepo.Create(rt); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenValue,
		ExpiresIn:    int(s.jwtService.GetAccessExpiry().Seconds()),
		User:         user,
		Roles:        roles,
		Permissions:  permissions,
	}, nil
}

// RefreshTokens handles step 3: exchange refresh token for new pair.
func (s *AuthService) RefreshTokens(refreshTokenValue string) (*domain.AuthResponse, error) {
	hash := jwtpkg.HashToken(refreshTokenValue)

	storedToken, err := s.tokenRepo.FindByHash(hash)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired refresh token")
	}

	// Revoke old token (rotation)
	if err := s.tokenRepo.Revoke(storedToken.ID); err != nil {
		return nil, fmt.Errorf("failed to revoke old token: %w", err)
	}

	user, err := s.userRepo.FindByID(storedToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	roles, permissions, err := s.userRepo.GetUserRolesAndPermissions(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.TenantID, user.Email, roles, permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshValue, newRefreshHash, err := s.jwtService.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	rt := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: newRefreshHash,
		ExpiresAt: time.Now().Add(s.jwtService.GetRefreshExpiry()),
		CreatedAt: time.Now(),
	}
	if err := s.tokenRepo.Create(rt); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshValue,
		ExpiresIn:    int(s.jwtService.GetAccessExpiry().Seconds()),
		User:         user,
		Roles:        roles,
		Permissions:  permissions,
	}, nil
}

// ValidateToken validates an access token and returns claims.
func (s *AuthService) ValidateToken(tokenString string) (*domain.ValidateResponse, error) {
	claims, err := s.jwtService.ValidateAccessToken(tokenString)
	if err != nil {
		return &domain.ValidateResponse{Valid: false, Message: err.Error()}, nil
	}
	return &domain.ValidateResponse{
		Valid:       true,
		UserID:      &claims.Subject,
		TenantID:    &claims.TenantID,
		Email:       &claims.Email,
		Roles:       claims.Roles,
		Permissions: claims.Permissions,
	}, nil
}

// Logout revokes a specific refresh token.
func (s *AuthService) Logout(refreshTokenValue string) error {
	hash := jwtpkg.HashToken(refreshTokenValue)
	token, err := s.tokenRepo.FindByHash(hash)
	if err != nil {
		return fmt.Errorf("invalid refresh token")
	}
	return s.tokenRepo.Revoke(token.ID)
}

// GetCurrentUser returns the user profile for the authenticated user.
func (s *AuthService) GetCurrentUser(userID uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	roles, _, err := s.userRepo.GetUserRolesAndPermissions(userID)
	if err == nil {
		for _, roleName := range roles {
			user.Roles = append(user.Roles, &domain.Role{Name: roleName})
		}
	}

	return user, nil
}

// UpdateCurrentUser updates the authenticated user's profile.
func (s *AuthService) UpdateCurrentUser(userID uuid.UUID, fullName string) (*domain.User, error) {
	return s.userRepo.UpdateProfile(userID, fullName)
}

// AdminLogin authenticates an admin user with username/password (admin panel only).
func (s *AuthService) AdminLogin(identifier, password string) (*domain.AuthResponse, error) {
	user, err := s.userRepo.FindByUsernameOrEmail(identifier)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if user.PasswordHash == nil {
		return nil, fmt.Errorf("password authentication not configured for this account")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	roles, permissions, err := s.userRepo.GetUserRolesAndPermissions(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	// Admin panel login requires admin or super_admin role
	isAdmin := false
	for _, r := range roles {
		if r == "admin" || r == "super_admin" {
			isAdmin = true
			break
		}
	}
	if !isAdmin {
		return nil, fmt.Errorf("insufficient privileges: admin or super_admin role required")
	}

	_ = s.userRepo.UpdateLastLogin(user.ID)

	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.TenantID, user.Email, roles, permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshTokenValue, refreshTokenHash, err := s.jwtService.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	rt := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: refreshTokenHash,
		ExpiresAt: time.Now().Add(s.jwtService.GetRefreshExpiry()),
		CreatedAt: time.Now(),
	}
	if err := s.tokenRepo.Create(rt); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenValue,
		ExpiresIn:    int(s.jwtService.GetAccessExpiry().Seconds()),
		User:         user,
		Roles:        roles,
		Permissions:  permissions,
	}, nil
}


// ProvisionUser ensures a user account exists for the given email (called by customer-service).
// If the user already exists and has an empty full_name, it is updated with the provided value.
func (s *AuthService) ProvisionUser(email, fullName string) (*domain.User, error) {
	existing, err := s.userRepo.FindByEmail(email)
	if err == nil {
		// Update full_name if auth-service record is out of sync (e.g. first OTP created empty name)
		if existing.FullName == "" && fullName != "" {
			updated, updateErr := s.userRepo.UpdateProfile(existing.ID, fullName)
			if updateErr == nil {
				return updated, nil
			}
			// Non-fatal: log and return existing if update fails
			log.Printf("ProvisionUser: failed to update full_name for %s: %v", email, updateErr)
		}
		return existing, nil
	}

	tenantID, err := s.userRepo.GetDefaultTenantID()
	if err != nil {
		return nil, fmt.Errorf("no default tenant found: %w", err)
	}
	roleID, err := s.userRepo.GetDefaultRoleID(tenantID)
	if err != nil {
		return nil, fmt.Errorf("no default role found: %w", err)
	}

	now := time.Now()
	user := &domain.User{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Email:     email,
		FullName:  fullName,
		IsActive:  true,
		RoleID:    &roleID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	if err := s.userRepo.AddUserRole(user.ID, roleID); err != nil {
		return nil, fmt.Errorf("failed to assign role: %w", err)
	}
	return user, nil
}

// LoginWithGoogle handles the complete Google OAuth login flow:
// finds or creates the user, links the OAuth identity, and returns JWT tokens.
func (s *AuthService) LoginWithGoogle(googleUser *GoogleUserInfo) (*domain.AuthResponse, error) {
	var user *domain.User

	// 1. Check if we already have an OAuth identity for this Google account
	identity, err := s.oauthRepo.FindByProvider("google", googleUser.ID)
	if err == nil {
		// Identity exists — load the linked user
		user, err = s.userRepo.FindByID(identity.UserID)
		if err != nil {
			return nil, fmt.Errorf("linked user not found: %w", err)
		}
	} else {
		// 2. No identity — find user by email (may have registered via OTP before)
		user, err = s.userRepo.FindByEmail(googleUser.Email)
		if err == sql.ErrNoRows {
			// 3. Brand new user — create account
			if err := s.createGoogleUser(googleUser); err != nil {
				return nil, fmt.Errorf("failed to create user: %w", err)
			}
			user, err = s.userRepo.FindByEmail(googleUser.Email)
			if err != nil {
				return nil, fmt.Errorf("failed to load new user: %w", err)
			}
		} else if err != nil {
			return nil, fmt.Errorf("failed to find user: %w", err)
		}

		// Link the Google identity to this user (new or existing)
		avatarURL := googleUser.Picture
		newIdentity := &domain.OAuthIdentity{
			UserID:     user.ID,
			Provider:   "google",
			ProviderID: googleUser.ID,
			Email:      googleUser.Email,
			AvatarURL:  &avatarURL,
		}
		if err := s.oauthRepo.Upsert(newIdentity); err != nil {
			log.Printf("[WARN] Failed to upsert OAuth identity: %v", err)
		}
	}

	// 4. Update avatar if user doesn't have one yet and Google provided it
	if user.AvatarURL == nil && googleUser.Picture != "" {
		avatarURL := googleUser.Picture
		req := domain.UserUpdateRequest{AvatarURL: &avatarURL}
		updated, err := s.userRepo.Update(user.ID, req)
		if err == nil {
			user = updated
		}
	}

	// 5. Update full name if empty
	if user.FullName == "" && googleUser.Name != "" {
		req := domain.UserUpdateRequest{FullName: &googleUser.Name}
		updated, err := s.userRepo.Update(user.ID, req)
		if err == nil {
			user = updated
		}
	}

	_ = s.userRepo.UpdateLastLogin(user.ID)

	roles, permissions, err := s.userRepo.GetUserRolesAndPermissions(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.TenantID, user.Email, roles, permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshTokenValue, refreshTokenHash, err := s.jwtService.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	rt := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: refreshTokenHash,
		ExpiresAt: time.Now().Add(s.jwtService.GetRefreshExpiry()),
		CreatedAt: time.Now(),
	}
	if err := s.tokenRepo.Create(rt); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenValue,
		ExpiresIn:    int(s.jwtService.GetAccessExpiry().Seconds()),
		User:         user,
		Roles:        roles,
		Permissions:  permissions,
	}, nil
}

func (s *AuthService) createGoogleUser(googleUser *GoogleUserInfo) error {
	tenantID, err := s.userRepo.GetDefaultTenantID()
	if err != nil {
		return fmt.Errorf("no default tenant: %w", err)
	}
	roleID, err := s.userRepo.GetDefaultRoleID(tenantID)
	if err != nil {
		return fmt.Errorf("no default role: %w", err)
	}

	now := time.Now()
	avatarURL := googleUser.Picture
	user := &domain.User{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Email:     googleUser.Email,
		FullName:  googleUser.Name,
		AvatarURL: &avatarURL,
		IsActive:  true,
		RoleID:    &roleID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.userRepo.Create(user); err != nil {
		return err
	}
	return s.userRepo.AddUserRole(user.ID, roleID)
}

func (s *AuthService) createDefaultUser(email string) error {
	tenantID, err := s.userRepo.GetDefaultTenantID()
	if err != nil {
		return fmt.Errorf("no default tenant found: %w", err)
	}

	roleID, err := s.userRepo.GetDefaultRoleID(tenantID)
	if err != nil {
		return fmt.Errorf("no default role found: %w", err)
	}

	now := time.Now()
	user := &domain.User{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Email:     email,
		FullName:  "",
		IsActive:  true,
		RoleID:    &roleID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.userRepo.Create(user); err != nil {
		return err
	}

	return s.userRepo.AddUserRole(user.ID, roleID)
}

// generateOTPCode generates a 6-digit OTP code (helper for WhatsApp until it's migrated)
func generateOTPCode() (string, error) {
	// For simplicity, generate a random 6-digit code
	// This is a temporary solution until WhatsApp is also migrated to notification service
	now := time.Now()
	code := fmt.Sprintf("%06d", (now.Unix()%1000000))
	return code, nil
}
