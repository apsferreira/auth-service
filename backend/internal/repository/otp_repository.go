package repository

import (
	"time"

	"github.com/apsferreira/auth-service/backend/internal/domain"
	"github.com/apsferreira/auth-service/backend/internal/pkg/database"
	"github.com/google/uuid"
)

type OTPRepository struct{}

func NewOTPRepository() *OTPRepository {
	return &OTPRepository{}
}

func (r *OTPRepository) Create(otp *domain.OTPCode) error {
	query := `INSERT INTO otp_codes (id, email, code_hash, channel, attempts, expires_at, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := database.DB.Exec(query,
		otp.ID, otp.Email, otp.CodeHash, otp.Channel,
		otp.Attempts, otp.ExpiresAt, otp.CreatedAt,
	)
	return err
}

func (r *OTPRepository) FindLatestByEmail(email string) (*domain.OTPCode, error) {
	otp := &domain.OTPCode{}
	query := `SELECT id, email, code_hash, channel, attempts, expires_at, created_at
	          FROM otp_codes
	          WHERE email = $1 AND expires_at > NOW()
	          ORDER BY created_at DESC LIMIT 1`
	err := database.DB.QueryRow(query, email).Scan(
		&otp.ID, &otp.Email, &otp.CodeHash, &otp.Channel,
		&otp.Attempts, &otp.ExpiresAt, &otp.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return otp, nil
}

func (r *OTPRepository) IncrementAttempts(id uuid.UUID) error {
	query := `UPDATE otp_codes SET attempts = attempts + 1 WHERE id = $1`
	_, err := database.DB.Exec(query, id)
	return err
}

func (r *OTPRepository) DeleteByEmail(email string) error {
	query := `DELETE FROM otp_codes WHERE email = $1`
	_, err := database.DB.Exec(query, email)
	return err
}

func (r *OTPRepository) DeleteExpired() error {
	query := `DELETE FROM otp_codes WHERE expires_at < NOW()`
	_, err := database.DB.Exec(query)
	return err
}

func (r *OTPRepository) CountRecentByEmail(email string, since time.Time) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM otp_codes WHERE email = $1 AND created_at > $2`
	err := database.DB.QueryRow(query, email, since).Scan(&count)
	return count, err
}
