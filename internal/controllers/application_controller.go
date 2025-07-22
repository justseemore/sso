package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/justseemore/sso/internal/models"
	"github.com/justseemore/sso/internal/services"
	"strconv"
)

type ApplicationController struct {
	appService *services.ApplicationService
}

func NewApplicationController() *ApplicationController {
	return &ApplicationController{
		appService: services.NewApplicationService(),
	}
}

// CreateApplication 创建应用
func (c *ApplicationController) CreateApplication(ctx *fiber.Ctx) error {
	app := new(models.Application)

	if err := ctx.BodyParser(app); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无法解析请求体",
		})
	}

	if err := c.appService.CreateApplication(app); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":     "应用创建成功",
		"application": app,
	})
}

// UpdateApplication 更新应用
func (c *ApplicationController) UpdateApplication(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的应用ID",
		})
	}

	app := new(models.Application)
	if err := ctx.BodyParser(app); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无法解析请求体",
		})
	}

	app.ID = uint(id)
	if err := c.appService.UpdateApplication(app); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "应用更新成功",
		"application": app,
	})
}

// DeleteApplication 删除应用
func (c *ApplicationController) DeleteApplication(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的应用ID",
		})
	}

	if err := c.appService.DeleteApplication(uint(id)); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "应用删除成功",
	})
}

// GetApplication 获取应用信息
func (c *ApplicationController) GetApplication(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的应用ID",
		})
	}

	app, err := c.appService.GetApplicationByID(uint(id))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "应用不存在",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"application": app,
	})
}

// ListApplications 获取应用列表
func (c *ApplicationController) ListApplications(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))

	apps, total, err := c.appService.ListApplications(page, limit)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"applications": apps,
		"total":        total,
		"page":         page,
		"limit":        limit,
	})
}

// RegenerateClientSecret 重新生成客户端密钥
func (c *ApplicationController) RegenerateClientSecret(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的应用ID",
		})
	}

	clientSecret, err := c.appService.RegenerateClientSecret(uint(id))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "客户端密钥重新生成成功",
		"client_secret": clientSecret,
	})
}

// UpdateApplicationTheme 更新应用主题
func (c *ApplicationController) UpdateApplicationTheme(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的应用ID",
		})
	}

	type ThemeInput struct {
		ThemeID uint `json:"theme_id"`
	}

	input := new(ThemeInput)
	if err := ctx.BodyParser(input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无法解析请求体",
		})
	}

	if err := c.appService.UpdateApplicationTheme(uint(id), input.ThemeID); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "应用主题更新成功",
	})
}

// UpdateRedirectURIs 更新重定向URI
func (c *ApplicationController) UpdateRedirectURIs(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的应用ID",
		})
	}

	type URIsInput struct {
		RedirectURIs []string `json:"redirect_uris"`
	}

	input := new(URIsInput)
	if err := ctx.BodyParser(input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无法解析请求体",
		})
	}

	if err := c.appService.UpdateRedirectURIs(uint(id), input.RedirectURIs); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "重定向URI更新成功",
	})
}

// UpdateAllowedScopes 更新允许的作用域
func (c *ApplicationController) UpdateAllowedScopes(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的应用ID",
		})
	}

	type ScopesInput struct {
		AllowedScopes []string `json:"allowed_scopes"`
	}

	input := new(ScopesInput)
	if err := ctx.BodyParser(input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无法解析请求体",
		})
	}

	if err := c.appService.UpdateAllowedScopes(uint(id), input.AllowedScopes); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "允许的作用域更新成功",
	})
}

// UpdateSettings 更新应用设置
func (c *ApplicationController) UpdateSettings(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的应用ID",
		})
	}

	type SettingsInput struct {
		Settings map[string]interface{} `json:"settings"`
	}

	input := new(SettingsInput)
	if err := ctx.BodyParser(input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无法解析请求体",
		})
	}

	if err := c.appService.UpdateSettings(uint(id), input.Settings); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "应用设置更新成功",
	})
}