package routes

import (
	"github.com/gofiber/fiber/v2"
	_ "github.com/gofiber/template/html/v2"
	"github.com/justseemore/sso/internal/controllers"
	"github.com/justseemore/sso/internal/middlewares"
)

// SetupRoutes 设置路由
func SetupRoutes(app *fiber.App) {
	// 控制器实例
	userController := controllers.NewUserController()
	roleController := controllers.NewRoleController()
	applicationController := controllers.NewApplicationController()
	themeController := controllers.NewThemeController()
	authController := controllers.NewAuthController()

	// API 路由组
	api := app.Group("/api")

	// 不需要认证的路由
	api.Post("/register", userController.Register)
	api.Post("/login", userController.Login)

	// 用户相关路由
	users := api.Group("/users", middlewares.AuthMiddleware())
	users.Get("/", middlewares.PermissionMiddleware("user", "list"), userController.ListUsers)
	users.Get("/:id", middlewares.PermissionMiddleware("user", "read"), userController.GetUser)
	users.Put("/:id", middlewares.PermissionMiddleware("user", "update"), userController.UpdateUser)
	users.Delete("/:id", middlewares.PermissionMiddleware("user", "delete"), userController.DeleteUser)
	users.Post("/:id/roles", middlewares.PermissionMiddleware("user", "assign_role"), userController.AssignRole)
	users.Delete("/:id/roles/:roleId", middlewares.PermissionMiddleware("user", "remove_role"), userController.RemoveRole)
	users.Put("/:id/password", middlewares.PermissionMiddleware("user", "change_password"), userController.ChangePassword)

	// 角色相关路由
	roles := api.Group("/roles", middlewares.AuthMiddleware())
	roles.Post("/", middlewares.PermissionMiddleware("role", "create"), roleController.CreateRole)
	roles.Get("/", middlewares.PermissionMiddleware("role", "list"), roleController.ListRoles)
	roles.Get("/:id", middlewares.PermissionMiddleware("role", "read"), roleController.GetRole)
	roles.Put("/:id", middlewares.PermissionMiddleware("role", "update"), roleController.UpdateRole)
	roles.Delete("/:id", middlewares.PermissionMiddleware("role", "delete"), roleController.DeleteRole)
	roles.Post("/:id/permissions", middlewares.PermissionMiddleware("role", "assign_permission"), roleController.AssignPermission)
	roles.Delete("/:id/permissions/:permissionId", middlewares.PermissionMiddleware("role", "remove_permission"), roleController.RemovePermission)

	// 应用相关路由
	applications := api.Group("/applications", middlewares.AuthMiddleware())
	applications.Post("/", middlewares.PermissionMiddleware("application", "create"), applicationController.CreateApplication)
	applications.Get("/", middlewares.PermissionMiddleware("application", "list"), applicationController.ListApplications)
	applications.Get("/:id", middlewares.PermissionMiddleware("application", "read"), applicationController.GetApplication)
	applications.Put("/:id", middlewares.PermissionMiddleware("application", "update"), applicationController.UpdateApplication)
	applications.Delete("/:id", middlewares.PermissionMiddleware("application", "delete"), applicationController.DeleteApplication)
	applications.Post("/:id/regenerate", middlewares.PermissionMiddleware("application", "update"), applicationController.RegenerateClientSecret)
	applications.Put("/:id/theme", middlewares.PermissionMiddleware("application", "update"), applicationController.UpdateApplicationTheme)
	applications.Put("/:id/redirect-uris", middlewares.PermissionMiddleware("application", "update"), applicationController.UpdateRedirectURIs)
	applications.Put("/:id/allowed-scopes", middlewares.PermissionMiddleware("application", "update"), applicationController.UpdateAllowedScopes)
	applications.Put("/:id/settings", middlewares.PermissionMiddleware("application", "update"), applicationController.UpdateSettings)

	// 主题相关路由
	themes := api.Group("/themes", middlewares.AuthMiddleware())
	themes.Post("/", middlewares.PermissionMiddleware("theme", "create"), themeController.CreateTheme)
	themes.Get("/", middlewares.PermissionMiddleware("theme", "list"), themeController.ListThemes)
	themes.Get("/:id", middlewares.PermissionMiddleware("theme", "read"), themeController.GetTheme)
	themes.Put("/:id", middlewares.PermissionMiddleware("theme", "update"), themeController.UpdateTheme)
	themes.Delete("/:id", middlewares.PermissionMiddleware("theme", "delete"), themeController.DeleteTheme)

	// OAuth 2.0 相关路由
	oauth := app.Group("/oauth")
	oauth.Get("/authorize", middlewares.OptionalAuthMiddleware(), authController.Authorize)
	oauth.Post("/token", authController.Token)

	// 用户信息端点
	app.Get("/userinfo", middlewares.AuthMiddleware(), authController.Userinfo)
}
