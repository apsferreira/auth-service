package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/apsferreira/auth-service/backend/internal/domain"
	"github.com/apsferreira/auth-service/backend/internal/pkg/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ServiceRepository struct{}

func NewServiceRepository() *ServiceRepository {
	return &ServiceRepository{}
}

func (r *ServiceRepository) Create(service *domain.Service) error {
	query := `INSERT INTO services (id, tenant_id, name, slug, description, redirect_urls, is_active, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := database.DB.Exec(query,
		service.ID, service.TenantID, service.Name, service.Slug,
		service.Description, pq.Array(service.RedirectURLs),
		service.IsActive, service.CreatedAt, service.UpdatedAt,
	)
	return err
}

func (r *ServiceRepository) FindByID(id uuid.UUID) (*domain.Service, error) {
	service := &domain.Service{}
	query := `SELECT id, tenant_id, name, slug, description, redirect_urls, is_active, created_at, updated_at
	          FROM services WHERE id = $1`

	var description sql.NullString
	var redirectURLs []string

	err := database.DB.QueryRow(query, id).Scan(
		&service.ID, &service.TenantID, &service.Name, &service.Slug,
		&description, pq.Array(&redirectURLs),
		&service.IsActive, &service.CreatedAt, &service.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if description.Valid {
		service.Description = description.String
	}
	service.RedirectURLs = redirectURLs
	return service, nil
}

func (r *ServiceRepository) FindBySlug(slug string) (*domain.Service, error) {
	service := &domain.Service{}
	query := `SELECT id, tenant_id, name, slug, description, redirect_urls, is_active, created_at, updated_at
	          FROM services WHERE slug = $1`

	var description sql.NullString
	var redirectURLs []string

	err := database.DB.QueryRow(query, slug).Scan(
		&service.ID, &service.TenantID, &service.Name, &service.Slug,
		&description, pq.Array(&redirectURLs),
		&service.IsActive, &service.CreatedAt, &service.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if description.Valid {
		service.Description = description.String
	}
	service.RedirectURLs = redirectURLs
	return service, nil
}

func (r *ServiceRepository) List(tenantID uuid.UUID) ([]*domain.Service, error) {
	query := `SELECT id, tenant_id, name, slug, description, redirect_urls, is_active, created_at, updated_at
	          FROM services WHERE tenant_id = $1 ORDER BY name ASC`

	rows, err := database.DB.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	services := make([]*domain.Service, 0)
	for rows.Next() {
		service := &domain.Service{}
		var description sql.NullString
		var redirectURLs []string

		err := rows.Scan(
			&service.ID, &service.TenantID, &service.Name, &service.Slug,
			&description, pq.Array(&redirectURLs),
			&service.IsActive, &service.CreatedAt, &service.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if description.Valid {
			service.Description = description.String
		}
		service.RedirectURLs = redirectURLs
		services = append(services, service)
	}

	return services, nil
}

func (r *ServiceRepository) Update(id uuid.UUID, req domain.ServiceUpdateRequest) (*domain.Service, error) {
	updates := []string{"updated_at = $1"}
	args := []interface{}{time.Now()}
	argIndex := 2

	if req.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *req.Name)
		argIndex++
	}
	if req.Description != nil {
		updates = append(updates, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *req.Description)
		argIndex++
	}
	if req.RedirectURLs != nil {
		updates = append(updates, fmt.Sprintf("redirect_urls = $%d", argIndex))
		args = append(args, pq.Array(req.RedirectURLs))
		argIndex++
	}
	if req.IsActive != nil {
		updates = append(updates, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *req.IsActive)
		argIndex++
	}

	query := fmt.Sprintf(`UPDATE services SET %s WHERE id = $%d
	          RETURNING id, tenant_id, name, slug, description, redirect_urls, is_active, created_at, updated_at`,
		strings.Join(updates, ", "), argIndex)
	args = append(args, id)

	service := &domain.Service{}
	var description sql.NullString
	var redirectURLs []string

	err := database.DB.QueryRow(query, args...).Scan(
		&service.ID, &service.TenantID, &service.Name, &service.Slug,
		&description, pq.Array(&redirectURLs),
		&service.IsActive, &service.CreatedAt, &service.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if description.Valid {
		service.Description = description.String
	}
	service.RedirectURLs = redirectURLs
	return service, nil
}

func (r *ServiceRepository) Delete(id uuid.UUID) error {
	query := `UPDATE services SET is_active = false, updated_at = $1 WHERE id = $2`
	result, err := database.DB.Exec(query, time.Now(), id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
