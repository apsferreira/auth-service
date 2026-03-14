package repository

import (
	"database/sql"

	"github.com/apsferreira/auth-service/backend/internal/domain"
	"github.com/apsferreira/auth-service/backend/internal/pkg/database"
	"github.com/google/uuid"
)

type PermissionRepository struct{}

func NewPermissionRepository() *PermissionRepository {
	return &PermissionRepository{}
}

func (r *PermissionRepository) Create(perm *domain.Permission) error {
	query := `INSERT INTO permissions (id, name, resource, action, description, service_id)
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := database.DB.Exec(query,
		perm.ID, perm.Name, perm.Resource, perm.Action, perm.Description, perm.ServiceID,
	)
	return err
}

func (r *PermissionRepository) FindByID(id uuid.UUID) (*domain.Permission, error) {
	perm := &domain.Permission{}
	query := `SELECT p.id, p.name, p.resource, p.action, p.description, p.service_id, COALESCE(s.slug, '')
	          FROM permissions p
	          LEFT JOIN services s ON s.id = p.service_id
	          WHERE p.id = $1`

	var description sql.NullString
	var serviceID sql.NullString
	var serviceSlug string

	err := database.DB.QueryRow(query, id).Scan(
		&perm.ID, &perm.Name, &perm.Resource, &perm.Action,
		&description, &serviceID, &serviceSlug,
	)
	if err != nil {
		return nil, err
	}

	if description.Valid {
		perm.Description = description.String
	}
	if serviceID.Valid {
		sid, _ := uuid.Parse(serviceID.String)
		perm.ServiceID = &sid
	}
	perm.ServiceSlug = serviceSlug
	return perm, nil
}

func (r *PermissionRepository) ListByService(serviceID uuid.UUID) ([]*domain.Permission, error) {
	query := `SELECT p.id, p.name, p.resource, p.action, p.description, p.service_id, COALESCE(s.slug, '')
	          FROM permissions p
	          LEFT JOIN services s ON s.id = p.service_id
	          WHERE p.service_id = $1
	          ORDER BY p.resource, p.action`

	rows, err := database.DB.Query(query, serviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	perms := make([]*domain.Permission, 0)
	for rows.Next() {
		perm := &domain.Permission{}
		var description sql.NullString
		var svcID sql.NullString
		var serviceSlug string

		err := rows.Scan(&perm.ID, &perm.Name, &perm.Resource, &perm.Action,
			&description, &svcID, &serviceSlug)
		if err != nil {
			return nil, err
		}

		if description.Valid {
			perm.Description = description.String
		}
		if svcID.Valid {
			sid, _ := uuid.Parse(svcID.String)
			perm.ServiceID = &sid
		}
		perm.ServiceSlug = serviceSlug
		perms = append(perms, perm)
	}

	return perms, nil
}

func (r *PermissionRepository) ListAll() ([]*domain.Permission, error) {
	query := `SELECT p.id, p.name, p.resource, p.action, p.description, p.service_id, COALESCE(s.slug, '')
	          FROM permissions p
	          LEFT JOIN services s ON s.id = p.service_id
	          ORDER BY s.name, p.resource, p.action`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	perms := make([]*domain.Permission, 0)
	for rows.Next() {
		perm := &domain.Permission{}
		var description sql.NullString
		var svcID sql.NullString
		var serviceSlug string

		err := rows.Scan(&perm.ID, &perm.Name, &perm.Resource, &perm.Action,
			&description, &svcID, &serviceSlug)
		if err != nil {
			return nil, err
		}

		if description.Valid {
			perm.Description = description.String
		}
		if svcID.Valid {
			sid, _ := uuid.Parse(svcID.String)
			perm.ServiceID = &sid
		}
		perm.ServiceSlug = serviceSlug
		perms = append(perms, perm)
	}

	return perms, nil
}

func (r *PermissionRepository) Delete(id uuid.UUID) error {
	// Also remove from role_permissions
	_, _ = database.DB.Exec(`DELETE FROM role_permissions WHERE permission_id = $1`, id)
	result, err := database.DB.Exec(`DELETE FROM permissions WHERE id = $1`, id)
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

func (r *PermissionRepository) GetRolePermissionIDs(roleID uuid.UUID) ([]uuid.UUID, error) {
	query := `SELECT permission_id FROM role_permissions WHERE role_id = $1`
	rows, err := database.DB.Query(query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]uuid.UUID, 0)
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *PermissionRepository) SetRolePermissions(roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clear existing
	_, err = tx.Exec(`DELETE FROM role_permissions WHERE role_id = $1`, roleID)
	if err != nil {
		return err
	}

	// Insert new
	for _, permID := range permissionIDs {
		_, err = tx.Exec(`INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			roleID, permID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
