package repository

import (
	"database/sql"
	"time"

	"github.com/apsferreira/auth-service/backend/internal/domain"
	"github.com/apsferreira/auth-service/backend/internal/pkg/database"
	"github.com/google/uuid"
)

type TokenRepository struct{}

func NewTokenRepository() *TokenRepository {
	return &TokenRepository{}
}

func (r *TokenRepository) Create(token *domain.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at)
	          VALUES ($1, $2, $3, $4, $5)`
	_, err := database.DB.Exec(query,
		token.ID, token.UserID, token.TokenHash, token.ExpiresAt, token.CreatedAt,
	)
	return err
}

func (r *TokenRepository) FindByHash(tokenHash string) (*domain.RefreshToken, error) {
	token := &domain.RefreshToken{}
	var revokedAt sql.NullTime

	query := `SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
	          FROM refresh_tokens
	          WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > NOW()`
	err := database.DB.QueryRow(query, tokenHash).Scan(
		&token.ID, &token.UserID, &token.TokenHash,
		&token.ExpiresAt, &revokedAt, &token.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if revokedAt.Valid {
		token.RevokedAt = &revokedAt.Time
	}

	return token, nil
}

func (r *TokenRepository) Revoke(id uuid.UUID) error {
	query := `UPDATE refresh_tokens SET revoked_at = $1 WHERE id = $2`
	_, err := database.DB.Exec(query, time.Now(), id)
	return err
}

func (r *TokenRepository) RevokeAllForUser(userID uuid.UUID) error {
	query := `UPDATE refresh_tokens SET revoked_at = $1 WHERE user_id = $2 AND revoked_at IS NULL`
	_, err := database.DB.Exec(query, time.Now(), userID)
	return err
}

func (r *TokenRepository) DeleteExpired() error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`
	_, err := database.DB.Exec(query)
	return err
}
