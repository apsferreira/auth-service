package jwt

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

// --- helpers ---

func newTestService(accessExpiry, refreshExpiry time.Duration) *JWTService {
	return NewJWTService("test-secret-key-for-unit-tests-only", accessExpiry, refreshExpiry)
}

func defaultClaims() (userID, tenantID uuid.UUID, email string, roles []string, permissions map[string][]string) {
	userID = uuid.New()
	tenantID = uuid.New()
	email = "user@example.com"
	roles = []string{"admin", "user"}
	permissions = map[string][]string{
		"books":  {"read", "write"},
		"users":  {"read"},
	}
	return
}

// --- GenerateAccessToken tests ---

func TestGenerateAccessToken_ReturnsNonEmptyToken(t *testing.T) {
	svc := newTestService(15*time.Minute, 7*24*time.Hour)
	userID, tenantID, email, roles, permissions := defaultClaims()

	token, err := svc.GenerateAccessToken(userID, tenantID, email, roles, permissions)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}
}

func TestGenerateAccessToken_IsValidJWTFormat(t *testing.T) {
	svc := newTestService(15*time.Minute, 7*24*time.Hour)
	userID, tenantID, email, roles, permissions := defaultClaims()

	token, _ := svc.GenerateAccessToken(userID, tenantID, email, roles, permissions)

	// A JWT has exactly 3 parts separated by dots
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Errorf("expected JWT with 3 parts, got %d parts: %s", len(parts), token)
	}
}

func TestGenerateAccessToken_DifferentCallsProduceDifferentTokens(t *testing.T) {
	svc := newTestService(15*time.Minute, 7*24*time.Hour)
	userID, tenantID, email, roles, permissions := defaultClaims()

	token1, _ := svc.GenerateAccessToken(userID, tenantID, email, roles, permissions)
	token2, _ := svc.GenerateAccessToken(userID, tenantID, email, roles, permissions)

	// Each token has a unique JTI, so they should be different
	if token1 == token2 {
		t.Error("expected different tokens for different calls (unique JTI per call)")
	}
}

// --- ValidateAccessToken tests ---

func TestValidateAccessToken_ValidToken_ReturnsClaims(t *testing.T) {
	svc := newTestService(15*time.Minute, 7*24*time.Hour)
	userID, tenantID, email, roles, permissions := defaultClaims()

	token, err := svc.GenerateAccessToken(userID, tenantID, email, roles, permissions)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	claims, err := svc.ValidateAccessToken(token)

	if err != nil {
		t.Fatalf("expected no error for valid token, got %v", err)
	}
	if claims == nil {
		t.Fatal("expected non-nil claims")
	}
}

func TestValidateAccessToken_ClaimsContainCorrectSubject(t *testing.T) {
	svc := newTestService(15*time.Minute, 7*24*time.Hour)
	userID, tenantID, email, roles, permissions := defaultClaims()

	token, _ := svc.GenerateAccessToken(userID, tenantID, email, roles, permissions)
	claims, err := svc.ValidateAccessToken(token)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if claims.Subject != userID.String() {
		t.Errorf("expected subject %s, got %s", userID.String(), claims.Subject)
	}
}

func TestValidateAccessToken_ClaimsContainCorrectTenantID(t *testing.T) {
	svc := newTestService(15*time.Minute, 7*24*time.Hour)
	userID, tenantID, email, roles, permissions := defaultClaims()

	token, _ := svc.GenerateAccessToken(userID, tenantID, email, roles, permissions)
	claims, err := svc.ValidateAccessToken(token)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if claims.TenantID != tenantID.String() {
		t.Errorf("expected TenantID %s, got %s", tenantID.String(), claims.TenantID)
	}
}

func TestValidateAccessToken_ClaimsContainCorrectEmail(t *testing.T) {
	svc := newTestService(15*time.Minute, 7*24*time.Hour)
	userID, tenantID, email, roles, permissions := defaultClaims()

	token, _ := svc.GenerateAccessToken(userID, tenantID, email, roles, permissions)
	claims, err := svc.ValidateAccessToken(token)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if claims.Email != email {
		t.Errorf("expected email %s, got %s", email, claims.Email)
	}
}

