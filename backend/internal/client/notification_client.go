package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// NotificationClient handles communication with the notification service
type NotificationClient struct {
	baseURL    string
	httpClient *http.Client
}

// SendOTPRequest represents the request to send an OTP
type SendOTPRequest struct {
	Email string `json:"email"`
}

// VerifyOTPRequest represents the request to verify an OTP
type VerifyOTPRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

// OTPResponse represents the response after sending an OTP
type OTPResponse struct {
	Message   string    `json:"message"`
	ExpiresAt time.Time `json:"expires_at"`
}

// VerifyOTPResponse represents the response after verifying an OTP
type VerifyOTPResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
}

// NewNotificationClient creates a new notification client
func NewNotificationClient(baseURL string) *NotificationClient {
	return &NotificationClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendOTP sends an OTP via the notification service
func (c *NotificationClient) SendOTP(email string) (*OTPResponse, error) {
	req := SendOTPRequest{Email: email}
	
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(
		c.baseURL+"/otp/send",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send OTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("notification service returned status %d", resp.StatusCode)
	}

	var otpResp OTPResponse
	if err := json.NewDecoder(resp.Body).Decode(&otpResp); err != nil {
		return nil, fmt.Errorf("failed to decode OTP response: %w", err)
	}

	return &otpResp, nil
}

// VerifyOTP verifies an OTP via the notification service
func (c *NotificationClient) VerifyOTP(email, code string) (*VerifyOTPResponse, error) {
	req := VerifyOTPRequest{
		Email: email,
		Code:  code,
	}
	
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(
		c.baseURL+"/otp/verify",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send verify request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("notification service returned status %d", resp.StatusCode)
	}

	var verifyResp VerifyOTPResponse
	if err := json.NewDecoder(resp.Body).Decode(&verifyResp); err != nil {
		return nil, fmt.Errorf("failed to decode verify response: %w", err)
	}

	return &verifyResp, nil
}