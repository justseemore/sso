package configs

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	DBHost             string `mapstructure:"DB_HOST"`
	DBPort             string `mapstructure:"DB_PORT"`
	DBUser             string `mapstructure:"DB_USER"`
	DBPass             string `mapstructure:"DB_PASS"`
	DBName             string `mapstructure:"DB_NAME"`
	APIPort            string `mapstructure:"API_PORT"`
	JWTSecret          string `mapstructure:"JWT_SECRET"`
	AccessTokenExpiry  int    `mapstructure:"ACCESS_TOKEN_EXPIRY"`
	RefreshTokenExpiry int    `mapstructure:"REFRESH_TOKEN_EXPIRY"`
	// 新增Redis配置
	RedisHost      string `mapstructure:"REDIS_HOST"`
	RedisPort      string `mapstructure:"REDIS_PORT"`
	RedisPassword  string `mapstructure:"REDIS_PASSWORD"`
	RedisDB        int    `mapstructure:"REDIS_DB"`
	AuthCodeExpiry int    `mapstructure:"AUTH_CODE_EXPIRY"` // 授权码过期时间（秒）
}

var AppConfig Config

// InitConfig 初始化配置
func InitConfig() {
	AppConfig = LoadConfig()
	log.Println("配置加载成功")
}

func LoadConfig() Config {
	// 加载.env文件
	err := godotenv.Load()
	if err != nil {
		log.Println("警告: 无法找到.env文件，将使用环境变量")
	}

	// 使用viper从环境变量加载配置
	viper.AutomaticEnv()

	// 读取环境变量
	AppConfig = Config{
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "3306"),
		DBUser:             getEnv("DB_USER", "root"),
		DBPass:             getEnv("DB_PASS", "password"),
		DBName:             getEnv("DB_NAME", "sso_db"),
		APIPort:            getEnv("API_PORT", "8080"),
		JWTSecret:          getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
		AccessTokenExpiry:  getEnvAsInt("ACCESS_TOKEN_EXPIRY", 15),
		RefreshTokenExpiry: getEnvAsInt("REFRESH_TOKEN_EXPIRY", 10080),
		// 新增Redis配置默认值
		RedisHost:      getEnv("REDIS_HOST", "localhost"),
		RedisPort:      getEnv("REDIS_PORT", "6379"),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
		RedisDB:        getEnvAsInt("REDIS_DB", 0),
		AuthCodeExpiry: getEnvAsInt("AUTH_CODE_EXPIRY", 600), // 默认10分钟
	}

	return AppConfig
}

func getEnv(key, defaultValue string) string {
	if value, ok := viper.Get(key).(string); ok && value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, ok := viper.Get(key).(int); ok {
		return value
	}
	return defaultValue
}
