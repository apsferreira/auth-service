package repository

import (
	"database/sql"
	"time"

	"github.com/apsferreira/auth-service/backend/internal/domain"
	"github.com/apsferreira/auth-service/backend/internal/pkg/database"
	"github.com/google/uuid"
)

type OAuthRepository struct{}

func NewOAuthRepository() *OAuthRepository {
	return &OAuthRepository{}
}

func (r *OAuthRepository) FindByProvider(provider, providerID string) (*domain.OAuthIdentity, error) {
	db := database.DB
	identity := &domain.OAuthIdentity{}
	var avatarURL sql.NullString

	err := db.QueryRow(`
		SELECT id, user_id, provider, provider_id, email, avatar_url, created_at, updated_at
		FROM oauth_identities
		WHERE provider = $1 AND provider_id = $2
	`, provider, providerID).Scan(
		&identity.ID,
		&identity.UserID,
		&identity.Provider,
		&identity.ProviderID,
		&identity.Email,
		&avatarURL,
		&identity.CreatedAt,
		&identity.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if avatarURL.Valid {
		identity.AvatarURL = &avatarURL.String
	}
	return identity, nil
}

func (r *OAuthRepository) Upsert(identity *domain.OAuthIdentity) error {
	db := database.DB

	if identity.ID == uuid.Nil {
		identity.ID = uuid.New()
	}
	now := time.Now()
	identity.UpdatedAt = now

	_, err := db.Exec(`
		INSERT INTO oauth_identities (id, user_id, provider, provider_id, email, avatar_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (provider, provider_id)
		DO UPDATE SET
			email      = EXCLUDED.email,
			avatar_url = EXCLUDED.avatar_url,
			updated_at = EXCLUDED.updated_at
	`,
		identity.ID,
		identity.UserID,
		identity.Provider,
		identity.ProviderID,
		identity.Email,
		identity.AvatarURL,
		now,
		now,
	)
	return err
}
