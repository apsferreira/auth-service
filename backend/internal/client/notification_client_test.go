package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── NewNotificationClient ────────────────────────────────────────────────────

func TestNewNotificationClient_SetsBaseURL(t *testing.T) {
	c := NewNotificationClient("http://notification-service:8080")
	assert.Equal(t, "http://notification-service:8080", c.baseURL)
	assert.NotNil(t, c.httpClient)
}

func TestNewNotificationClient_HTTPClientHasTimeout(t *testing.T) {
	c := NewNotificationClient("http://example.com")
	assert.Equal(t, 30*time.Second, c.httpClient.Timeout)
}

// ─── SendOTP ─────────────────────────────────────────────────────────────────

func TestSendOTP_Success_ReturnsOTPResponse(t *testing.T) {
	expiresAt := time.Now().Add(10 * time.Minute)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/otp/send", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var req SendOTPRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "user@example.com", req.Email)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(OTPResponse{
			Message:   "OTP sent",
			ExpiresAt: expiresAt,
		})
	}))
	defer server.Close()

	c := NewNotificationClient(server.URL)
	resp, err := c.SendOTP("user@example.com")
	require.NoError(t, err)
	assert.Equal(t, "OTP sent", resp.Message)
}

func TestSendOTP_ServerReturns500_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	c := NewNotificationClient(server.URL)
	_, err := c.SendOTP("user@example.com")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "500")
}

func TestSendOTP_ServerReturns404_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	c := NewNotificationClient(server.URL)
	_, err := c.SendOTP("user@example.com")
	assert.Error(t, err)
}

func TestSendOTP_ServerReturnsInvalidJSON_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not valid json {"))
	}))
	defer server.Close()

	c := NewNotificationClient(server.URL)
	_, err := c.SendOTP("user@example.com")
	assert.Error(t, err)
}

func TestSendOTP_Unreachable_ReturnsError(t *testing.T) {
	c := NewNotificationClient("http://127.0.0.1:0")
	_, err := c.SendOTP("user@example.com")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send OTP request")
}

// ─── VerifyOTP ────────────────────────────────────────────────────────────────

func TestVerifyOTP_ValidCode_ReturnsValidTrue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/otp/verify", r.URL.Path)

		var req VerifyOTPRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "user@example.com", req.Email)
		assert.Equal(t, "123456", req.Code)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(VerifyOTPResponse{Valid: true, Message: "OTP valid"})
	}))
	defer server.Close()

	c := NewNotificationClient(server.URL)
	resp, err := c.VerifyOTP("user@example.com", "123456")
	require.NoError(t, err)
	assert.True(t, resp.Valid)
	assert.Equal(t, "OTP valid", resp.Message)
}

func TestVerifyOTP_InvalidCode_ReturnsValidFalse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(VerifyOTPResponse{Valid: false, Message: "invalid or expired OTP"})
	}))
	defer server.Close()

	c := NewNotificationClient(server.URL)
	resp, err := c.VerifyOTP("user@example.com", "000000")
	require.NoError(t, err)
	assert.False(t, resp.Valid)
	assert.Equal(t, "invalid or expired OTP", resp.Message)
}

func TestVerifyOTP_ServerReturns500_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	c := NewNotificationClient(server.URL)
	_, err := c.VerifyOTP("user@example.com", "123456")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "500")
}

func TestVerifyOTP_ServerReturnsInvalidJSON_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{bad json"))
	}))
	defer server.Close()

	c := NewNotificationClient(server.URL)
	_, err := c.VerifyOTP("user@example.com", "123456")
	assert.Error(t, err)
}

func TestVerifyOTP_Unreachable_ReturnsError(t *testing.T) {
	c := NewNotificationClient("http://127.0.0.1:0")
	_, err := c.VerifyOTP("user@example.com", "123456")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send verify request")
}

func TestVerifyOTP_SendsCorrectPayload(t *testing.T) {
	var receivedEmail, receivedCode string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req VerifyOTPRequest
		json.NewDecoder(r.Body).Decode(&req)
		receivedEmail = req.Email
		receivedCode = req.Code
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(VerifyOTPResponse{Valid: true})
	}))
	defer server.Close()

	c := NewNotificationClient(server.URL)
	c.VerifyOTP("specific@example.com", "987654")

	assert.Equal(t, "specific@example.com", receivedEmail)
	assert.Equal(t, "987654", receivedCode)
}
