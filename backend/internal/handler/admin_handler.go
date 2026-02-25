package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/apsferreira/auth-service/backend/internal/domain"
	"github.com/apsferreira/auth-service/backend/internal/service"
)

type AdminHandler struct {
	adminService *service.AdminService
	eventService *service.EventService
}

func NewAdminHandler(adminService *service.AdminService, eventService *service.EventService) *AdminHandler {
	return &AdminHandler{adminService: adminService, eventService: eventService}
}

// --- Services ---

func (h *AdminHandler) ListServices(c *fiber.Ctx) error {
	tenantID, err := uuid.Parse(c.Locals("tenantID").(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid tenant"})
	}

	services, err := h.adminService.ListServices(tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list services"})
	}

	return c.JSON(services)
}

func (h *AdminHandler) CreateService(c *fiber.Ctx) error {
	tenantID, err := uuid.Parse(c.Locals("tenantID").(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid tenant"})
	}

	var req domain.ServiceCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	svc, err := h.adminService.CreateService(tenantID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(svc)
}

func (h *AdminHandler) GetService(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid service id"})
	}

	svc, err := h.adminService.GetService(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "service not found"})
	}

	return c.JSON(svc)
}

func (h *AdminHandler) UpdateService(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid service id"})
	}

	var req domain.ServiceUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	svc, err := h.adminService.UpdateService(id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(svc)
}

func (h *AdminHandler) DeleteService(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid service id"})
	}

	if err := h.adminService.DeleteService(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// --- Permissions ---

func (h *AdminHandler) ListServicePermissions(c *fiber.Ctx) error {
	serviceID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid service id"})
	}

	perms, err := h.adminService.ListPermissionsByService(serviceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list permissions"})
	}

	return c.JSON(perms)
}

func (h *AdminHandler) CreateServicePermission(c *fiber.Ctx) error {
	serviceID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid service id"})
	}

	var req domain.PermissionCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	perm, err := h.adminService.CreatePermission(serviceID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(perm)
}

func (h *AdminHandler) ListAllPermissions(c *fiber.Ctx) error {
	perms, err := h.adminService.ListAllPermissions()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list permissions"})
	}

	return c.JSON(perms)
}

func (h *AdminHandler) DeletePermission(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid permission id"})
	}

	if err := h.adminService.DeletePermission(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// --- Roles ---

func (h *AdminHandler) ListRoles(c *fiber.Ctx) error {
	tenantID, err := uuid.Parse(c.Locals("tenantID").(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid tenant"})
	}

	roles, err := h.adminService.ListRoles(tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list roles"})
	}

	return c.JSON(roles)
}

func (h *AdminHandler) CreateRole(c *fiber.Ctx) error {
	tenantID, err := uuid.Parse(c.Locals("tenantID").(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid tenant"})
	}

	var req domain.RoleCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	role, err := h.adminService.CreateRole(tenantID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(role)
}

func (h *AdminHandler) UpdateRole(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid role id"})
	}

	var req domain.RoleUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	role, err := h.adminService.UpdateRole(id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(role)
}

func (h *AdminHandler) GetRolePermissions(c *fiber.Ctx) error {
	roleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid role id"})
	}

	ids, err := h.adminService.GetRolePermissions(roleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get role permissions"})
	}

	return c.JSON(fiber.Map{"permission_ids": ids})
}

func (h *AdminHandler) SetRolePermissions(c *fiber.Ctx) error {
	roleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid role id"})
	}

	var req domain.RolePermissionsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	permIDs := make([]uuid.UUID, 0, len(req.PermissionIDs))
	for _, idStr := range req.PermissionIDs {
		pid, err := uuid.Parse(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid permission id: " + idStr})
		}
		permIDs = append(permIDs, pid)
	}

	if err := h.adminService.SetRolePermissions(roleID, permIDs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "role permissions updated"})
}

// --- Audit Events ---

// ListEvents handles GET /api/v1/admin/events
func (h *AdminHandler) ListEvents(c *fiber.Ctx) error {
	filter := domain.AuthEventFilter{
		EventType: c.Query("event_type"),
		Email:     c.Query("email"),
	}

	if limitStr := c.Query("limit", "50"); limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil && v > 0 {
			filter.Limit = v
		}
	}
	if offsetStr := c.Query("offset", "0"); offsetStr != "" {
		if v, err := strconv.Atoi(offsetStr); err == nil && v >= 0 {
			filter.Offset = v
		}
	}
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if uid, err := uuid.Parse(userIDStr); err == nil {
			filter.UserID = &uid
		}
	}

	resp, err := h.eventService.List(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list events"})
	}

	return c.JSON(resp)
}
