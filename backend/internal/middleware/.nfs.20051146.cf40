package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	jwtpkg "github.com/apsferreira/auth-service/backend/internal/pkg/jwt"
)

func Auth(jwtService *jwtpkg.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "missing authorization header"})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{"error": "invalid authorization header format"})
		}

		claims, err := jwtService.ValidateAccessToken(parts[1])
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid or expired token"})
		}

		c.Locals("userID", claims.Subject)
		c.Locals("tenantID", claims.TenantID)
		c.Locals("email", claims.Email)
		c.Locals("roles", claims.Roles)
		c.Locals("permissions", claims.Permissions)

		return c.Next()
	}
}

func RequireRole(requiredRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoles, ok := c.Locals("roles").([]string)
		if !ok {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
		for _, ur := range userRoles {
			for _, required := range requiredRoles {
				if ur == required {
					return c.Next()
				}
			}
		}
		return c.Status(403).JSON(fiber.Map{"error": "insufficient role"})
	}
}

// RequirePermission checks if the user has a specific permission in any service.
// Permission format: "resource.action" (e.g., "books.create", "users.manage")
func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		permissions, ok := c.Locals("permissions").(map[string][]string)
		if !ok {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
		for _, perms := range permissions {
			for _, p := range perms {
				if p == permission {
					return c.Next()
				}
			}
		}
		return c.Status(403).JSON(fiber.Map{"error": "insufficient permissions"})
	}
}

// RequireServicePermission checks if the user has a specific permission for a specific service.
// Usage: RequireServicePermission("my-library", "books.create")
func RequireServicePermission(serviceSlug, permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		permissions, ok := c.Locals("permissions").(map[string][]string)
		if !ok {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
		servicePerms, exists := permissions[serviceSlug]
		if !exists {
			return c.Status(403).JSON(fiber.Map{"error": "no permissions for service"})
		}
		for _, p := range servicePerms {
			if p == permission {
				return c.Next()
			}
		}
		return c.Status(403).JSON(fiber.Map{"error": "insufficient permissions"})
	}
}
