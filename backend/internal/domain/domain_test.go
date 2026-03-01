package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

// --- Event type constant tests ---

func TestEventConstants_HaveExpectedValues(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"OTPRequested", EventOTPRequested, "otp_requested"},
		{"LoginSuccess", EventLoginSuccess, "login_success"},
		{"LoginFailed", EventLoginFailed, "login_failed"},
		{"Logout", EventLogout, "logout"},
		{"TokenRefreshed", EventTokenRefreshed, "token_refreshed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("expected %s = %q, got %q", tt.name, tt.expected, tt.constant)
			}
		})
	}
}

// --- OTPCode tests ---

func TestOTPCode_Fields(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	otp := &OTPCode{
		ID:        id,
		Email:     "user@example.com",
		CodeHash:  "$2b$10$hash",
		Channel:   "email",
		Attempts:  0,
		ExpiresAt: now.Add(10 * time.Minute),
		CreatedAt: now,
	}

	if otp.ID != id {
		t.Errorf("expected ID %s, got %s", id, otp.ID)
	}
	if otp.Email != "user@example.com" {
		t.Errorf("expected email 'user@example.com', got %s", otp.Email)
	}
	if otp.Channel != "email" {
		t.Errorf("expected channel 'email', got %s", otp.Channel)
	}
	if otp.Attempts != 0 {
		t.Errorf("expected 0 attempts, got %d", otp.Attempts)
	}
	if otp.ExpiresAt.Before(now) {
		t.Error("expected ExpiresAt to be in the future")
	}
}

func TestOTPCode_IsExpired(t *testing.T) {
	past := time.Now().Add(-1 * time.Minute)
	otp := &OTPCode{
		ExpiresAt: past,
	}

	if !time.Now().After(otp.ExpiresAt) {
		t.Error("expected OTP to be expired")
	}
}

func TestOTPCode_IsNotExpired(t *testing.T) {
	future := time.Now().Add(10 * time.Minute)
	otp := &OTPCode{
		ExpiresAt: future,
	}

	if time.Now().After(otp.ExpiresAt) {
		t.Error("expected OTP to not be expired yet")
	}
}

// --- RefreshToken tests ---

func TestRefreshToken_Fields(t *testing.T) {
	id := uuid.New()
	userID := uuid.New()
	now := time.Now()

	rt := &RefreshToken{
		ID:        id,
		UserID:    userID,
		TokenHash: "sha256hashhere",
		ExpiresAt: now.Add(7 * 24 * time.Hour),
		CreatedAt: now,
	}

	if rt.ID != id {
		t.Errorf("expected ID %s, got %s", id, rt.ID)
	}
	if rt.UserID != userID {
		t.Errorf("expected UserID %s, got %s", userID, rt.UserID)
	}
	if rt.RevokedAt != nil {
		t.Error("expected RevokedAt to be nil for new refresh token")
	}
	if rt.ExpiresAt.Before(now) {
		t.Error("expected ExpiresAt to be in the future")
	}
}

func TestRefreshToken_IsRevoked(t *testing.T) {
	now := time.Now()
	rt := &RefreshToken{
		RevokedAt: &now,
	}

	if rt.RevokedAt == nil {
		t.Error("expected RevokedAt to be set")
	}
}

// --- User struct tests ---

