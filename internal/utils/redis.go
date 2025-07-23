package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/justseemore/sso/configs"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// InitRedis 初始化Redis客户端
func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", configs.AppConfig.RedisHost, configs.AppConfig.RedisPort),
		Password: configs.AppConfig.RedisPassword,
		DB:       configs.AppConfig.RedisDB,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("无法连接到Redis: %v", err)
	}

	log.Println("Redis连接成功")
}