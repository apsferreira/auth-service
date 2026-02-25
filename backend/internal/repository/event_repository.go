package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/apsferreira/auth-service/backend/internal/domain"
	"github.com/apsferreira/auth-service/backend/internal/pkg/database"
	"github.com/google/uuid"
)

type EventRepository struct{}

func NewEventRepository() *EventRepository {
	return &EventRepository{}
}

func (r *EventRepository) Create(e *domain.AuthEvent) error {
	meta, err := json.Marshal(e.Metadata)
	if err != nil {
		meta = []byte("{}")
	}
	query := `INSERT INTO auth_events (id, event_type, user_id, email, ip_address, user_agent, metadata, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = database.DB.Exec(query,
		e.ID, e.EventType, e.UserID, e.Email, e.IPAddress, e.UserAgent, meta, e.CreatedAt,
	)
	return err
}

func (r *EventRepository) List(f domain.AuthEventFilter) ([]*domain.AuthEvent, int64, error) {
	where := "WHERE 1=1"
	args := []interface{}{}
	idx := 1

	if f.EventType != "" {
		where += fmt.Sprintf(" AND event_type = $%d", idx)
		args = append(args, f.EventType)
		idx++
	}
	if f.Email != "" {
		where += fmt.Sprintf(" AND email ILIKE $%d", idx)
		args = append(args, "%"+f.Email+"%")
		idx++
	}
	if f.UserID != nil {
		where += fmt.Sprintf(" AND user_id = $%d", idx)
		args = append(args, *f.UserID)
		idx++
	}

	// Count
	var total int64
	countQuery := "SELECT COUNT(*) FROM auth_events " + where
	if err := database.DB.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// List with JOIN on users for display name
	if f.Limit <= 0 {
		f.Limit = 50
	}
	listQuery := `
		SELECT ae.id, ae.event_type, ae.user_id, ae.email, ae.ip_address, ae.user_agent, ae.metadata, ae.created_at,
		       COALESCE(u.email, ''), COALESCE(u.full_name, '')
		FROM auth_events ae
		LEFT JOIN users u ON u.id = ae.user_id
		` + where + fmt.Sprintf(`
		ORDER BY ae.created_at DESC
		LIMIT $%d OFFSET $%d`, idx, idx+1)
	args = append(args, f.Limit, f.Offset)

	rows, err := database.DB.Query(listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var events []*domain.AuthEvent
	for rows.Next() {
		e := &domain.AuthEvent{}
		var metaRaw []byte
		if err := rows.Scan(
			&e.ID, &e.EventType, &e.UserID, &e.Email, &e.IPAddress, &e.UserAgent, &metaRaw, &e.CreatedAt,
			&e.UserEmail, &e.UserFullName,
		); err != nil {
			return nil, 0, err
		}
		if len(metaRaw) > 0 {
			_ = json.Unmarshal(metaRaw, &e.Metadata)
		}
		events = append(events, e)
	}
	return events, total, rows.Err()
}

func (r *EventRepository) DeleteOlderThan(days int) error {
	_, err := database.DB.Exec(
		`DELETE FROM auth_events WHERE created_at < NOW() - INTERVAL '1 day' * $1`, days,
	)
	return err
}

// Helper to build a new AuthEvent ready to persist
func NewAuthEvent(eventType, email, ip, userAgent string, userID *uuid.UUID, meta map[string]interface{}) *domain.AuthEvent {
	if meta == nil {
		meta = map[string]interface{}{}
	}
	return &domain.AuthEvent{
		ID:        uuid.New(),
		EventType: eventType,
		UserID:    userID,
		Email:     email,
		IPAddress: ip,
		UserAgent: userAgent,
		Metadata:  meta,
		CreatedAt: time.Now(),
	}
}