func TestUser_Fields(t *testing.T) {
	id := uuid.New()
	tenantID := uuid.New()
	now := time.Now()

	user := &User{
		ID:       id,
		TenantID: tenantID,
		Email:    "admin@example.com",
		FullName: "Admin User",
		IsActive: true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if user.ID != id {
		t.Errorf("expected ID %s, got %s", id, user.ID)
	}
	if user.TenantID != tenantID {
		t.Errorf("expected TenantID %s, got %s", tenantID, user.TenantID)
	}
	if user.Email != "admin@example.com" {
		t.Errorf("expected email 'admin@example.com', got %s", user.Email)
	}
	if !user.IsActive {
		t.Error("expected IsActive to be true")
	}
	if user.PasswordHash != nil {
		t.Error("PasswordHash should be nil when not set")
	}
	if user.AvatarURL != nil {
		t.Error("AvatarURL should be nil when not set")
	}
	if user.LastLoginAt != nil {
		t.Error("LastLoginAt should be nil for new user")
	}
}

func TestUser_InactiveUser(t *testing.T) {
	user := &User{
		ID:       uuid.New(),
		Email:    "inactive@example.com",
		IsActive: false,
	}

	if user.IsActive {
		t.Error("expected IsActive to be false for inactive user")
	}
}

// --- OTPRequest tests ---

func TestOTPRequest_DefaultChannel(t *testing.T) {
	req := &OTPRequest{
		Email:   "user@example.com",
		Channel: "",
	}

	// The channel defaults to "email" — enforced in service layer
	// Here we just verify the struct can be constructed with empty channel
	if req.Email != "user@example.com" {
		t.Errorf("expected email 'user@example.com', got %s", req.Email)
	}
}

func TestOTPRequest_TelegramChannel(t *testing.T) {
	req := &OTPRequest{
		Email:   "user@example.com",
		Channel: "telegram",
	}

	if req.Channel != "telegram" {
		t.Errorf("expected channel 'telegram', got %s", req.Channel)
	}
}

func TestOTPRequest_WhatsAppChannel(t *testing.T) {
	req := &OTPRequest{
		Email:   "user@example.com",
		Channel: "whatsapp",
	}

	if req.Channel != "whatsapp" {
		t.Errorf("expected channel 'whatsapp', got %s", req.Channel)
	}
}

// --- OTPVerifyRequest tests ---

func TestOTPVerifyRequest_Fields(t *testing.T) {
	req := &OTPVerifyRequest{
		Email: "user@example.com",
		Code:  "123456",
	}

	if req.Email != "user@example.com" {
		t.Errorf("expected email 'user@example.com', got %s", req.Email)
	}
	if req.Code != "123456" {
		t.Errorf("expected code '123456', got %s", req.Code)
	}
	if len(req.Code) != 6 {
		t.Errorf("expected 6-digit code, got %d chars", len(req.Code))
	}
}

// --- AuthResponse tests ---

func TestAuthResponse_Fields(t *testing.T) {
	user := &User{ID: uuid.New(), Email: "user@example.com"}
	resp := &AuthResponse{
		AccessToken:  "access.jwt.token",
		RefreshToken: "refresh-token-hex",
		ExpiresIn:    900,
		User:         user,
		Roles:        []string{"admin"},
		Permissions:  map[string][]string{"resource": {"read"}},
	}

	if resp.AccessToken == "" {
		t.Error("expected non-empty AccessToken")
	}
	if resp.RefreshToken == "" {
		t.Error("expected non-empty RefreshToken")
	}
	if resp.ExpiresIn != 900 {
		t.Errorf("expected ExpiresIn 900, got %d", resp.ExpiresIn)
	}
	if resp.User == nil {
		t.Error("expected non-nil User")
	}
	if len(resp.Roles) != 1 {
		t.Errorf("expected 1 role, got %d", len(resp.Roles))
	}
}

// --- ValidateResponse tests ---

func TestValidateResponse_ValidToken(t *testing.T) {
	userID := "user-uuid-here"
	tenantID := "tenant-uuid-here"
	email := "user@example.com"

	resp := &ValidateResponse{
		Valid:     true,
		UserID:    &userID,
		TenantID:  &tenantID,
		Email:     &email,
		Roles:     []string{"user"},
		Permissions: map[string][]string{},
	}

	if !resp.Valid {
		t.Error("expected Valid to be true")
	}
	if resp.UserID == nil || *resp.UserID != userID {
		t.Errorf("expected UserID %s", userID)
	}
}

func TestValidateResponse_InvalidToken(t *testing.T) {
	resp := &ValidateResponse{
		Valid:   false,
		Message: "invalid or expired token",
	}

	if resp.Valid {
		t.Error("expected Valid to be false for invalid token")
	}
	if resp.Message == "" {
		t.Error("expected non-empty Message for invalid response")
	}
	if resp.UserID != nil {
		t.Error("expected nil UserID for invalid token response")
	}
}

// --- AuthEventFilter tests ---

func TestAuthEventFilter_DefaultValues(t *testing.T) {
	filter := &AuthEventFilter{}

	if filter.EventType != "" {
		t.Errorf("expected empty EventType, got %s", filter.EventType)
	}
	if filter.Email != "" {
		t.Errorf("expected empty Email, got %s", filter.Email)
	}
	if filter.UserID != nil {
		t.Error("expected nil UserID")
	}
	if filter.Limit != 0 {
		t.Errorf("expected Limit 0, got %d", filter.Limit)
	}
}

func TestAuthEvent_Fields(t *testing.T) {
	id := uuid.New()
	userID := uuid.New()
	now := time.Now()

	event := &AuthEvent{
		ID:        id,
		EventType: EventLoginSuccess,
		UserID:    &userID,
		Email:     "user@example.com",
		IPAddress: "127.0.0.1",
		UserAgent: "Mozilla/5.0",
		CreatedAt: now,
	}

	if event.ID != id {
		t.Errorf("expected ID %s, got %s", id, event.ID)
	}
	if event.EventType != EventLoginSuccess {
		t.Errorf("expected EventType %s, got %s", EventLoginSuccess, event.EventType)
	}
	if event.Email != "user@example.com" {
		t.Errorf("expected email 'user@example.com', got %s", event.Email)
	}
}

// --- Tenant tests ---

func TestTenant_Fields(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	tenant := &Tenant{
		ID:        id,
		Name:      "My Academy",
		Slug:      "my-academy",
		Plan:      "pro",
		Settings:  map[string]interface{}{"theme": "dark"},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if tenant.ID != id {
		t.Errorf("expected ID %s, got %s", id, tenant.ID)
	}
	if tenant.Slug != "my-academy" {
		t.Errorf("expected slug 'my-academy', got %s", tenant.Slug)
	}
	if tenant.Plan != "pro" {
		t.Errorf("expected plan 'pro', got %s", tenant.Plan)
	}
	theme, ok := tenant.Settings["theme"]
	if !ok || theme != "dark" {
		t.Error("expected Settings['theme'] = 'dark'")
	}
}

// --- Role tests ---

func TestRole_Fields(t *testing.T) {
	id := uuid.New()

	role := &Role{
		ID:       id,
		Name:     "admin",
		Level:    10,
		IsSystem: true,
	}

	if role.ID != id {
		t.Errorf("expected ID %s, got %s", id, role.ID)
	}
	if role.Name != "admin" {
		t.Errorf("expected name 'admin', got %s", role.Name)
	}
	if role.Level != 10 {
		t.Errorf("expected level 10, got %d", role.Level)
	}
	if !role.IsSystem {
		t.Error("expected IsSystem to be true")
	}
}

// --- Permission tests ---

func TestPermission_Fields(t *testing.T) {
	id := uuid.New()

	perm := &Permission{
		ID:       id,
		Name:     "read:books",
		Resource: "books",
		Action:   "read",
	}

	if perm.ID != id {
		t.Errorf("expected ID %s, got %s", id, perm.ID)
	}
	if perm.Resource != "books" {
		t.Errorf("expected resource 'books', got %s", perm.Resource)
	}
	if perm.Action != "read" {
		t.Errorf("expected action 'read', got %s", perm.Action)
	}
}

// --- AdminLoginRequest tests ---

func TestAdminLoginRequest_Fields(t *testing.T) {
	req := &AdminLoginRequest{
		Identifier: "admin@example.com",
		Password:   "securepassword",
	}

	if req.Identifier != "admin@example.com" {
		t.Errorf("expected identifier 'admin@example.com', got %s", req.Identifier)
	}
	if req.Password != "securepassword" {
		t.Errorf("expected password 'securepassword', got %s", req.Password)
	}
}
