package service

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/apsferreira/auth-service/backend/internal/domain"
	jwtpkg "github.com/apsferreira/auth-service/backend/internal/pkg/jwt"
	"github.com/apsferreira/auth-service/backend/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo         repository.UserRepositoryInterface
	otpService       *OTPService
	emailService     *EmailService
	telegramNotifier *TelegramNotifier
	whatsappService  *WhatsAppService
	tokenRepo        repository.TokenRepositoryInterface
	jwtService       *jwtpkg.JWTService
	oauthRepo        repository.OAuthRepositoryInterface
}

func NewAuthService(
	userRepo repository.UserRepositoryInterface,
	otpService *OTPService,
	emailService *EmailService,
	telegramNotifier *TelegramNotifier,
	whatsappService *WhatsAppService,
	tokenRepo repository.TokenRepositoryInterface,
	jwtService *jwtpkg.JWTService,
	oauthRepo repository.OAuthRepositoryInterface,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		otpService:       otpService,
		emailService:     emailService,
		telegramNotifier: telegramNotifier,
		whatsappService:  whatsappService,
		tokenRepo:        tokenRepo,
		jwtService:       jwtService,
		oauthRepo:        oauthRepo,
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

	code, err := s.otpService.GenerateAndStore(email)
	if err != nil {
		return nil, "", err
	}

	switch channel {
	case "telegram":
		if s.telegramNotifier == nil || !s.telegramNotifier.IsConfigured() {
			return nil, "", fmt.Errorf("Telegram não está configurado neste serviço")
		}
		if err := s.telegramNotifier.SendOTP(email, code); err != nil {
			return nil, "", fmt.Errorf("falha ao enviar OTP via Telegram: %w", err)
		}
	case "whatsapp":
		if s.whatsappService == nil || !s.whatsappService.IsConfigured() {
			return nil, "", fmt.Errorf("WhatsApp não está configurado neste serviço")
		}
		if err := s.whatsappService.SendOTP(email, code); err != nil {
			return nil, "", fmt.Errorf("falha ao enviar OTP via WhatsApp: %w", err)
		}
	default: // "email"
		if err := s.emailService.SendOTP(email, code); err != nil {
			log.Printf("[WARN] Email OTP delivery failed: %v", err)
			return nil, "", fmt.Errorf("falha ao enviar OTP por email: %w", err)
		}
	}

	channelLabels := map[string]string{
		"email":    "email",
		"telegram": "Telegram",
		"whatsapp": "WhatsApp",
	}
	label := channelLabels[channel]

	return &domain.OTPResponse{
		Message:   "Código enviado via " + label,
		ExpiresIn: s.otpService.GetExpiryMinutes() * 60,
		Channel:   channel,
	}, code, nil
}

// VerifyOTP handles step 2: user submits email + code, we return JWT tokens.
func (s *AuthService) VerifyOTP(email, code string) (*domain.AuthResponse, error) {
	if err := s.otpService.Verify(email, code); err != nil {
		return nil, err
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
func (s *AuthService) ProvisionUser(email, fullName string) (*domain.User, error) {
	existing, err := s.userRepo.FindByEmail(email)
	if err == nil {
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
