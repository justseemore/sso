package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/justseemore/sso/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接并返回数据库实例
func InitDatabase() (*gorm.DB, error) {
	ConnectDatabase()
	return DB, nil
}

func ConnectDatabase() {
	config := configs.AppConfig

	// 构建DSN（数据源名称）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUser, config.DBPass, config.DBHost, config.DBPort, config.DBName)

	// 配置GORM日志
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // 慢SQL阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,       // 忽略记录未找到错误
			Colorful:                  true,       // 启用彩色打印
		},
	)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	log.Println("数据库连接成功")

	// 设置连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取数据库连接池失败: %v", err)
	}

	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(10)
	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(100)
	// 设置连接最大生命周期
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
}