package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/justseemore/sso/configs"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

// 自定义JWT声明结构
type Claims struct {
	jwt.RegisteredClaims
	UserID uint   `json:"user_id"`
	UUID   string `json:"uuid"`
}

// GenerateTokens 生成访问令牌和刷新令牌
func GenerateTokens(userID uint) (*TokenDetails, error) {
	config := configs.AppConfig
	td := &TokenDetails{}

	// 设置过期时间
	td.AtExpires = time.Now().Add(time.Minute * time.Duration(config.AccessTokenExpiry)).Unix()
	td.RtExpires = time.Now().Add(time.Minute * time.Duration(config.RefreshTokenExpiry)).Unix()

	// 创建唯一标识符
	td.AccessUUID = fmt.Sprintf("%d-%v", userID, time.Now().Unix())
	td.RefreshUUID = fmt.Sprintf("%d-%v-refresh", userID, time.Now().Unix())

	// 创建访问令牌
	atClaims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(td.AtExpires, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        td.AccessUUID,
		},
		UserID: userID,
		UUID:   td.AccessUUID,
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	var err error
	td.AccessToken, err = at.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return nil, err
	}

	// 创建刷新令牌
	rtClaims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(td.RtExpires, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        td.RefreshUUID,
		},
		UserID: userID,
		UUID:   td.RefreshUUID,
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return nil, err
	}

	return td, nil
}

// ValidateToken 验证令牌并返回声明
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return []byte(configs.AppConfig.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("无效的令牌")
	}

	return claims, nil
}

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) (string, error) {
	b := make([]byte, length/2)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
