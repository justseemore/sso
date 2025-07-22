package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/justseemore/sso/internal/services"
	"strings"
)

// AuthMiddleware 用于验证JWT令牌
func AuthMiddleware() fiber.Handler {
	authService := services.NewAuthService()

	return func(c *fiber.Ctx) error {
		// 获取Authorization头
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "未提供授权令牌",
			})
		}

		// 解析Bearer令牌
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "授权格式无效",
			})
		}

		tokenString := parts[1]

		// 验证令牌
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// 将用户ID存储在上下文中
		c.Locals("userID", claims.UserID)
		c.Locals("uuid", claims.UUID)

		return c.Next()
	}
}

// OptionalAuthMiddleware 尝试验证令牌，但不会阻止请求
func OptionalAuthMiddleware() fiber.Handler {
	authService := services.NewAuthService()

	return func(c *fiber.Ctx) error {
		// 获取Authorization头
		authHeader := c.Get("Authorization")
		if authHeader != "" {
			// 解析Bearer令牌
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString := parts[1]

				// 验证令牌
				claims, err := authService.ValidateToken(tokenString)
				if err == nil {
					// 将用户ID存储在上下文中
					c.Locals("userID", claims.UserID)
					c.Locals("uuid", claims.UUID)
				}
			}
		}

		return c.Next()
	}
}

// PermissionMiddleware 用于检查用户权限
func PermissionMiddleware(resource, action string) fiber.Handler {
	authService := services.NewAuthService()

	return func(c *fiber.Ctx) error {
		// 获取用户ID
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "未授权，请先登录",
			})
		}

		// 检查权限
		hasPermission, err := authService.CheckPermission(userID, resource, action)
		if err != nil || !hasPermission {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "没有执行此操作的权限",
			})
		}

		return c.Next()
	}
}