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

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, tenant_id, email, full_name, avatar_url, is_active, role_id,
	          created_at, updated_at, last_login_at
	          FROM users WHERE email = $1 AND is_active = true`

	var avatarURL sql.NullString
	var roleID sql.NullString
	var lastLoginAt sql.NullTime

	err := database.DB.QueryRow(query, email).Scan(
		&user.ID, &user.TenantID, &user.Email, &user.FullName,
		&avatarURL, &user.IsActive, &roleID,
		&user.CreatedAt, &user.UpdatedAt, &lastLoginAt,
	)
	if err != nil {
		return nil, err
	}

	if avatarURL.Valid {
		user.AvatarURL = &avatarURL.String
	}
	if roleID.Valid {
		id, _ := uuid.Parse(roleID.String)
		user.RoleID = &id
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return user, nil
}

// FindByUsernameOrEmail fetches a user by username OR email, including password_hash.
// Used exclusively for admin panel password-based authentication.
func (r *UserRepository) FindByUsernameOrEmail(identifier string) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, tenant_id, email, full_name, avatar_url, is_active, role_id,
	          username, password_hash, created_at, updated_at, last_login_at
	          FROM users WHERE (email = $1 OR username = $1) AND is_active = true LIMIT 1`

	var avatarURL sql.NullString
	var roleID sql.NullString
	var username sql.NullString
	var passwordHash sql.NullString
	var lastLoginAt sql.NullTime

	err := database.DB.QueryRow(query, identifier).Scan(
		&user.ID, &user.TenantID, &user.Email, &user.FullName,
		&avatarURL, &user.IsActive, &roleID,
		&username, &passwordHash,
		&user.CreatedAt, &user.UpdatedAt, &lastLoginAt,
	)
	if err != nil {
		return nil, err
	}

	if avatarURL.Valid {
		user.AvatarURL = &avatarURL.String
	}
	if roleID.Valid {
		id, _ := uuid.Parse(roleID.String)
		user.RoleID = &id
	}
	if username.Valid {
		user.Username = &username.String
	}
	if passwordHash.Valid {
		user.PasswordHash = &passwordHash.String
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return user, nil
}

