package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/apsferreira/auth-service/backend/internal/client"
	"github.com/apsferreira/auth-service/backend/internal/domain"
	jwtpkg "github.com/apsferreira/auth-service/backend/internal/pkg/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// ─── Mocks ───────────────────────────────────────────────────────────────────

type mockUserRepo struct {
	findByEmailFn           func(email string) (*domain.User, error)
	findByIDFn              func(id uuid.UUID) (*domain.User, error)
	findByUsernameOrEmailFn func(identifier string) (*domain.User, error)
	createFn                func(user *domain.User) error
	updateFn                func(id uuid.UUID, req domain.UserUpdateRequest) (*domain.User, error)
	deleteFn                func(id uuid.UUID) error
	listFn                  func(tenantID uuid.UUID) ([]*domain.User, error)
	setPasswordFn           func(userID uuid.UUID, hash string) error
	updateLastLoginFn       func(userID uuid.UUID) error
	updateProfileFn         func(id uuid.UUID, fullName string) (*domain.User, error)
	addUserRoleFn           func(userID, roleID uuid.UUID) error
	getRolesAndPermsFn      func(userID uuid.UUID) ([]string, map[string][]string, error)
	getDefaultTenantIDFn    func() (uuid.UUID, error)
	getDefaultRoleIDFn      func(tenantID uuid.UUID) (uuid.UUID, error)
}

func (m *mockUserRepo) FindByEmail(email string) (*domain.User, error) {
	if m.findByEmailFn != nil {
		return m.findByEmailFn(email)
	}
	return nil, sql.ErrNoRows
}
func (m *mockUserRepo) FindByID(id uuid.UUID) (*domain.User, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(id)
	}
	return nil, errors.New("not found")
}
func (m *mockUserRepo) FindByUsernameOrEmail(identifier string) (*domain.User, error) {
	if m.findByUsernameOrEmailFn != nil {
		return m.findByUsernameOrEmailFn(identifier)
	}
	return nil, sql.ErrNoRows
}
func (m *mockUserRepo) Create(user *domain.User) error {
	if m.createFn != nil {
		return m.createFn(user)
	}
	return nil
}
func (m *mockUserRepo) Update(id uuid.UUID, req domain.UserUpdateRequest) (*domain.User, error) {
	if m.updateFn != nil {
		return m.updateFn(id, req)
	}
	return nil, nil
}
func (m *mockUserRepo) Delete(id uuid.UUID) error {
	if m.deleteFn != nil {
		return m.deleteFn(id)
	}
	return nil
}
func (m *mockUserRepo) List(tenantID uuid.UUID) ([]*domain.User, error) {
	if m.listFn != nil {
		return m.listFn(tenantID)
	}
	return nil, nil
}
func (m *mockUserRepo) SetPassword(userID uuid.UUID, hash string) error {
	if m.setPasswordFn != nil {
		return m.setPasswordFn(userID, hash)
	}
	return nil
}
func (m *mockUserRepo) UpdateLastLogin(userID uuid.UUID) error {
	if m.updateLastLoginFn != nil {
		return m.updateLastLoginFn(userID)
	}
	return nil
}
func (m *mockUserRepo) UpdateProfile(id uuid.UUID, fullName string) (*domain.User, error) {
	if m.updateProfileFn != nil {
		return m.updateProfileFn(id, fullName)
	}
	return nil, nil
}
func (m *mockUserRepo) AddUserRole(userID, roleID uuid.UUID) error {
	if m.addUserRoleFn != nil {
		return m.addUserRoleFn(userID, roleID)
	}
	return nil
}
func (m *mockUserRepo) GetUserRolesAndPermissions(userID uuid.UUID) ([]string, map[string][]string, error) {
	if m.getRolesAndPermsFn != nil {
		return m.getRolesAndPermsFn(userID)
	}
	return []string{"user"}, map[string][]string{}, nil
}
func (m *mockUserRepo) GetDefaultTenantID() (uuid.UUID, error) {
	if m.getDefaultTenantIDFn != nil {
		return m.getDefaultTenantIDFn()
	}
	return uuid.New(), nil
}
func (m *mockUserRepo) GetDefaultRoleID(tenantID uuid.UUID) (uuid.UUID, error) {
	if m.getDefaultRoleIDFn != nil {
		return m.getDefaultRoleIDFn(tenantID)
	}
	return uuid.New(), nil
}

type mockTokenRepo struct {
	createFn           func(token *domain.RefreshToken) error
	findByHashFn       func(tokenHash string) (*domain.RefreshToken, error)
	revokeFn           func(id uuid.UUID) error
	revokeAllForUserFn func(userID uuid.UUID) error
	deleteExpiredFn    func() error
}