func TestValidateAccessToken_ClaimsContainCorrectRoles(t *testing.T) {
	svc := newTestService(15*time.Minute, 7*24*time.Hour)
	userID, tenantID, email, roles, permissions := defaultClaims()

	token, _ := svc.GenerateAccessToken(userID, tenantID, email, roles, permissions)
	claims, err := svc.ValidateAccessToken(token)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(claims.Roles) != len(roles) {
		t.Errorf("expected %d roles, got %d", len(roles), len(claims.Roles))
	}
}

func TestValidateAccessToken_ClaimsContainPermissions(t *testing.T) {
	svc := newTestService(15*time.Minute, 7*24*time.Hour)
	userID, tenantID, email, roles, permissions := defaultClaims()

	token, _ := svc.GenerateAccessToken(userID, tenantID, email, roles, permissions)
	claims, err := svc.ValidateAccessToken(token)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(claims.Permissions) != len(permissions) {
		t.Errorf("expected %d permission keys, got %d", len(permissions), len(claims.Permissions))
	}
}

func TestValidateAccessToken_ExpiredToken_ReturnsError(t *testing.T) {
	// Use very short expiry (1 nanosecond)
	svc := newTestService(1*time.Nanosecond, 7*24*time.Hour)
	userID, tenantID, email, roles, permissions := defaultClaims()

	token, err := svc.GenerateAccessToken(userID, tenantID, email, roles, permissions)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	_, err = svc.ValidateAccessToken(token)

	if err == nil {
		t.Error("expected error for expired token")
	}
}

func TestValidateAccessToken_WrongSecret_ReturnsError(t *testing.T) {
	svc1 := NewJWTService("secret-one", 15*time.Minute, 7*24*time.Hour)
	svc2 := NewJWTService("secret-two", 15*time.Minute, 7*24*time.Hour)

	userID, tenantID, email, roles, permissions := defaultClaims()
	token, _ := svc1.GenerateAccessToken(userID, tenantID, email, roles, permissions)

	_, err := svc2.ValidateAccessToken(token)

	if err == nil {
		t.Error("expected error when validating with wrong secret")
	}
}

func TestValidateAccessToken_MalformedToken_ReturnsError(t *testing.T) {
	svc := newTestService(15*time.Minute, 7*24*time.Hour)

	_, err := svc.ValidateAccessToken("not.a.valid.jwt.token")

	if err == nil {
		t.Error("expected error for malformed token")
	}
}

func TestValidateAccessToken_EmptyToken_ReturnsError(t *testing.T) {
	svc := newTestService(15*time.Minute, 7*24*time.Hour)

	_, err := svc.ValidateAccessToken("")

	if err == nil {
		t.Error("expected error for empty token")
	}
}

func TestValidateAccessToken_TamperedToken_ReturnsError(t *testing.T) {
	svc := newTestService(15*time.Minute, 7*24*time.Hour)
	userID, tenantID, email, roles, permissions := defaultClaims()

	token, _ := svc.GenerateAccessToken(userID, tenantID, email, roles, permissions)

	// Tamper with the payload part
	parts := strings.Split(token, ".")
	parts[1] = parts[1] + "tampered"
	tamperedToken := strings.Join(parts, ".")

	_, err := svc.ValidateAccessToken(tamperedToken)

	if err == nil {
		t.Error("expected error for tampered token")
	}
}

// --- GenerateRefreshToken tests ---

func TestGenerateRefreshToken_ReturnsTokenAndHash(t *testing.T) {
	svc := newTestService(15*time.Minute, 7*24*time.Hour)

	token, hash, err := svc.GenerateRefreshToken()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}
	if hash == "" {
		t.Error("expected non-empty hash")
	}
}

func TestGenerateRefreshToken_TokenIsHexEncoded64Chars(t *testing.T) {
	svc := newTestService(15*time.Minute, 7*24*time.Hour)

	token, _, _ := svc.GenerateRefreshToken()

	// 32 bytes of random data → 64 hex chars
	if len(token) != 64 {
		t.Errorf("expected 64-char hex token, got %d chars", len(token))
	}
}

func TestGenerateRefreshToken_HashIsSHA256HexString(t *testing.T) {
	svc := newTestService(15*time.Minute, 7*24*time.Hour)

	_, hash, _ := svc.GenerateRefreshToken()

	// SHA-256 produces 32 bytes → 64 hex chars
	if len(hash) != 64 {
		t.Errorf("expected 64-char SHA-256 hash, got %d chars", len(hash))
	}
}

