package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/justseemore/sso/internal/models"
	"github.com/justseemore/sso/internal/services"
	"strconv"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
	}
}

// Register 注册用户
func (c *UserController) Register(ctx *fiber.Ctx) error {
	user := new(models.User)

	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无法解析请求体",
		})
	}

	if err := c.userService.Register(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "用户注册成功",
		"user":    user,
	})
}

// Login 用户登录
func (c *UserController) Login(ctx *fiber.Ctx) error {
	type LoginInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	input := new(LoginInput)

	if err := ctx.BodyParser(input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无法解析请求体",
		})
	}

	user, tokens, err := c.userService.Login(input.Username, input.Password)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "登录成功",
		"user":          user,
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"expires_at":    tokens.AtExpires,
	})
}

// GetUser 获取用户信息
func (c *UserController) GetUser(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的用户ID",
		})
	}

	user, err := c.userService.GetUserByID(uint(id))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "用户不存在",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": user,
	})
}

// UpdateUser 更新用户信息
func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的用户ID",
		})
	}

	user := new(models.User)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无法解析请求体",
		})
	}

	user.ID = uint(id)
	if err := c.userService.UpdateUser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "用户更新成功",
		"user":    user,
	})
}

// DeleteUser 删除用户
func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的用户ID",
		})
	}

	if err := c.userService.DeleteUser(uint(id)); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "用户删除成功",
	})
}

// ListUsers 获取用户列表
func (c *UserController) ListUsers(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))

	users, total, err := c.userService.ListUsers(page, limit)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// AssignRole 分配角色
func (c *UserController) AssignRole(ctx *fiber.Ctx) error {
	userID, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的用户ID",
		})
	}

	type RoleInput struct {
		RoleID uint `json:"role_id"`
	}

	input := new(RoleInput)
	if err := ctx.BodyParser(input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无法解析请求体",
		})
	}

	if err := c.userService.AssignRole(uint(userID), input.RoleID); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "角色分配成功",
	})
}

// RemoveRole 移除角色
func (c *UserController) RemoveRole(ctx *fiber.Ctx) error {
	userID, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的用户ID",
		})
	}

	type RoleInput struct {
		RoleID uint `json:"role_id"`
	}

	input := new(RoleInput)
	if err := ctx.BodyParser(input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无法解析请求体",
		})
	}

	if err := c.userService.RemoveRole(uint(userID), input.RoleID); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "角色移除成功",
	})
}

// ChangePassword 修改密码
func (c *UserController) ChangePassword(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的用户ID",
		})
	}

	type PasswordInput struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	input := new(PasswordInput)
	if err := ctx.BodyParser(input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无法解析请求体",
		})
	}

	if err := c.userService.ChangePassword(uint(id), input.OldPassword, input.NewPassword); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "密码修改成功",
	})
}