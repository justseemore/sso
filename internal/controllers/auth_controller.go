package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/justseemore/sso/internal/services"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		authService: services.NewAuthService(),
	}
}

// Authorize 授权端点
func (c *AuthController) Authorize(ctx *fiber.Ctx) error {
	// 获取请求参数
	clientID := ctx.Query("client_id")
	redirectURI := ctx.Query("redirect_uri")
	responseType := ctx.Query("response_type")
	scope := ctx.Query("scope")
	state := ctx.Query("state")

	// 验证客户端和重定向URI
	app, err := c.authService.ValidateClientCredentials(clientID, redirectURI)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":             "invalid_request",
			"error_description": err.Error(),
		})
	}

	// 如果是授权码模式
	if responseType == "code" {
		// 如果用户已登录，则直接授权
		userID := ctx.Locals("userID")
		if userID != nil {
			// 生成授权码
			code, err := c.authService.AuthorizeUser(userID.(uint), clientID, scope)
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":             "server_error",
					"error_description": err.Error(),
				})
			}

			// 重定向到客户端
			redirectURL := redirectURI + "?code=" + code
			if state != "" {
				redirectURL += "&state=" + state
			}

			return ctx.Redirect(redirectURL)
		}

		// 如果用户未登录，则渲染登录页面
		return ctx.Render("login", fiber.Map{
			"clientID":     clientID,
			"redirectURI":  redirectURI,
			"responseType": responseType,
			"scope":        scope,
			"state":        state,
			"app":          app,
		})
	}

	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error":             "unsupported_response_type",
		"error_description": "响应类型不支持",
	})
}

// Token 令牌端点
func (c *AuthController) Token(ctx *fiber.Ctx) error {
	// 获取请求参数
	grantType := ctx.FormValue("grant_type")
	clientID := ctx.FormValue("client_id")
	clientSecret := ctx.FormValue("client_secret")

	// 验证客户端凭证
	app, err := c.authService.ValidateClientCredentials(clientID, "")
	if err != nil || app.ClientSecret != clientSecret {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":             "invalid_client",
			"error_description": "客户端凭证无效",
		})
	}

	// 根据授权类型处理
	switch grantType {
	case "authorization_code":
		// 授权码模式
		code := ctx.FormValue("code")
		redirectURI := ctx.FormValue("redirect_uri")

		// 使用授权码交换令牌
		tokens, err := c.authService.ExchangeCodeForTokens(code, clientID, redirectURI)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":             "invalid_grant",
				"error_description": err.Error(),
			})
		}

		return ctx.JSON(tokens)

	case "refresh_token":
		// 刷新令牌模式
		refreshToken := ctx.FormValue("refresh_token")

		// 使用刷新令牌获取新的访问令牌
		tokens, err := c.authService.RefreshTokens(refreshToken, clientID)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":             "invalid_grant",
				"error_description": err.Error(),
			})
		}

		return ctx.JSON(tokens)

	default:
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":             "unsupported_grant_type",
			"error_description": "授权类型不支持",
		})
	}
}

// Userinfo 用户信息端点
func (c *AuthController) Userinfo(ctx *fiber.Ctx) error {
	// 从请求头中获取访问令牌
	authHeader := ctx.Get("Authorization")
	var accessToken string
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		accessToken = authHeader[7:]
	} else {
		accessToken = ctx.Query("access_token")
	}

	if accessToken == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":             "invalid_token",
			"error_description": "访问令牌缺失",
		})
	}

	// 验证访问令牌
	userInfo, err := c.authService.ValidateToken(accessToken)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":             "invalid_token",
			"error_description": err.Error(),
		})
	}

	return ctx.JSON(userInfo)
}