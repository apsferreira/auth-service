package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user entity
type User struct {
	ID           uuid.UUID  `json:"id"`
	TenantID     uuid.UUID  `json:"tenant_id"`
	Email        string     `json:"email"`
	FullName     string     `json:"full_name"`
	AvatarURL    *string    `json:"avatar_url,omitempty"`
	IsActive     bool       `json:"is_active"`
	RoleID       *uuid.UUID `json:"role_id,omitempty"`
	Username     *string    `json:"username,omitempty"`
	PasswordHash *string    `json:"-"` // never serialized
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`

	// Relations (optional, loaded on demand)
	Tenant *Tenant `json:"tenant,omitempty"`
	Roles  []*Role `json:"roles,omitempty"`
}

// AdminLoginRequest for username/password authentication (admin panel only)
type AdminLoginRequest struct {
	Identifier string `json:"identifier"` // username or email
	Password   string `json:"password"`
}

// Tenant represents a tenant (organization)
type Tenant struct {
	ID        uuid.UUID              `json:"id"`
	Name      string                 `json:"name"`
	Slug      string                 `json:"slug"`
	Plan      string                 `json:"plan"`
	Settings  map[string]interface{} `json:"settings"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// Role represents a role entity
type Role struct {
	ID          uuid.UUID  `json:"id"`
	TenantID    *uuid.UUID `json:"tenant_id,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Level       int        `json:"level"`
	IsSystem    bool       `json:"is_system"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relations
	Permissions []*Permission `json:"permissions,omitempty"`
}

// Permission represents a permission entity
type Permission struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Resource    string     `json:"resource"`
	Action      string     `json:"action"`
	Description string     `json:"description,omitempty"`
	ServiceID   *uuid.UUID `json:"service_id,omitempty"`
	ServiceSlug string     `json:"service_slug,omitempty"`
}

// Service represents a registered application/service
type Service struct {
	ID           uuid.UUID `json:"id"`
	TenantID     uuid.UUID `json:"tenant_id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Description  string    `json:"description,omitempty"`
	RedirectURLs []string  `json:"redirect_urls,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}



// OAuthIdentity links a user to a social login provider account
type OAuthIdentity struct {
	ID         uuid.UUID  `json:"id"`
	UserID     uuid.UUID  `json:"user_id"`
	Provider   string     `json:"provider"`    // "google", "github", etc.
	ProviderID string     `json:"provider_id"` // provider's unique user ID
	Email      string     `json:"email"`
	AvatarURL  *string    `json:"avatar_url,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// RefreshToken represents a stored refresh token
type RefreshToken struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	TokenHash string     `json:"-"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}
