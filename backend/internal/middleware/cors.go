package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func CORS(allowedOrigins string) fiber.Handler {
	origins := strings.Split(allowedOrigins, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	return func(c *fiber.Ctx) error {
		origin := c.Get("Origin")
		for _, allowed := range origins {
			if allowed == "*" || allowed == origin {
				c.Set("Access-Control-Allow-Origin", origin)
				break
			}
		}
		c.Set("Access-Control-Allow-Credentials", "true")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(204)
		}

		return c.Next()
	}
}
