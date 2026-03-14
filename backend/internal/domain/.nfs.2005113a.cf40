package domain

// OTPRequest is the body for POST /auth/request-otp
type OTPRequest struct {
	Email   string `json:"email"`
	Channel string `json:"channel"` // email | telegram | whatsapp — defaults to "email"
}

// OTPVerifyRequest is the body for POST /auth/verify-otp
type OTPVerifyRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

// RefreshRequest is the body for POST /auth/refresh
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// AuthResponse is returned after successful verify-otp or refresh
type AuthResponse struct {
	AccessToken  string              `json:"access_token"`
	RefreshToken string              `json:"refresh_token"`
	ExpiresIn    int                 `json:"expires_in"`
	User         *User               `json:"user"`
	Roles        []string            `json:"roles"`
	Permissions  map[string][]string `json:"permissions"`
}

// OTPResponse is returned after requesting an OTP
type OTPResponse struct {
	Message   string `json:"message"`
	ExpiresIn int    `json:"expires_in"`
	Channel   string `json:"channel"` // the channel used to deliver the OTP
}

// ValidateRequest represents a token validation request
type ValidateRequest struct {
	Token string `json:"token"`
}

// ValidateResponse represents a validation response
type ValidateResponse struct {
	Valid       bool                `json:"valid"`
	UserID      *string             `json:"user_id,omitempty"`
	TenantID    *string             `json:"tenant_id,omitempty"`
	Email       *string             `json:"email,omitempty"`
	Roles       []string            `json:"roles,omitempty"`
	Permissions map[string][]string `json:"permissions,omitempty"`
	Message     string              `json:"message,omitempty"`
}

// GoogleLoginURLResponse is returned by GET /auth/google to initiate the OAuth flow
type GoogleLoginURLResponse struct {
	URL string `json:"url"`
}

// ErrorResponse is a standard error envelope
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// UserCreateRequest for admin user creation
type UserCreateRequest struct {
	Email    string  `json:"email"`
	FullName string  `json:"full_name"`
	RoleID   *string `json:"role_id,omitempty"`
	TenantID string  `json:"tenant_id"`
}

// UserUpdateRequest for admin user updates
type UserUpdateRequest struct {
	FullName  *string `json:"full_name,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
	RoleID    *string `json:"role_id,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

// ServiceCreateRequest for registering a new service
type ServiceCreateRequest struct {
	Name         string   `json:"name"`
	Slug         string   `json:"slug"`
	Description  string   `json:"description,omitempty"`
	RedirectURLs []string `json:"redirect_urls,omitempty"`
}

// ServiceUpdateRequest for updating a service
type ServiceUpdateRequest struct {
	Name         *string  `json:"name,omitempty"`
	Description  *string  `json:"description,omitempty"`
	RedirectURLs []string `json:"redirect_urls,omitempty"`
	IsActive     *bool    `json:"is_active,omitempty"`
}

// PermissionCreateRequest for creating a service permission
type PermissionCreateRequest struct {
	Name        string `json:"name"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description,omitempty"`
}

// RoleCreateRequest for creating a custom role
type RoleCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Level       int    `json:"level"`
}

// RoleUpdateRequest for updating a role
type RoleUpdateRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Level       *int    `json:"level,omitempty"`
}

// RolePermissionsRequest for setting role permissions
type RolePermissionsRequest struct {
	PermissionIDs []string `json:"permission_ids"`
}
