package repository

import (
	"time"

	"github.com/apsferreira/auth-service/backend/internal/domain"
	"github.com/google/uuid"
)

// UserRepositoryInterface defines the contract for user repository operations
type UserRepositoryInterface interface {
	FindByEmail(email string) (*domain.User, error)
	FindByID(id uuid.UUID) (*domain.User, error)
	FindByUsernameOrEmail(identifier string) (*domain.User, error)
	Create(user *domain.User) error
	Update(id uuid.UUID, req domain.UserUpdateRequest) (*domain.User, error)
	Delete(id uuid.UUID) error
	List(tenantID uuid.UUID) ([]*domain.User, error)
	SetPassword(userID uuid.UUID, hash string) error
	UpdateLastLogin(userID uuid.UUID) error
	UpdateProfile(id uuid.UUID, fullName string) (*domain.User, error)
	AddUserRole(userID, roleID uuid.UUID) error
	GetUserRolesAndPermissions(userID uuid.UUID) (roles []string, permissions map[string][]string, err error)
	GetDefaultTenantID() (uuid.UUID, error)
	GetDefaultRoleID(tenantID uuid.UUID) (uuid.UUID, error)
}



// TokenRepositoryInterface defines the contract for token repository operations
type TokenRepositoryInterface interface {
	Create(token *domain.RefreshToken) error
	FindByHash(tokenHash string) (*domain.RefreshToken, error)
	Revoke(id uuid.UUID) error
	RevokeAllForUser(userID uuid.UUID) error
	DeleteExpired() error
}

// EventRepositoryInterface defines the contract for auth event repository operations
type EventRepositoryInterface interface {
	Create(e *domain.AuthEvent) error
	List(f domain.AuthEventFilter) ([]*domain.AuthEvent, int64, error)
	DeleteOlderThan(days int) error
}

// RoleRepositoryInterface defines the contract for role repository operations
type RoleRepositoryInterface interface {
	Create(role *domain.Role) error
	FindByID(id uuid.UUID) (*domain.Role, error)
	List(tenantID uuid.UUID) ([]*domain.Role, error)
	Update(id uuid.UUID, req domain.RoleUpdateRequest) (*domain.Role, error)
}

// PermissionRepositoryInterface defines the contract for permission repository operations
type PermissionRepositoryInterface interface {
	Create(perm *domain.Permission) error
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (*domain.Permission, error)
	ListAll() ([]*domain.Permission, error)
	ListByService(serviceID uuid.UUID) ([]*domain.Permission, error)
	GetRolePermissionIDs(roleID uuid.UUID) ([]uuid.UUID, error)
	SetRolePermissions(roleID uuid.UUID, permissionIDs []uuid.UUID) error
}

// OAuthRepositoryInterface defines the contract for OAuth identity operations
type OAuthRepositoryInterface interface {
	FindByProvider(provider, providerID string) (*domain.OAuthIdentity, error)
	Upsert(identity *domain.OAuthIdentity) error
}

// ServiceRepositoryInterface defines the contract for service repository operations
type ServiceRepositoryInterface interface {
	Create(service *domain.Service) error
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (*domain.Service, error)
	FindBySlug(slug string) (*domain.Service, error)
	List(tenantID uuid.UUID) ([]*domain.Service, error)
	Update(id uuid.UUID, req domain.ServiceUpdateRequest) (*domain.Service, error)
}
