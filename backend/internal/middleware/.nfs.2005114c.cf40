package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type rateLimitEntry struct {
	count   int
	resetAt time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	entries map[string]*rateLimitEntry
	max     int
	window  time.Duration
}

func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		entries: make(map[string]*rateLimitEntry),
		max:     maxRequests,
		window:  window,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) Handler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()

		rl.mu.Lock()
		entry, exists := rl.entries[ip]

		if !exists || time.Now().After(entry.resetAt) {
			rl.entries[ip] = &rateLimitEntry{count: 1, resetAt: time.Now().Add(rl.window)}
			rl.mu.Unlock()
			return c.Next()
		}

		if entry.count >= rl.max {
			retryAfter := int(time.Until(entry.resetAt).Seconds())
			rl.mu.Unlock()
			return c.Status(429).JSON(fiber.Map{
				"error":       "too many requests",
				"retry_after": retryAfter,
			})
		}

		entry.count++
		rl.mu.Unlock()
		return c.Next()
	}
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for k, v := range rl.entries {
			if now.After(v.resetAt) {
				delete(rl.entries, k)
			}
		}
		rl.mu.Unlock()
	}
}
