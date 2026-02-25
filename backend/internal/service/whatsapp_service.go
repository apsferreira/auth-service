package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// WhatsAppService delivers OTP codes via Evolution API (self-hosted WhatsApp gateway).
// Evolution API runs on the shared-infra network at shared-evolution-api:8080.
type WhatsAppService struct {
	apiURL       string
	apiKey       string
	instanceName string
	defaultPhone string // E.164 without '+', e.g. "5511999990000"
	client       *http.Client
}

func NewWhatsAppService(apiURL, apiKey, instanceName, defaultPhone string) *WhatsAppService {
	// Skip TLS verification — VPN intercepts HTTPS with self-signed cert
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
	}
	return &WhatsAppService{
		apiURL:       apiURL,
		apiKey:       apiKey,
		instanceName: instanceName,
		defaultPhone: defaultPhone,
		client:       &http.Client{Timeout: 15 * time.Second, Transport: transport},
	}
}

// IsConfigured returns true when all required fields are set.
func (w *WhatsAppService) IsConfigured() bool {
	return w.apiURL != "" && w.apiKey != "" && w.instanceName != "" && w.defaultPhone != ""
}

// SendOTP sends the OTP code to the configured WhatsApp number.
func (w *WhatsAppService) SendOTP(toEmail, code string) error {
	if !w.IsConfigured() {
		return nil
	}

	text := fmt.Sprintf(
		"🔐 *Código de acesso*\n\nUsuário: %s\nCódigo: *%s*\n\n⏱ Válido por 10 minutos.",
		toEmail, code,
	)

	payload := map[string]interface{}{
		"number":  w.defaultPhone,
		"text":    text,
		"options": map[string]interface{}{"delay": 0},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("whatsapp: marshal error: %w", err)
	}

	url := fmt.Sprintf("%s/message/sendText/%s", w.apiURL, w.instanceName)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("whatsapp: request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", w.apiKey)

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("whatsapp: send error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("whatsapp: API returned %d: %s", resp.StatusCode, string(respBody))
	}

	log.Printf("[WHATSAPP] OTP delivered to %s for %s", w.defaultPhone, toEmail)
	return nil
}

// GetInstanceQRCode retrieves the QR code URL to connect a WhatsApp number.
// Call this once to pair your WhatsApp with the Evolution API instance.
func (w *WhatsAppService) GetInstanceQRCode() (string, error) {
	url := fmt.Sprintf("%s/instance/connect/%s", w.apiURL, w.instanceName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("apikey", w.apiKey)

	resp, err := w.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("whatsapp: QR code error: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if qr, ok := result["code"].(string); ok {
		return qr, nil
	}
	if base64, ok := result["base64"].(string); ok {
		return base64, nil
	}

	return fmt.Sprintf("%v", result), nil
}
