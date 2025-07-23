package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/justseemore/sso/internal/models"
	"github.com/justseemore/sso/internal/services"
	"strconv"
)

type RoleController struct {
	roleService *services.RoleService
}

func NewRoleController() *RoleController {
	return &RoleController{
		roleService: services.NewRoleService(),
	}
}

// CreateRole 创建角色
func (c *RoleController) CreateRole(ctx *fiber.Ctx) error {
	role := new(models.Role)

	if err := ctx.BodyParser(role); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无法解析请求体",
		})
	}

	if err := c.roleService.CreateRole(role); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "角色创建成功",
		"role":    role,
	})
}

// UpdateRole 更新角色
func (c *RoleController) UpdateRole(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的角色ID",
		})
	}

	role := new(models.Role)
	if err := ctx.BodyParser(role); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无法解析请求体",
		})
	}

	role.ID = uint(id)
	if err := c.roleService.UpdateRole(role); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "角色更新成功",
		"role":    role,
	})
}

// DeleteRole 删除角色
func (c *RoleController) DeleteRole(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的角色ID",
		})
	}

	if err := c.roleService.DeleteRole(uint(id)); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "角色删除成功",
	})
}

// GetRole 获取角色信息
func (c *RoleController) GetRole(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的角色ID",
		})
	}

	role, err := c.roleService.GetRoleByID(uint(id))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "角色不存在",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"role": role,
	})
}

// ListRoles 获取角色列表
func (c *RoleController) ListRoles(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))

	roles, total, err := c.roleService.ListRoles(page, limit)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"roles": roles,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// AssignPermission 分配权限
func (c *RoleController) AssignPermission(ctx *fiber.Ctx) error {
	roleID, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的角色ID",
		})
	}

	type PermissionInput struct {
		PermissionID uint `json:"permission_id"`
	}

	input := new(PermissionInput)
	if err := ctx.BodyParser(input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无法解析请求体",
		})
	}

	if err := c.roleService.AssignPermission(uint(roleID), input.PermissionID); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "权限分配成功",
	})
}

// RemovePermission 移除权限
func (c *RoleController) RemovePermission(ctx *fiber.Ctx) error {
	roleID, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的角色ID",
		})
	}

	permissionID, err := strconv.ParseUint(ctx.Params("permissionId"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的权限ID",
		})
	}

	if err := c.roleService.RemovePermission(uint(roleID), uint(permissionID)); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "权限移除成功",
	})
}
