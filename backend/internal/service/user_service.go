package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/apsferreira/auth-service/backend/internal/domain"
	"github.com/apsferreira/auth-service/backend/internal/repository"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) List(tenantID uuid.UUID) ([]*domain.User, error) {
	return s.userRepo.List(tenantID)
}

func (s *UserService) GetByID(userID uuid.UUID) (*domain.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *UserService) Create(req domain.UserCreateRequest) (*domain.User, error) {
	if req.Email == "" || req.TenantID == "" {
		return nil, fmt.Errorf("email and tenant_id are required")
	}

	tenantID, err := uuid.Parse(req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant_id: %w", err)
	}

	// Check if user already exists
	existing, err := s.userRepo.FindByEmail(req.Email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("user with this email already exists")
	}

	now := time.Now()
	user := &domain.User{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Email:     req.Email,
		FullName:  req.FullName,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if req.RoleID != nil {
		roleID, err := uuid.Parse(*req.RoleID)
		if err == nil {
			user.RoleID = &roleID
		}
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Assign role if provided
	if user.RoleID != nil {
		_ = s.userRepo.AddUserRole(user.ID, *user.RoleID)
	}

	return s.userRepo.FindByID(user.ID)
}

func (s *UserService) Update(id uuid.UUID, req domain.UserUpdateRequest) (*domain.User, error) {
	return s.userRepo.Update(id, req)
}

func (s *UserService) Delete(id uuid.UUID) error {
	return s.userRepo.Delete(id)
}
