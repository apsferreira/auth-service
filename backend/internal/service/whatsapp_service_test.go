package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── IsConfigured ────────────────────────────────────────────────────────────

func TestWhatsAppService_IsConfigured_AllFieldsPresent_ReturnsTrue(t *testing.T) {
	svc := NewWhatsAppService("http://localhost:8080", "api-key", "instance-1", "5511999990000")
	assert.True(t, svc.IsConfigured())
}

func TestWhatsAppService_IsConfigured_MissingAPIURL_ReturnsFalse(t *testing.T) {
	svc := NewWhatsAppService("", "api-key", "instance-1", "5511999990000")
	assert.False(t, svc.IsConfigured())
}

func TestWhatsAppService_IsConfigured_MissingAPIKey_ReturnsFalse(t *testing.T) {
	svc := NewWhatsAppService("http://localhost:8080", "", "instance-1", "5511999990000")
	assert.False(t, svc.IsConfigured())
}

func TestWhatsAppService_IsConfigured_MissingInstanceName_ReturnsFalse(t *testing.T) {
	svc := NewWhatsAppService("http://localhost:8080", "api-key", "", "5511999990000")
	assert.False(t, svc.IsConfigured())
}

func TestWhatsAppService_IsConfigured_MissingDefaultPhone_ReturnsFalse(t *testing.T) {
	svc := NewWhatsAppService("http://localhost:8080", "api-key", "instance-1", "")
	assert.False(t, svc.IsConfigured())
}

func TestWhatsAppService_IsConfigured_AllEmpty_ReturnsFalse(t *testing.T) {
	svc := NewWhatsAppService("", "", "", "")
	assert.False(t, svc.IsConfigured())
}

// ─── SendOTP — not configured ─────────────────────────────────────────────

func TestWhatsAppService_SendOTP_NotConfigured_ReturnsNilError(t *testing.T) {
	// When not configured, SendOTP should silently return nil
	svc := NewWhatsAppService("", "", "", "")
	err := svc.SendOTP("user@example.com", "123456")
	assert.NoError(t, err)
}

// ─── SendOTP — with mock HTTP server ─────────────────────────────────────────

func TestWhatsAppService_SendOTP_ServerReturnsOK_Succeeds(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Contains(t, r.URL.Path, "/message/sendText/")
		assert.Equal(t, "test-api-key", r.Header.Get("apikey"))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	svc := NewWhatsAppService(server.URL, "test-api-key", "test-instance", "5511999990000")
	err := svc.SendOTP("user@example.com", "123456")
	assert.NoError(t, err)
}

func TestWhatsAppService_SendOTP_ServerReturnsCreated_Succeeds(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	svc := NewWhatsAppService(server.URL, "api-key", "instance", "5511999990000")
	err := svc.SendOTP("user@example.com", "654321")
	assert.NoError(t, err)
}

func TestWhatsAppService_SendOTP_ServerReturns500_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "server error"}`))
	}))
	defer server.Close()

	svc := NewWhatsAppService(server.URL, "api-key", "instance", "5511999990000")
	err := svc.SendOTP("user@example.com", "999999")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "500")
}

func TestWhatsAppService_SendOTP_ServerReturns401_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	svc := NewWhatsAppService(server.URL, "bad-key", "instance", "5511999990000")
	err := svc.SendOTP("user@example.com", "000000")
	assert.Error(t, err)
}

func TestWhatsAppService_SendOTP_InvalidURL_ReturnsError(t *testing.T) {
	svc := NewWhatsAppService("http://127.0.0.1:0", "api-key", "instance", "5511999990000")
	err := svc.SendOTP("user@example.com", "123456")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "whatsapp: send error")
}

func TestWhatsAppService_SendOTP_URLContainsInstanceName(t *testing.T) {
	var capturedPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	svc := NewWhatsAppService(server.URL, "key", "my-instance", "5511999990000")
	err := svc.SendOTP("test@example.com", "111111")
	require.NoError(t, err)
	assert.Contains(t, capturedPath, "my-instance")
}

func TestWhatsAppService_SendOTP_SendsAPIKeyHeader(t *testing.T) {
	var receivedAPIKey string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAPIKey = r.Header.Get("apikey")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	svc := NewWhatsAppService(server.URL, "secret-api-key", "instance", "5511999990000")
	err := svc.SendOTP("test@example.com", "222222")
	require.NoError(t, err)
	assert.Equal(t, "secret-api-key", receivedAPIKey)
}

// ─── NewWhatsAppService ───────────────────────────────────────────────────────

func TestNewWhatsAppService_SetsFields(t *testing.T) {
	svc := NewWhatsAppService("http://api.example.com", "my-key", "my-instance", "+5511999990000")
	assert.Equal(t, "http://api.example.com", svc.apiURL)
	assert.Equal(t, "my-key", svc.apiKey)
	assert.Equal(t, "my-instance", svc.instanceName)
	assert.Equal(t, "+5511999990000", svc.defaultPhone)
	assert.NotNil(t, svc.client)
}