func (m *mockTokenRepo) Create(token *domain.RefreshToken) error {
	if m.createFn != nil {
		return m.createFn(token)
	}
	return nil
}
func (m *mockTokenRepo) FindByHash(tokenHash string) (*domain.RefreshToken, error) {
	if m.findByHashFn != nil {
		return m.findByHashFn(tokenHash)
	}
	return nil, errors.New("not found")
}
func (m *mockTokenRepo) Revoke(id uuid.UUID) error {
	if m.revokeFn != nil {
		return m.revokeFn(id)
	}
	return nil
}
func (m *mockTokenRepo) RevokeAllForUser(userID uuid.UUID) error {
	if m.revokeAllForUserFn != nil {
		return m.revokeAllForUserFn(userID)
	}
	return nil
}
func (m *mockTokenRepo) DeleteExpired() error {
	if m.deleteExpiredFn != nil {
		return m.deleteExpiredFn()
	}
	return nil
}

type mockOAuthRepo struct {
	findByProviderFn func(provider, providerID string) (*domain.OAuthIdentity, error)
	upsertFn         func(identity *domain.OAuthIdentity) error
}

func (m *mockOAuthRepo) FindByProvider(provider, providerID string) (*domain.OAuthIdentity, error) {
	if m.findByProviderFn != nil {
		return m.findByProviderFn(provider, providerID)
	}
	return nil, errors.New("not found")
}
func (m *mockOAuthRepo) Upsert(identity *domain.OAuthIdentity) error {
	if m.upsertFn != nil {
		return m.upsertFn(identity)
	}
	return nil
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func newTestJWT() *jwtpkg.JWTService {
	return jwtpkg.NewJWTService("test-secret", 15*time.Minute, 7*24*time.Hour)
}

func newFixedUser() *domain.User {
	tenantID := uuid.New()
	roleID := uuid.New()
	return &domain.User{
		ID:       uuid.New(),
		TenantID: tenantID,
		Email:    "user@example.com",
		FullName: "Test User",
		IsActive: true,
		RoleID:   &roleID,
	}
}

// buildSvc wires up an AuthService without notificationClient and publisher
// so it never makes real HTTP calls. Tests that exercise OTP send paths
// must supply alternate wiring.
func buildSvc(ur *mockUserRepo, tr *mockTokenRepo, or_ *mockOAuthRepo) *AuthService {
	return NewAuthService(ur, nil, nil, tr, newTestJWT(), or_, nil, "Test Service")
}

// hashPassword generates a bcrypt hash at min cost for test speed.
func hashPassword(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	require.NoError(t, err)
	return string(hash)
}

// ─── NewAuthService ───────────────────────────────────────────────────────────

func TestNewAuthService_DefaultServiceName(t *testing.T) {
	svc := NewAuthService(&mockUserRepo{}, nil, nil, &mockTokenRepo{}, newTestJWT(), &mockOAuthRepo{}, nil, "")
	assert.Equal(t, "Instituto Itinerante", svc.serviceName)
}

func TestNewAuthService_CustomServiceName(t *testing.T) {
	svc := NewAuthService(&mockUserRepo{}, nil, nil, &mockTokenRepo{}, newTestJWT(), &mockOAuthRepo{}, nil, "My App")
	assert.Equal(t, "My App", svc.serviceName)
}

// ─── ValidateToken ───────────────────────────────────────────────────────────

func TestValidateToken_ValidToken_ReturnsValid(t *testing.T) {
	user := newFixedUser()
	svc := buildSvc(&mockUserRepo{}, &mockTokenRepo{}, &mockOAuthRepo{})

	token, err := newTestJWT().GenerateAccessToken(user.ID, user.TenantID, user.Email, []string{"user"}, nil)
	require.NoError(t, err)

	resp, err := svc.ValidateToken(token)
	require.NoError(t, err)
	assert.True(t, resp.Valid)
	assert.Equal(t, user.Email, *resp.Email)
}

func TestValidateToken_InvalidToken_ReturnsInvalid(t *testing.T) {
	svc := buildSvc(&mockUserRepo{}, &mockTokenRepo{}, &mockOAuthRepo{})

	resp, err := svc.ValidateToken("not-a-valid-token")
	require.NoError(t, err)
	assert.False(t, resp.Valid)
	assert.NotEmpty(t, resp.Message)
}

func TestValidateToken_EmptyToken_ReturnsInvalid(t *testing.T) {
	svc := buildSvc(&mockUserRepo{}, &mockTokenRepo{}, &mockOAuthRepo{})

	resp, err := svc.ValidateToken("")
	require.NoError(t, err)
	assert.False(t, resp.Valid)
}

func TestValidateToken_WrongSecretToken_ReturnsInvalid(t *testing.T) {
	// generate with different secret
	otherJWT := jwtpkg.NewJWTService("other-secret", 15*time.Minute, 7*24*time.Hour)
	user := newFixedUser()
	token, err := otherJWT.GenerateAccessToken(user.ID, user.TenantID, user.Email, nil, nil)
	require.NoError(t, err)

	svc := buildSvc(&mockUserRepo{}, &mockTokenRepo{}, &mockOAuthRepo{})
	resp, err := svc.ValidateToken(token)
	require.NoError(t, err)
	assert.False(t, resp.Valid)
}

// ─── GetCurrentUser ───────────────────────────────────────────────────────────

func TestGetCurrentUser_UserExists_ReturnsUserWithRoles(t *testing.T) {
	user := newFixedUser()
	ur := &mockUserRepo{
		findByIDFn: func(id uuid.UUID) (*domain.User, error) { return user, nil },
		getRolesAndPermsFn: func(userID uuid.UUID) ([]string, map[string][]string, error) {
			return []string{"admin", "editor"}, nil, nil
		},
	}
	svc := buildSvc(ur, &mockTokenRepo{}, &mockOAuthRepo{})

	result, err := svc.GetCurrentUser(user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.Email, result.Email)
	assert.Len(t, result.Roles, 2)
}

func TestGetCurrentUser_UserNotFound_ReturnsError(t *testing.T) {
	ur := &mockUserRepo{
		findByIDFn: func(id uuid.UUID) (*domain.User, error) { return nil, errors.New("not found") },
	}
	svc := buildSvc(ur, &mockTokenRepo{}, &mockOAuthRepo{})

	_, err := svc.GetCurrentUser(uuid.New())
	assert.Error(t, err)
}

func TestGetCurrentUser_RolesError_ReturnsUserWithoutRoles(t *testing.T) {
	user := newFixedUser()
	ur := &mockUserRepo{
		findByIDFn: func(id uuid.UUID) (*domain.User, error) { return user, nil },
		getRolesAndPermsFn: func(userID uuid.UUID) ([]string, map[string][]string, error) {
			return nil, nil, errors.New("db error")
		},
	}
	svc := buildSvc(ur, &mockTokenRepo{}, &mockOAuthRepo{})

	result, err := svc.GetCurrentUser(user.ID)
	require.NoError(t, err)
	assert.Empty(t, result.Roles)
}

// ─── UpdateCurrentUser ────────────────────────────────────────────────────────

func TestUpdateCurrentUser_Success(t *testing.T) {
	user := newFixedUser()
	user.FullName = "New Name"
	ur := &mockUserRepo{
		updateProfileFn: func(id uuid.UUID, fullName string) (*domain.User, error) { return user, nil },
	}
	svc := buildSvc(ur, &mockTokenRepo{}, &mockOAuthRepo{})

	result, err := svc.UpdateCurrentUser(user.ID, "New Name")
	require.NoError(t, err)
	assert.Equal(t, "New Name", result.FullName)
}

// ─── Logout ──────────────────────────────────────────────────────────────────

func TestLogout_ValidToken_Succeeds(t *testing.T) {
	stored := &domain.RefreshToken{ID: uuid.New(), UserID: uuid.New(), ExpiresAt: time.Now().Add(time.Hour)}
	tr := &mockTokenRepo{
		findByHashFn: func(hash string) (*domain.RefreshToken, error) { return stored, nil },
		revokeFn:     func(id uuid.UUID) error { return nil },
	}
	svc := buildSvc(&mockUserRepo{}, tr, &mockOAuthRepo{})

	err := svc.Logout("some-refresh-token")
	assert.NoError(t, err)
}

func TestLogout_InvalidToken_ReturnsError(t *testing.T) {
	tr := &mockTokenRepo{
		findByHashFn: func(hash string) (*domain.RefreshToken, error) { return nil, errors.New("not found") },
	}
	svc := buildSvc(&mockUserRepo{}, tr, &mockOAuthRepo{})

	err := svc.Logout("bad-token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid refresh token")
}

// ─── RefreshTokens ────────────────────────────────────────────────────────────

func TestRefreshTokens_ValidToken_ReturnsNewPair(t *testing.T) {
	user := newFixedUser()
	stored := &domain.RefreshToken{ID: uuid.New(), UserID: user.ID, ExpiresAt: time.Now().Add(time.Hour)}

	ur := &mockUserRepo{
		findByIDFn:         func(id uuid.UUID) (*domain.User, error) { return user, nil },
		getRolesAndPermsFn: func(userID uuid.UUID) ([]string, map[string][]string, error) { return []string{"user"}, nil, nil },
		updateLastLoginFn:  func(userID uuid.UUID) error { return nil },
	}
	tr := &mockTokenRepo{
		findByHashFn: func(hash string) (*domain.RefreshToken, error) { return stored, nil },
		revokeFn:     func(id uuid.UUID) error { return nil },
		createFn:     func(token *domain.RefreshToken) error { return nil },
	}
	svc := buildSvc(ur, tr, &mockOAuthRepo{})

	resp, err := svc.RefreshTokens("old-token")
	require.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
}

func TestRefreshTokens_InvalidToken_ReturnsError(t *testing.T) {
	tr := &mockTokenRepo{
		findByHashFn: func(hash string) (*domain.RefreshToken, error) { return nil, errors.New("not found") },
	}
	svc := buildSvc(&mockUserRepo{}, tr, &mockOAuthRepo{})

	_, err := svc.RefreshTokens("bad-token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid or expired refresh token")
}

func TestRefreshTokens_RevokeError_ReturnsError(t *testing.T) {
	stored := &domain.RefreshToken{ID: uuid.New(), UserID: uuid.New(), ExpiresAt: time.Now().Add(time.Hour)}
	tr := &mockTokenRepo{
		findByHashFn: func(hash string) (*domain.RefreshToken, error) { return stored, nil },
		revokeFn:     func(id uuid.UUID) error { return errors.New("db error") },
	}
	svc := buildSvc(&mockUserRepo{}, tr, &mockOAuthRepo{})

	_, err := svc.RefreshTokens("any-token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to revoke old token")
}

func TestRefreshTokens_UserNotFound_ReturnsError(t *testing.T) {
	stored := &domain.RefreshToken{ID: uuid.New(), UserID: uuid.New(), ExpiresAt: time.Now().Add(time.Hour)}
	tr := &mockTokenRepo{
		findByHashFn: func(hash string) (*domain.RefreshToken, error) { return stored, nil },
		revokeFn:     func(id uuid.UUID) error { return nil },
	}
	ur := &mockUserRepo{
		findByIDFn: func(id uuid.UUID) (*domain.User, error) { return nil, errors.New("gone") },
	}
	svc := buildSvc(ur, tr, &mockOAuthRepo{})

	_, err := svc.RefreshTokens("token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestRefreshTokens_TokenCreateError_ReturnsError(t *testing.T) {
	user := newFixedUser()
	stored := &domain.RefreshToken{ID: uuid.New(), UserID: user.ID, ExpiresAt: time.Now().Add(time.Hour)}

	ur := &mockUserRepo{
		findByIDFn:         func(id uuid.UUID) (*domain.User, error) { return user, nil },
		getRolesAndPermsFn: func(userID uuid.UUID) ([]string, map[string][]string, error) { return []string{"user"}, nil, nil },
	}
	tr := &mockTokenRepo{
		findByHashFn: func(hash string) (*domain.RefreshToken, error) { return stored, nil },
		revokeFn:     func(id uuid.UUID) error { return nil },
		createFn:     func(token *domain.RefreshToken) error { return errors.New("store error") },
	}
	svc := buildSvc(ur, tr, &mockOAuthRepo{})

	_, err := svc.RefreshTokens("token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to store refresh token")
}

// ─── AdminLogin ───────────────────────────────────────────────────────────────

func TestAdminLogin_UserNotFound_ReturnsInvalidCredentials(t *testing.T) {
	ur := &mockUserRepo{
		findByUsernameOrEmailFn: func(id string) (*domain.User, error) { return nil, sql.ErrNoRows },
	}
	svc := buildSvc(ur, &mockTokenRepo{}, &mockOAuthRepo{})

	_, err := svc.AdminLogin("admin@example.com", "secret")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestAdminLogin_DBError_ReturnsError(t *testing.T) {
	ur := &mockUserRepo{
		findByUsernameOrEmailFn: func(id string) (*domain.User, error) { return nil, errors.New("db error") },
	}
	svc := buildSvc(ur, &mockTokenRepo{}, &mockOAuthRepo{})

	_, err := svc.AdminLogin("admin@example.com", "secret")
	assert.Error(t, err)
}

func TestAdminLogin_NoPasswordHash_ReturnsError(t *testing.T) {
	user := newFixedUser()
	user.PasswordHash = nil
	ur := &mockUserRepo{
		findByUsernameOrEmailFn: func(id string) (*domain.User, error) { return user, nil },
	}
	svc := buildSvc(ur, &mockTokenRepo{}, &mockOAuthRepo{})

	_, err := svc.AdminLogin("admin@example.com", "pass")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "password authentication not configured")
}

func TestAdminLogin_WrongPassword_ReturnsInvalidCredentials(t *testing.T) {
	user := newFixedUser()
	hash := hashPassword(t, "correct-password")
	user.PasswordHash = &hash
	ur := &mockUserRepo{
		findByUsernameOrEmailFn: func(id string) (*domain.User, error) { return user, nil },
	}
	svc := buildSvc(ur, &mockTokenRepo{}, &mockOAuthRepo{})

	_, err := svc.AdminLogin("admin@example.com", "wrong-password")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestAdminLogin_NonAdminRole_ReturnsInsufficientPrivileges(t *testing.T) {
	user := newFixedUser()
	hash := hashPassword(t, "pass")
	user.PasswordHash = &hash
	ur := &mockUserRepo{
		findByUsernameOrEmailFn: func(id string) (*domain.User, error) { return user, nil },
		getRolesAndPermsFn:      func(userID uuid.UUID) ([]string, map[string][]string, error) { return []string{"user"}, nil, nil },
	}
	svc := buildSvc(ur, &mockTokenRepo{}, &mockOAuthRepo{})

	_, err := svc.AdminLogin("admin@example.com", "pass")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient privileges")
}

func TestAdminLogin_AdminRole_ReturnsTokens(t *testing.T) {
	user := newFixedUser()
	hash := hashPassword(t, "adminpass")
	user.PasswordHash = &hash

	ur := &mockUserRepo{
		findByUsernameOrEmailFn: func(id string) (*domain.User, error) { return user, nil },
		getRolesAndPermsFn:      func(userID uuid.UUID) ([]string, map[string][]string, error) { return []string{"admin"}, nil, nil },
		updateLastLoginFn:       func(userID uuid.UUID) error { return nil },
	}
	tr := &mockTokenRepo{createFn: func(token *domain.RefreshToken) error { return nil }}
	svc := buildSvc(ur, tr, &mockOAuthRepo{})

	resp, err := svc.AdminLogin("admin@example.com", "adminpass")
	require.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
}

func TestAdminLogin_SuperAdminRole_ReturnsTokens(t *testing.T) {
	user := newFixedUser()
	hash := hashPassword(t, "superpass")
	user.PasswordHash = &hash

	ur := &mockUserRepo{
		findByUsernameOrEmailFn: func(id string) (*domain.User, error) { return user, nil },
		getRolesAndPermsFn:      func(userID uuid.UUID) ([]string, map[string][]string, error) { return []string{"super_admin"}, nil, nil },
		updateLastLoginFn:       func(userID uuid.UUID) error { return nil },
	}
	tr := &mockTokenRepo{createFn: func(token *domain.RefreshToken) error { return nil }}
	svc := buildSvc(ur, tr, &mockOAuthRepo{})

	resp, err := svc.AdminLogin("super@example.com", "superpass")
	require.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
}

// ─── createDefaultUser ────────────────────────────────────────────────────────

func TestCreateDefaultUser_NoTenant_ReturnsError(t *testing.T) {
	ur := &mockUserRepo{
		getDefaultTenantIDFn: func() (uuid.UUID, error) { return uuid.Nil, errors.New("no tenant") },
	}
	svc := buildSvc(ur, &mockTokenRepo{}, &mockOAuthRepo{})

	err := svc.createDefaultUser("new@example.com")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no default tenant found")
}

func TestCreateDefaultUser_NoRole_ReturnsError(t *testing.T) {
	ur := &mockUserRepo{
		getDefaultTenantIDFn: func() (uuid.UUID, error) { return uuid.New(), nil },
		getDefaultRoleIDFn:   func(tenantID uuid.UUID) (uuid.UUID, error) { return uuid.Nil, errors.New("no role") },
	}
	svc := buildSvc(ur, &mockTokenRepo{}, &mockOAuthRepo{})

	err := svc.createDefaultUser("new@example.com")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no default role found")
}

func TestCreateDefaultUser_CreateFails_ReturnsError(t *testing.T) {
	ur := &mockUserRepo{
		getDefaultTenantIDFn: func() (uuid.UUID, error) { return uuid.New(), nil },
		getDefaultRoleIDFn:   func(tenantID uuid.UUID) (uuid.UUID, error) { return uuid.New(), nil },
		createFn:             func(user *domain.User) error { return errors.New("db error") },
	}
	svc := buildSvc(ur, &mockTokenRepo{}, &mockOAuthRepo{})

	err := svc.createDefaultUser("fail@example.com")
	assert.Error(t, err)
}

func TestCreateDefaultUser_Success(t *testing.T) {
	var createdUser *domain.User
	ur := &mockUserRepo{
		getDefaultTenantIDFn: func() (uuid.UUID, error) { return uuid.New(), nil },
		getDefaultRoleIDFn:   func(tenantID uuid.UUID) (uuid.UUID, error) { return uuid.New(), nil },
		createFn: func(user *domain.User) error {
			createdUser = user
			return nil
		},
		addUserRoleFn: func(userID, roleID uuid.UUID) error { return nil },
	}
	svc := buildSvc(ur, &mockTokenRepo{}, &mockOAuthRepo{})

	err := svc.createDefaultUser("new@example.com")
	require.NoError(t, err)
	require.NotNil(t, createdUser)
	assert.Equal(t, "new@example.com", createdUser.Email)
	assert.True(t, createdUser.IsActive)
	assert.Empty(t, createdUser.FullName)
}

// ─── RequestOTP (paths that don't need real HTTP) ────────────────────────────

func TestRequestOTP_WhatsAppNotConfigured_ReturnsError(t *testing.T) {
	user := newFixedUser()
	ur := &mockUserRepo{
		findByEmailFn: func(email string) (*domain.User, error) { return user, nil },
	}
	svc := buildSvc(ur, &mockTokenRepo{}, &mockOAuthRepo{})
	// whatsappService is nil → not configured

	_, _, err := svc.RequestOTP("user@example.com", "whatsapp")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "não está configurado")
}

func TestRequestOTP_UserLookupDBError_ReturnsError(t *testing.T) {
	ur := &mockUserRepo{
		findByEmailFn: func(email string) (*domain.User, error) { return nil, errors.New("db connection error") },
	}
	svc := buildSvc(ur, &mockTokenRepo{}, &mockOAuthRepo{})

	_, _, err := svc.RequestOTP("user@example.com", "email")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find user")
}

func TestRequestOTP_NewUserCreateFails_ReturnsError(t *testing.T) {
	ur := &mockUserRepo{
		findByEmailFn:        func(email string) (*domain.User, error) { return nil, sql.ErrNoRows },
		getDefaultTenantIDFn: func() (uuid.UUID, error) { return uuid.New(), nil },
		getDefaultRoleIDFn:   func(tenantID uuid.UUID) (uuid.UUID, error) { return uuid.New(), nil },
		createFn:             func(user *domain.User) error { return errors.New("insert failed") },
	}
	svc := buildSvc(ur, &mockTokenRepo{}, &mockOAuthRepo{})

	_, _, err := svc.RequestOTP("new@example.com", "email")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create user")
}

func TestRequestOTP_EmailChannelDefault(t *testing.T) {
	// channel == "" should default to "email" — user lookup errors out before HTTP call
	ur := &mockUserRepo{
		findByEmailFn: func(email string) (*domain.User, error) { return nil, errors.New("db error") },
	}
	svc := buildSvc(ur, &mockTokenRepo{}, &mockOAuthRepo{})

	_, _, err := svc.RequestOTP("user@example.com", "")
	// The error here is from userRepo, which means the channel defaulted and execution proceeded past channel check
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find user")
}

// ─── ProvisionUser ────────────────────────────────────────────────────────────

func TestProvisionUser_ExistingUserWithName_ReturnsExisting(t *testing.T) {
	user := newFixedUser()
	user.FullName = "Already Set"
	ur := &mockUserRepo{
		findByEmailFn: func(email string) (*domain.User, error) { return user, nil },
	}
	svc := buildSvc(ur, &mockTokenRepo{}, &mockOAuthRepo{})

	result, err := svc.ProvisionUser(user.Email, "Any Name")
	require.NoError(t, err)
	assert.Equal(t, user.ID, result.ID)
}

func TestProvisionUser_ExistingUserEmptyName_UpdatesProfile(t *testing.T) {
	user := newFixedUser()
	user.FullName = ""
	updated := newFixedUser()
	updated.ID = user.ID
	updated.FullName = "New Name"

	ur := &mockUserRepo{
		findByEmailFn:   func(email string) (*domain.User, error) { return user, nil },
		updateProfileFn: func(id uuid.UUID, fullName string) (*domain.User, error) { return updated, nil },
	}
	svc := buildSvc(ur, &mockTokenRepo{}, &mockOAuthRepo{})

	result, err := svc.ProvisionUser(user.Email, "New Name")
	require.NoError(t, err)
	assert.Equal(t, "New Name", result.FullName)
}

func TestProvisionUser_NewUser_CreatesAccount(t *testing.T) {
	ur := &mockUserRepo{
		findByEmailFn:        func(email string) (*domain.User, error) { return nil, errors.New("not found") },
		getDefaultTenantIDFn: func() (uuid.UUID, error) { return uuid.New(), nil },
		getDefaultRoleIDFn:   func(tenantID uuid.UUID) (uuid.UUID, error) { return uuid.New(), nil },
		createFn:             func(user *domain.User) error { return nil },
		addUserRoleFn:        func(userID, roleID uuid.UUID) error { return nil },
	}
	svc := buildSvc(ur, &mockTokenRepo{}, &mockOAuthRepo{})

	result, err := svc.ProvisionUser("brand-new@example.com", "Brand New")
	require.NoError(t, err)
	assert.Equal(t, "brand-new@example.com", result.Email)
	assert.Equal(t, "Brand New", result.FullName)
}

func TestProvisionUser_NewUser_NoTenant_ReturnsError(t *testing.T) {
	ur := &mockUserRepo{
		findByEmailFn:        func(email string) (*domain.User, error) { return nil, errors.New("not found") },
		getDefaultTenantIDFn: func() (uuid.UUID, error) { return uuid.Nil, errors.New("no tenant") },
	}
	svc := buildSvc(ur, &mockTokenRepo{}, &mockOAuthRepo{})

	_, err := svc.ProvisionUser("x@example.com", "X")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no default tenant found")
}

func TestProvisionUser_NewUser_NoRole_ReturnsError(t *testing.T) {
	ur := &mockUserRepo{
		findByEmailFn:        func(email string) (*domain.User, error) { return nil, errors.New("not found") },
		getDefaultTenantIDFn: func() (uuid.UUID, error) { return uuid.New(), nil },
		getDefaultRoleIDFn:   func(tenantID uuid.UUID) (uuid.UUID, error) { return uuid.Nil, errors.New("no role") },
	}
	svc := buildSvc(ur, &mockTokenRepo{}, &mockOAuthRepo{})

	_, err := svc.ProvisionUser("x@example.com", "X")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no default role found")
}

// ─── LoginWithGoogle ──────────────────────────────────────────────────────────

func TestLoginWithGoogle_ExistingIdentity_ReturnsTokens(t *testing.T) {
	user := newFixedUser()
	identity := &domain.OAuthIdentity{UserID: user.ID, Provider: "google", ProviderID: "g-123"}

	or_ := &mockOAuthRepo{
		findByProviderFn: func(provider, providerID string) (*domain.OAuthIdentity, error) { return identity, nil },
	}
	ur := &mockUserRepo{
		findByIDFn:         func(id uuid.UUID) (*domain.User, error) { return user, nil },
		getRolesAndPermsFn: func(userID uuid.UUID) ([]string, map[string][]string, error) { return []string{"user"}, nil, nil },
		updateLastLoginFn:  func(userID uuid.UUID) error { return nil },
		updateFn:           func(id uuid.UUID, req domain.UserUpdateRequest) (*domain.User, error) { return user, nil },
	}
	tr := &mockTokenRepo{createFn: func(token *domain.RefreshToken) error { return nil }}
	svc := buildSvc(ur, tr, or_)

	resp, err := svc.LoginWithGoogle(&GoogleUserInfo{ID: "g-123", Email: user.Email, Name: user.FullName})
	require.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
}

func TestLoginWithGoogle_LinkedUserMissing_ReturnsError(t *testing.T) {
	identity := &domain.OAuthIdentity{UserID: uuid.New(), Provider: "google", ProviderID: "g-456"}
	or_ := &mockOAuthRepo{
		findByProviderFn: func(provider, providerID string) (*domain.OAuthIdentity, error) { return identity, nil },
	}
	ur := &mockUserRepo{
		findByIDFn: func(id uuid.UUID) (*domain.User, error) { return nil, errors.New("deleted") },
	}
	svc := buildSvc(ur, &mockTokenRepo{}, or_)

	_, err := svc.LoginWithGoogle(&GoogleUserInfo{ID: "g-456", Email: "x@example.com"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "linked user not found")
}

func TestLoginWithGoogle_NewUser_CreatesAccount(t *testing.T) {
	newUser := newFixedUser()
	newUser.FullName = ""
	newUser.AvatarURL = nil

	callCount := 0
	or_ := &mockOAuthRepo{
		findByProviderFn: func(provider, providerID string) (*domain.OAuthIdentity, error) { return nil, errors.New("not found") },
		upsertFn:         func(identity *domain.OAuthIdentity) error { return nil },
	}
	ur := &mockUserRepo{
		findByEmailFn: func(email string) (*domain.User, error) {
			callCount++
			if callCount == 1 {
				return nil, sql.ErrNoRows
			}
			return newUser, nil
		},
		getDefaultTenantIDFn: func() (uuid.UUID, error) { return newUser.TenantID, nil },
		getDefaultRoleIDFn:   func(tenantID uuid.UUID) (uuid.UUID, error) { return *newUser.RoleID, nil },
		createFn:             func(user *domain.User) error { return nil },
		addUserRoleFn:        func(userID, roleID uuid.UUID) error { return nil },
		findByIDFn:           func(id uuid.UUID) (*domain.User, error) { return newUser, nil },
		getRolesAndPermsFn:   func(userID uuid.UUID) ([]string, map[string][]string, error) { return []string{"user"}, nil, nil },
		updateLastLoginFn:    func(userID uuid.UUID) error { return nil },
		updateFn:             func(id uuid.UUID, req domain.UserUpdateRequest) (*domain.User, error) { return newUser, nil },
	}
	tr := &mockTokenRepo{createFn: func(token *domain.RefreshToken) error { return nil }}
	svc := buildSvc(ur, tr, or_)

	resp, err := svc.LoginWithGoogle(&GoogleUserInfo{
		ID:      "g-new",
		Email:   "new@example.com",
		Name:    "New User",
		Picture: "https://example.com/pic.jpg",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
}

func TestLoginWithGoogle_ExistingEmailUser_LinksIdentity(t *testing.T) {
	// Google identity not found, but user already exists by email → link should be upserted
	user := newFixedUser()
	upsertCalled := false
	or_ := &mockOAuthRepo{
		findByProviderFn: func(provider, providerID string) (*domain.OAuthIdentity, error) { return nil, errors.New("not found") },
		upsertFn: func(identity *domain.OAuthIdentity) error {
			upsertCalled = true
			return nil
		},
	}
	ur := &mockUserRepo{
		findByEmailFn:      func(email string) (*domain.User, error) { return user, nil },
		getRolesAndPermsFn: func(userID uuid.UUID) ([]string, map[string][]string, error) { return []string{"user"}, nil, nil },
		updateLastLoginFn:  func(userID uuid.UUID) error { return nil },
		updateFn:           func(id uuid.UUID, req domain.UserUpdateRequest) (*domain.User, error) { return user, nil },
	}
	tr := &mockTokenRepo{createFn: func(token *domain.RefreshToken) error { return nil }}
	svc := buildSvc(ur, tr, or_)

	resp, err := svc.LoginWithGoogle(&GoogleUserInfo{ID: "g-existing", Email: user.Email, Name: user.FullName})
	require.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
	assert.True(t, upsertCalled, "expected OAuth identity upsert to be called")
}

// ─── generateOTPCode ─────────────────────────────────────────────────────────

func TestGenerateOTPCode_AlwaysSixDigits(t *testing.T) {
	for i := 0; i < 20; i++ {
		code, err := generateOTPCode()
		require.NoError(t, err)
		assert.Len(t, code, 6, "expected 6-digit code, got %q", code)
	}
}

func TestGenerateOTPCode_OnlyDigits(t *testing.T) {
	code, err := generateOTPCode()
	require.NoError(t, err)
	for _, c := range code {
		assert.True(t, c >= '0' && c <= '9', "unexpected character %c in code %s", c, code)
	}
}

// ─── VerifyOTP (requires mock notification server) ───────────────────────────

// buildSvcWithNotifServer creates an AuthService backed by a mock HTTP notification server.
func buildSvcWithNotifServer(
	t *testing.T,
	server *httptest.Server,
	ur *mockUserRepo,
	tr *mockTokenRepo,
) *AuthService {
	t.Helper()
	nc := client.NewNotificationClient(server.URL)
	return NewAuthService(ur, nc, nil, tr, newTestJWT(), &mockOAuthRepo{}, nil, "Test Service")
}

func TestVerifyOTP_InvalidCode_Returns401(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"valid":   false,
			"message": "invalid or expired OTP",
		})
	}))
	defer server.Close()

	svc := buildSvcWithNotifServer(t, server, &mockUserRepo{}, &mockTokenRepo{})

	_, err := svc.VerifyOTP("user@example.com", "000000")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid or expired OTP")
}