// SetPassword updates the password_hash for a user.
func (r *UserRepository) SetPassword(userID uuid.UUID, hash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3`
	_, err := database.DB.Exec(query, hash, time.Now(), userID)
	return err
}

func (r *UserRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, tenant_id, email, full_name, avatar_url, is_active, role_id,
	          created_at, updated_at, last_login_at
	          FROM users WHERE id = $1`

	var avatarURL sql.NullString
	var roleID sql.NullString
	var lastLoginAt sql.NullTime

	err := database.DB.QueryRow(query, id).Scan(
		&user.ID, &user.TenantID, &user.Email, &user.FullName,
		&avatarURL, &user.IsActive, &roleID,
		&user.CreatedAt, &user.UpdatedAt, &lastLoginAt,
	)
	if err != nil {
		return nil, err
	}

	if avatarURL.Valid {
		user.AvatarURL = &avatarURL.String
	}
	if roleID.Valid {
		id, _ := uuid.Parse(roleID.String)
		user.RoleID = &id
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return user, nil
}

func (r *UserRepository) Create(user *domain.User) error {
	query := `INSERT INTO users (id, tenant_id, email, full_name, is_active, role_id, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := database.DB.Exec(query,
		user.ID, user.TenantID, user.Email, user.FullName,
		user.IsActive, user.RoleID, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *UserRepository) UpdateLastLogin(userID uuid.UUID) error {
	query := `UPDATE users SET last_login_at = $1, updated_at = $1 WHERE id = $2`
	_, err := database.DB.Exec(query, time.Now(), userID)
	return err
}

func (r *UserRepository) UpdateProfile(id uuid.UUID, fullName string) (*domain.User, error) {
	query := `UPDATE users SET full_name = $1, updated_at = $2 WHERE id = $3`
	_, err := database.DB.Exec(query, fullName, time.Now(), id)
	if err != nil {
		return nil, err
	}
	return r.FindByID(id)
}

func (r *UserRepository) List(tenantID uuid.UUID) ([]*domain.User, error) {
	query := `SELECT id, tenant_id, email, full_name, avatar_url, is_active, role_id,
	          created_at, updated_at, last_login_at
	          FROM users WHERE tenant_id = $1 ORDER BY created_at DESC`

	rows, err := database.DB.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*domain.User, 0)
	for rows.Next() {
		user := &domain.User{}
		var avatarURL sql.NullString
		var roleID sql.NullString
		var lastLoginAt sql.NullTime

		err := rows.Scan(
			&user.ID, &user.TenantID, &user.Email, &user.FullName,
			&avatarURL, &user.IsActive, &roleID,
			&user.CreatedAt, &user.UpdatedAt, &lastLoginAt,
		)
		if err != nil {
			return nil, err
		}

		if avatarURL.Valid {
			user.AvatarURL = &avatarURL.String
		}
		if roleID.Valid {
			id, _ := uuid.Parse(roleID.String)
			user.RoleID = &id
		}
		if lastLoginAt.Valid {
			user.LastLoginAt = &lastLoginAt.Time
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) Update(id uuid.UUID, req domain.UserUpdateRequest) (*domain.User, error) {
	updates := []string{"updated_at = $1"}
	args := []interface{}{time.Now()}
	argIndex := 2

	if req.FullName != nil {
		updates = append(updates, fmt.Sprintf("full_name = $%d", argIndex))
		args = append(args, *req.FullName)
		argIndex++
	}
	if req.IsActive != nil {
		updates = append(updates, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *req.IsActive)
		argIndex++
	}
	if req.RoleID != nil {
		updates = append(updates, fmt.Sprintf("role_id = $%d", argIndex))
		args = append(args, *req.RoleID)
		argIndex++
	}
	if req.AvatarURL != nil {
		updates = append(updates, fmt.Sprintf("avatar_url = $%d", argIndex))
		args = append(args, *req.AvatarURL)
		argIndex++
	}

	query := fmt.Sprintf(`UPDATE users SET %s WHERE id = $%d
	          RETURNING id, tenant_id, email, full_name, avatar_url, is_active, role_id, created_at, updated_at, last_login_at`,
		strings.Join(updates, ", "), argIndex)
	args = append(args, id)

	user := &domain.User{}
	var avatarURL sql.NullString
	var roleID sql.NullString
	var lastLoginAt sql.NullTime

	err := database.DB.QueryRow(query, args...).Scan(
		&user.ID, &user.TenantID, &user.Email, &user.FullName,
		&avatarURL, &user.IsActive, &roleID,
		&user.CreatedAt, &user.UpdatedAt, &lastLoginAt,
	)
	if err != nil {
		return nil, err
	}

	if avatarURL.Valid {
		user.AvatarURL = &avatarURL.String
	}
	if roleID.Valid {
		rid, _ := uuid.Parse(roleID.String)
		user.RoleID = &rid
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return user, nil
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	query := `UPDATE users SET is_active = false, updated_at = $1 WHERE id = $2`
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

func (r *UserRepository) GetUserRolesAndPermissions(userID uuid.UUID) (roles []string, permissions map[string][]string, err error) {
	rolesQuery := `SELECT r.name FROM roles r
	               INNER JOIN user_roles ur ON ur.role_id = r.id
	               WHERE ur.user_id = $1`
	rows, err := database.DB.Query(rolesQuery, userID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	roles = make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, nil, err
		}
		roles = append(roles, name)
	}

	permQuery := `SELECT DISTINCT COALESCE(s.slug, 'global'), p.name
	              FROM permissions p
	              INNER JOIN role_permissions rp ON rp.permission_id = p.id
	              INNER JOIN user_roles ur ON ur.role_id = rp.role_id
	              LEFT JOIN services s ON s.id = p.service_id
	              WHERE ur.user_id = $1
	              ORDER BY 1, 2`
	pRows, err := database.DB.Query(permQuery, userID)
	if err != nil {
		return nil, nil, err
	}
	defer pRows.Close()

	permissions = make(map[string][]string)
	for pRows.Next() {
		var serviceSlug, permName string
		if err := pRows.Scan(&serviceSlug, &permName); err != nil {
			return nil, nil, err
		}
		permissions[serviceSlug] = append(permissions[serviceSlug], permName)
	}

	return roles, permissions, nil
}

func (r *UserRepository) AddUserRole(userID, roleID uuid.UUID) error {
	query := `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := database.DB.Exec(query, userID, roleID)
	return err
}

// GetDefaultTenantID returns the first tenant (used for auto-registration)
func (r *UserRepository) GetDefaultTenantID() (uuid.UUID, error) {
	var tenantID uuid.UUID
	query := `SELECT id FROM tenants ORDER BY created_at ASC LIMIT 1`
	err := database.DB.QueryRow(query).Scan(&tenantID)
	return tenantID, err
}

// GetDefaultRoleID returns the "user" role for the given tenant
func (r *UserRepository) GetDefaultRoleID(tenantID uuid.UUID) (uuid.UUID, error) {
	var roleID uuid.UUID
	query := `SELECT id FROM roles WHERE tenant_id = $1 AND name = 'user' LIMIT 1`
	err := database.DB.QueryRow(query, tenantID).Scan(&roleID)
	return roleID, err
}
