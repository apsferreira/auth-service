package service

import (
	"log"
	"time"

	"github.com/apsferreira/auth-service/backend/internal/domain"
	"github.com/apsferreira/auth-service/backend/internal/repository"
	"github.com/google/uuid"
)

// EventService logs authentication events for auditing.
type EventService struct {
	repo  *repository.EventRepository
	isDev bool
}

func NewEventService(repo *repository.EventRepository, isDev bool) *EventService {
	svc := &EventService{repo: repo, isDev: isDev}
	svc.startCleanupWorker(15)
	return svc
}

// startCleanupWorker deletes audit events older than retentionDays every 24 hours.
func (s *EventService) startCleanupWorker(retentionDays int) {
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		// Run immediately on startup
		if err := s.repo.DeleteOlderThan(retentionDays); err != nil {
			log.Printf("[events] cleanup error: %v", err)
		} else {
			log.Printf("[events] cleaned up audit events older than %d days", retentionDays)
		}
		for range ticker.C {
			if err := s.repo.DeleteOlderThan(retentionDays); err != nil {
				log.Printf("[events] cleanup error: %v", err)
			} else {
				log.Printf("[events] cleaned up audit events older than %d days", retentionDays)
			}
		}
	}()
}

// Log persists an auth event asynchronously so it never blocks the main request.
func (s *EventService) Log(eventType, email, ip, userAgent string, userID *uuid.UUID, meta map[string]interface{}) {
	go func() {
		e := repository.NewAuthEvent(eventType, email, ip, userAgent, userID, meta)
		if err := s.repo.Create(e); err != nil {
			log.Printf("[events] failed to persist %s for %s: %v", eventType, email, err)
		}
	}()
}

// List returns paginated auth events for the admin panel.
func (s *EventService) List(f domain.AuthEventFilter) (*domain.AuthEventsResponse, error) {
	if f.Limit <= 0 {
		f.Limit = 50
	}
	events, total, err := s.repo.List(f)
	if err != nil {
		return nil, err
	}
	if events == nil {
		events = []*domain.AuthEvent{}
	}
	return &domain.AuthEventsResponse{
		Events: events,
		Total:  total,
		Limit:  f.Limit,
		Offset: f.Offset,
	}, nil
}