func TestVerifyOTP_NotificationServerDown_ReturnsError(t *testing.T) {
	// server closed immediately — simulates notification service unavailable
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Close()

	svc := buildSvcWithNotifServer(t, server, &mockUserRepo{}, &mockTokenRepo{})

	_, err := svc.VerifyOTP("user@example.com", "123456")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to verify OTP")
}

func TestVerifyOTP_ValidCode_UserNotFound_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"valid":   true,
			"message": "OTP verified",
		})
	}))
	defer server.Close()

	ur := &mockUserRepo{
		findByEmailFn: func(email string) (*domain.User, error) {
			return nil, errors.New("user not found")
		},
	}
	svc := buildSvcWithNotifServer(t, server, ur, &mockTokenRepo{})

	_, err := svc.VerifyOTP("unknown@example.com", "123456")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestVerifyOTP_ValidCode_Success_ReturnsAuthResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"valid":   true,
			"message": "OTP verified",
		})
	}))
	defer server.Close()

	user := newFixedUser()
	ur := &mockUserRepo{
		findByEmailFn: func(email string) (*domain.User, error) { return user, nil },
		updateLastLoginFn: func(userID uuid.UUID) error { return nil },
		getRolesAndPermsFn: func(userID uuid.UUID) ([]string, map[string][]string, error) {
			return []string{"user"}, nil, nil
		},
	}
	tr := &mockTokenRepo{
		createFn: func(token *domain.RefreshToken) error { return nil },
	}
	svc := buildSvcWithNotifServer(t, server, ur, tr)

	resp, err := svc.VerifyOTP(user.Email, "123456")
	require.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	assert.Equal(t, user.Email, resp.User.Email)
}

func TestVerifyOTP_ValidCode_TokenStoreFails_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"valid": true, "message": "ok"})
	}))
	defer server.Close()

	user := newFixedUser()
	ur := &mockUserRepo{
		findByEmailFn:     func(email string) (*domain.User, error) { return user, nil },
		updateLastLoginFn: func(userID uuid.UUID) error { return nil },
		getRolesAndPermsFn: func(userID uuid.UUID) ([]string, map[string][]string, error) {
			return []string{"user"}, nil, nil
		},
	}
	tr := &mockTokenRepo{
		createFn: func(token *domain.RefreshToken) error { return errors.New("store failed") },
	}
	svc := buildSvcWithNotifServer(t, server, ur, tr)

	_, err := svc.VerifyOTP(user.Email, "123456")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to store refresh token")
}
