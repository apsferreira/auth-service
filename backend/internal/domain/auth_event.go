package domain

import (
	"time"

	"github.com/google/uuid"
)

// Event type constants
const (
	EventOTPRequested   = "otp_requested"
	EventLoginSuccess   = "login_success"
	EventLoginFailed    = "login_failed"
	EventLogout         = "logout"
	EventTokenRefreshed = "token_refreshed"
)

// AuthEvent represents a single authentication audit entry.
type AuthEvent struct {
	ID        uuid.UUID              `json:"id"`
	EventType string                 `json:"event_type"`
	UserID    *uuid.UUID             `json:"user_id,omitempty"`
	Email     string                 `json:"email,omitempty"`
	IPAddress string                 `json:"ip_address,omitempty"`
	UserAgent string                 `json:"user_agent,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time              `json:"created_at"`

	// Loaded on demand
	UserEmail    string `json:"user_email,omitempty"`
	UserFullName string `json:"user_full_name,omitempty"`
}

// AuthEventFilter for listing events
type AuthEventFilter struct {
	EventType string
	Email     string
	UserID    *uuid.UUID
	Limit     int
	Offset    int
}

// AuthEventsResponse wraps the paginated list
type AuthEventsResponse struct {
	Events []*AuthEvent `json:"events"`
	Total  int64        `json:"total"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
}
