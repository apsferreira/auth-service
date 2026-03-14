package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)

		log.Printf("[%s] %s %s - %d - %v",
			c.Method(), c.Path(), c.IP(), c.Response().StatusCode(), latency)

		return err
	}
}