func TestGenerateRefreshToken_DifferentCallsProduceDifferentTokens(t *testing.T) {
	svc := newTestService(15*time.Minute, 7*24*time.Hour)

	token1, _, _ := svc.GenerateRefreshToken()
	token2, _, _ := svc.GenerateRefreshToken()

	if token1 == token2 {
		t.Error("expected different tokens on different calls")
	}
}

func TestGenerateRefreshToken_HashMatchesToken(t *testing.T) {
	svc := newTestService(15*time.Minute, 7*24*time.Hour)

	token, hash, _ := svc.GenerateRefreshToken()

	// The hash should equal HashToken(token)
	expectedHash := HashToken(token)
	if hash != expectedHash {
		t.Errorf("expected hash %s, got %s", expectedHash, hash)
	}
}

// --- HashToken tests ---

func TestHashToken_Deterministic_SameInputSameOutput(t *testing.T) {
	input := "my-test-token-value"

	hash1 := HashToken(input)
	hash2 := HashToken(input)

	if hash1 != hash2 {
		t.Errorf("expected same hash for same input, got %s and %s", hash1, hash2)
	}
}

func TestHashToken_DifferentInputsDifferentOutputs(t *testing.T) {
	hash1 := HashToken("token-one")
	hash2 := HashToken("token-two")

	if hash1 == hash2 {
		t.Error("expected different hashes for different inputs")
	}
}

func TestHashToken_ReturnsHexString(t *testing.T) {
	hash := HashToken("any-token")

	// SHA-256 output is 32 bytes = 64 hex chars
	if len(hash) != 64 {
		t.Errorf("expected 64-char hex SHA-256 hash, got %d chars: %s", len(hash), hash)
	}

	// Verify it's valid hex
	validHex := "0123456789abcdef"
	for _, c := range hash {
		if !strings.ContainsRune(validHex, c) {
			t.Errorf("expected lowercase hex string, found char '%c' in hash: %s", c, hash)
			break
		}
	}
}

func TestHashToken_EmptyStringProducesKnownSHA256(t *testing.T) {
	hash := HashToken("")
	// SHA-256("") = e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
	expected := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	if hash != expected {
		t.Errorf("expected SHA-256 of empty string to be %s, got %s", expected, hash)
	}
}

// --- GetAccessExpiry / GetRefreshExpiry tests ---

func TestGetAccessExpiry_ReturnsConfiguredValue(t *testing.T) {
	expectedExpiry := 30 * time.Minute
	svc := NewJWTService("secret", expectedExpiry, 7*24*time.Hour)

	if svc.GetAccessExpiry() != expectedExpiry {
		t.Errorf("expected access expiry %v, got %v", expectedExpiry, svc.GetAccessExpiry())
	}
}

func TestGetRefreshExpiry_ReturnsConfiguredValue(t *testing.T) {
	expectedExpiry := 14 * 24 * time.Hour
	svc := NewJWTService("secret", 15*time.Minute, expectedExpiry)

	if svc.GetRefreshExpiry() != expectedExpiry {
		t.Errorf("expected refresh expiry %v, got %v", expectedExpiry, svc.GetRefreshExpiry())
	}
}

// --- Integration: Generate then Validate ---

func TestGenerateAndValidate_RoundTrip(t *testing.T) {
	svc := newTestService(15*time.Minute, 7*24*time.Hour)

	userID := uuid.New()
	tenantID := uuid.New()
	email := "roundtrip@example.com"
	roles := []string{"editor"}
	permissions := map[string][]string{"articles": {"read", "write", "delete"}}

	token, err := svc.GenerateAccessToken(userID, tenantID, email, roles, permissions)
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	claims, err := svc.ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("validate failed: %v", err)
	}

	if claims.Subject != userID.String() {
		t.Errorf("userID mismatch: expected %s, got %s", userID.String(), claims.Subject)
	}
	if claims.TenantID != tenantID.String() {
		t.Errorf("tenantID mismatch: expected %s, got %s", tenantID.String(), claims.TenantID)
	}
	if claims.Email != email {
		t.Errorf("email mismatch: expected %s, got %s", email, claims.Email)
	}
	if len(claims.Roles) != 1 || claims.Roles[0] != "editor" {
		t.Errorf("roles mismatch: expected [editor], got %v", claims.Roles)
	}
	articlePerms, ok := claims.Permissions["articles"]
	if !ok || len(articlePerms) != 3 {
		t.Errorf("permissions mismatch: expected 3 article permissions, got %v", claims.Permissions)
	}
}
