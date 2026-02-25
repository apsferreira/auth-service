package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/apsferreira/auth-service/backend/internal/domain"
	"github.com/apsferreira/auth-service/backend/internal/pkg/database"
	"github.com/google/uuid"
)

type RoleRepository struct{}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{}
}

func (r *RoleRepository) Create(role *domain.Role) error {
	query := `INSERT INTO roles (id, tenant_id, name, description, level, is_system, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := database.DB.Exec(query,
		role.ID, role.TenantID, role.Name, role.Description,
		role.Level, role.IsSystem, role.CreatedAt, role.UpdatedAt,
	)
	return err
}

func (r *RoleRepository) FindByID(id uuid.UUID) (*domain.Role, error) {
	role := &domain.Role{}
	query := `SELECT id, tenant_id, name, description, level, is_system, created_at, updated_at
	          FROM roles WHERE id = $1`

	var tenantID sql.NullString
	var description sql.NullString

	err := database.DB.QueryRow(query, id).Scan(
		&role.ID, &tenantID, &role.Name, &description,
		&role.Level, &role.IsSystem, &role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if tenantID.Valid {
		tid, _ := uuid.Parse(tenantID.String)
		role.TenantID = &tid
	}
	if description.Valid {
		role.Description = description.String
	}
	return role, nil
}

func (r *RoleRepository) List(tenantID uuid.UUID) ([]*domain.Role, error) {
	query := `SELECT id, tenant_id, name, description, level, is_system, created_at, updated_at
	          FROM roles WHERE tenant_id = $1 ORDER BY level DESC`

	rows, err := database.DB.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]*domain.Role, 0)
	for rows.Next() {
		role := &domain.Role{}
		var tid sql.NullString
		var description sql.NullString

		err := rows.Scan(
			&role.ID, &tid, &role.Name, &description,
			&role.Level, &role.IsSystem, &role.CreatedAt, &role.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if tid.Valid {
			tenantUUID, _ := uuid.Parse(tid.String)
			role.TenantID = &tenantUUID
		}
		if description.Valid {
			role.Description = description.String
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (r *RoleRepository) Update(id uuid.UUID, req domain.RoleUpdateRequest) (*domain.Role, error) {
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
	if req.Level != nil {
		updates = append(updates, fmt.Sprintf("level = $%d", argIndex))
		args = append(args, *req.Level)
		argIndex++
	}

	query := fmt.Sprintf(`UPDATE roles SET %s WHERE id = $%d
	          RETURNING id, tenant_id, name, description, level, is_system, created_at, updated_at`,
		strings.Join(updates, ", "), argIndex)
	args = append(args, id)

	role := &domain.Role{}
	var tid sql.NullString
	var description sql.NullString

	err := database.DB.QueryRow(query, args...).Scan(
		&role.ID, &tid, &role.Name, &description,
		&role.Level, &role.IsSystem, &role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if tid.Valid {
		tenantUUID, _ := uuid.Parse(tid.String)
		role.TenantID = &tenantUUID
	}
	if description.Valid {
		role.Description = description.String
	}
	return role, nil
}
