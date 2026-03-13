package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// GoogleUserInfo holds the profile data returned by Google's userinfo endpoint
type GoogleUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type GoogleOAuthService struct {
	clientID     string
	clientSecret string
	redirectURI  string
}

func NewGoogleOAuthService(clientID, clientSecret, redirectURI string) *GoogleOAuthService {
	return &GoogleOAuthService{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
	}
}

func (s *GoogleOAuthService) IsConfigured() bool {
	return s.clientID != "" && s.clientSecret != "" && s.redirectURI != ""
}

// GetAuthURL builds the Google OAuth2 authorization URL.
// state is passed through unchanged and should be validated on callback.
func (s *GoogleOAuthService) GetAuthURL(state string) string {
	params := url.Values{}
	params.Set("client_id", s.clientID)
	params.Set("redirect_uri", s.redirectURI)
	params.Set("response_type", "code")
	params.Set("scope", "openid email profile")
	params.Set("state", state)
	params.Set("access_type", "online")
	params.Set("prompt", "select_account")
	return "https://accounts.google.com/o/oauth2/v2/auth?" + params.Encode()
}

// ExchangeCode exchanges an authorization code for a Google access token,
// then fetches the user's profile from Google's userinfo endpoint.
func (s *GoogleOAuthService) ExchangeCode(code string) (*GoogleUserInfo, error) {
	// 1. Exchange code for access token
	body := url.Values{}
	body.Set("code", code)
	body.Set("client_id", s.clientID)
	body.Set("client_secret", s.clientSecret)
	body.Set("redirect_uri", s.redirectURI)
	body.Set("grant_type", "authorization_code")

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", body)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed (%d): %s", resp.StatusCode, raw)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}
	if err := json.Unmarshal(raw, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}
	if tokenResp.Error != "" {
		return nil, fmt.Errorf("token error: %s", tokenResp.Error)
	}

	// 2. Fetch user profile
	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)

	client := &http.Client{}
	userResp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}
	defer userResp.Body.Close()

	userRaw, _ := io.ReadAll(userResp.Body)
	if userResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo request failed (%d): %s", userResp.StatusCode, userRaw)
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(userRaw, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	userInfo.Email = strings.ToLower(strings.TrimSpace(userInfo.Email))
	if userInfo.Email == "" {
		return nil, fmt.Errorf("Google account has no email address")
	}

	return &userInfo, nil
}
