package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/apsferreira/auth-service/backend/internal/domain"
	"github.com/apsferreira/auth-service/backend/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// List handles GET /api/v1/users
func (h *UserHandler) List(c *fiber.Ctx) error {
	if h.userService == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "user service not configured"})
	}

	tenantIDStr, ok := c.Locals("tenantID").(string)
	if !ok {
		return c.Status(401).JSON(domain.ErrorResponse{Error: "unauthorized"})
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.Status(400).JSON(domain.ErrorResponse{Error: "invalid tenant ID"})
	}

	users, err := h.userService.List(tenantID)
	if err != nil {
		return c.Status(500).JSON(domain.ErrorResponse{Error: "failed to list users"})
	}

	return c.Status(200).JSON(fiber.Map{"users": users, "count": len(users)})
}

// Create handles POST /api/v1/users
func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req domain.UserCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(domain.ErrorResponse{Error: "invalid request body"})
	}

	if req.Email == "" {
		return c.Status(400).JSON(domain.ErrorResponse{Error: "email is required"})
	}

	if h.userService == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "user service not configured"})
	}

	// Default to caller's tenant if not specified
	if req.TenantID == "" {
		if tid, ok := c.Locals("tenantID").(string); ok {
			req.TenantID = tid
		}
	}

	user, err := h.userService.Create(req)
	if err != nil {
		return c.Status(400).JSON(domain.ErrorResponse{Error: err.Error()})
	}

	return c.Status(201).JSON(user)
}

// GetByID handles GET /api/v1/users/:id
func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(domain.ErrorResponse{Error: "invalid user ID"})
	}

	if h.userService == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "user service not configured"})
	}

	user, err := h.userService.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(domain.ErrorResponse{Error: "user not found"})
	}

	return c.Status(200).JSON(user)
}

// Update handles PUT /api/v1/users/:id
func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(domain.ErrorResponse{Error: "invalid user ID"})
	}

	var req domain.UserUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(domain.ErrorResponse{Error: "invalid request body"})
	}

	if h.userService == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "user service not configured"})
	}

	user, err := h.userService.Update(id, req)
	if err != nil {
		return c.Status(400).JSON(domain.ErrorResponse{Error: err.Error()})
	}

	return c.Status(200).JSON(user)
}

// Delete handles DELETE /api/v1/users/:id
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(domain.ErrorResponse{Error: "invalid user ID"})
	}

	if h.userService == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "user service not configured"})
	}

	if err := h.userService.Delete(id); err != nil {
		return c.Status(400).JSON(domain.ErrorResponse{Error: err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"message": "user deleted successfully"})
}
