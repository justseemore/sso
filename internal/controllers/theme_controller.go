package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/justseemore/sso/internal/models"
	"github.com/justseemore/sso/internal/services"
	"strconv"
)

type ThemeController struct {
	themeService *services.ThemeService
}

func NewThemeController() *ThemeController {
	return &ThemeController{
		themeService: services.NewThemeService(),
	}
}

// CreateTheme 创建主题
func (c *ThemeController) CreateTheme(ctx *fiber.Ctx) error {
	theme := new(models.Theme)

	if err := ctx.BodyParser(theme); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无法解析请求体",
		})
	}

	if err := c.themeService.CreateTheme(theme); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "主题创建成功",
		"theme":   theme,
	})
}

// UpdateTheme 更新主题
func (c *ThemeController) UpdateTheme(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的主题ID",
		})
	}

	theme := new(models.Theme)
	if err := ctx.BodyParser(theme); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无法解析请求体",
		})
	}

	theme.ID = uint(id)
	if err := c.themeService.UpdateTheme(theme); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "主题更新成功",
		"theme":   theme,
	})
}

// DeleteTheme 删除主题
func (c *ThemeController) DeleteTheme(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的主题ID",
		})
	}

	if err := c.themeService.DeleteTheme(uint(id)); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "主题删除成功",
	})
}

// GetTheme 获取主题信息
func (c *ThemeController) GetTheme(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的主题ID",
		})
	}

	theme, err := c.themeService.GetThemeByID(uint(id))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "主题不存在",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"theme": theme,
	})
}

// ListThemes 获取主题列表
func (c *ThemeController) ListThemes(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))

	themes, total, err := c.themeService.ListThemes(page, limit)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"themes": themes,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}