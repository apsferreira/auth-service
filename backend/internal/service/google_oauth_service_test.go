package service

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ─── IsConfigured ────────────────────────────────────────────────────────────

func TestGoogleOAuthService_IsConfigured_AllFieldsPresent_ReturnsTrue(t *testing.T) {
	svc := NewGoogleOAuthService("client-id", "client-secret", "https://example.com/callback")
	assert.True(t, svc.IsConfigured())
}

func TestGoogleOAuthService_IsConfigured_MissingClientID_ReturnsFalse(t *testing.T) {
	svc := NewGoogleOAuthService("", "client-secret", "https://example.com/callback")
	assert.False(t, svc.IsConfigured())
}

func TestGoogleOAuthService_IsConfigured_MissingClientSecret_ReturnsFalse(t *testing.T) {
	svc := NewGoogleOAuthService("client-id", "", "https://example.com/callback")
	assert.False(t, svc.IsConfigured())
}

func TestGoogleOAuthService_IsConfigured_MissingRedirectURI_ReturnsFalse(t *testing.T) {
	svc := NewGoogleOAuthService("client-id", "client-secret", "")
	assert.False(t, svc.IsConfigured())
}

func TestGoogleOAuthService_IsConfigured_AllEmpty_ReturnsFalse(t *testing.T) {
	svc := NewGoogleOAuthService("", "", "")
	assert.False(t, svc.IsConfigured())
}

// ─── GetAuthURL ──────────────────────────────────────────────────────────────

func TestGoogleOAuthService_GetAuthURL_ContainsGoogleEndpoint(t *testing.T) {
	svc := NewGoogleOAuthService("my-client-id", "my-secret", "https://myapp.com/callback")
	url := svc.GetAuthURL("my-state-token")

	assert.Contains(t, url, "https://accounts.google.com/o/oauth2/v2/auth")
}

func TestGoogleOAuthService_GetAuthURL_ContainsClientID(t *testing.T) {
	svc := NewGoogleOAuthService("my-client-id", "my-secret", "https://myapp.com/callback")
	url := svc.GetAuthURL("some-state")

	assert.Contains(t, url, "my-client-id")
}

func TestGoogleOAuthService_GetAuthURL_ContainsRedirectURI(t *testing.T) {
	svc := NewGoogleOAuthService("cid", "csec", "https://myapp.com/oauth/callback")
	url := svc.GetAuthURL("state-xyz")

	assert.Contains(t, url, "myapp.com")
}

func TestGoogleOAuthService_GetAuthURL_ContainsState(t *testing.T) {
	svc := NewGoogleOAuthService("cid", "csec", "https://myapp.com/callback")
	state := "unique-state-value"
	url := svc.GetAuthURL(state)

	assert.Contains(t, url, state)
}

func TestGoogleOAuthService_GetAuthURL_ContainsRequiredScopes(t *testing.T) {
	svc := NewGoogleOAuthService("cid", "csec", "https://myapp.com/callback")
	url := svc.GetAuthURL("state")

	// The URL must include openid, email, profile scopes (URL-encoded)
	assert.Contains(t, url, "openid")
	assert.Contains(t, url, "email")
	assert.Contains(t, url, "profile")
}

func TestGoogleOAuthService_GetAuthURL_ContainsResponseTypeCode(t *testing.T) {
	svc := NewGoogleOAuthService("cid", "csec", "https://myapp.com/callback")
	url := svc.GetAuthURL("state")

	assert.Contains(t, url, "response_type=code")
}

func TestGoogleOAuthService_GetAuthURL_DifferentStates_ProduceDifferentURLs(t *testing.T) {
	svc := NewGoogleOAuthService("cid", "csec", "https://myapp.com/callback")
	url1 := svc.GetAuthURL("state-aaa")
	url2 := svc.GetAuthURL("state-bbb")

	assert.NotEqual(t, url1, url2)
}

func TestGoogleOAuthService_GetAuthURL_IsValidURL(t *testing.T) {
	svc := NewGoogleOAuthService("cid", "csec", "https://myapp.com/callback")
	url := svc.GetAuthURL("state")

	assert.True(t, strings.HasPrefix(url, "https://"), "URL should use HTTPS")
	assert.Contains(t, url, "?", "URL should contain query parameters")
}

// ─── ExchangeCode — only reachable with real HTTP; just ensure constructor works ─

func TestGoogleOAuthService_NewGoogleOAuthService_SetsFields(t *testing.T) {
	svc := NewGoogleOAuthService("cid", "csec", "https://example.com/cb")
	assert.Equal(t, "cid", svc.clientID)
	assert.Equal(t, "csec", svc.clientSecret)
	assert.Equal(t, "https://example.com/cb", svc.redirectURI)
}
