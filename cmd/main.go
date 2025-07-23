package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/justseemore/sso/configs"
	"github.com/justseemore/sso/internal/routes"
	"github.com/justseemore/sso/internal/utils"
	"github.com/joho/godotenv"
)



func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("警告: 未找到.env文件")
	}

	// 初始化配置
	configs.InitConfig()

	// 初始化数据库连接
	db, err := utils.InitDatabase()
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()
     // 在database初始化后添加
    utils.InitRedis()
	// 初始化视图引擎
	viewsEngine := html.New("./web/views", ".html")

	// 创建 Fiber 实例
	app := fiber.New(fiber.Config{
		Views: viewsEngine,
	})

	// 注册中间件
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// 静态文件
	app.Static("/static", "./web/static")

	// 设置路由
	routes.SetupRoutes(app)

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("服务器启动在 http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}