package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/justseemore/sso/configs"
	"github.com/justseemore/sso/internal/auth"
	"github.com/justseemore/sso/internal/models"
	"github.com/justseemore/sso/internal/repositories"
	"github.com/justseemore/sso/internal/utils"
)

// 定义Redis中使用的键前缀
const (
	AuthCodePrefix              = "auth_code:"
	RefreshTokenPrefix          = "refresh_token:"
	RefreshTokenBlacklistPrefix = "blacklist:refresh_token:"
)

// AuthCodeData 授权码关联的数据结构
type AuthCodeData struct {
	UserID    uint      `json:"user_id"`
	ClientID  string    `json:"client_id"`
	Scopes    []string  `json:"scopes"`
	ExpiredAt time.Time `json:"expired_at"`
}

// RefreshTokenData 刷新令牌关联的数据结构
type RefreshTokenData struct {
	UserID    uint      `json:"user_id"`
	ClientID  string    `json:"client_id"`
	ExpiredAt time.Time `json:"expired_at"`
}

type AuthService struct {
	userRepo *repositories.UserRepository
	appRepo  *repositories.ApplicationRepository
	roleRepo *repositories.RoleRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repositories.NewUserRepository(),
		appRepo:  repositories.NewApplicationRepository(),
		roleRepo: repositories.NewRoleRepository(),
	}
}

// ValidateClientCredentials 验证客户端凭证
func (s *AuthService) ValidateClientCredentials(clientID, clientSecret string) (*models.Application, error) {
	// 获取应用
	app, err := s.appRepo.FindByClientID(clientID)
	if err != nil {
		return nil, errors.New("客户端ID无效")
	}

	// 验证客户端密钥
	if app.ClientSecret != clientSecret {
		return nil, errors.New("客户端密钥无效")
	}

	// 检查应用状态
	if !app.Active {
		return nil, errors.New("应用已被禁用")
	}

	return app, nil
}

// AuthorizeUser 授权用户访问应用
func (s *AuthService) AuthorizeUser(userID uint, clientID string, scopes []string) (string, error) {
	// 获取用户
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", errors.New("用户不存在")
	}

	// 检查用户状态
	if !user.Active {
		return "", errors.New("用户已被禁用")
	}

	// 获取应用
	app, err := s.appRepo.FindByClientID(clientID)
	if err != nil {
		return "", errors.New("应用不存在")
	}

	// 检查应用状态
	if !app.Active {
		return "", errors.New("应用已被禁用")
	}

	// 验证作用域
	allowedScopes, err := app.GetAllowedScopes()
	if err != nil {
		return "", err
	}

	// 检查请求的作用域是否为空
	if len(scopes) == 0 {
		return "", errors.New("请求的作用域不能为空")
	}

	validScopes := []string{}
	for _, scope := range scopes {
		for _, allowedScope := range allowedScopes {
			if scope == allowedScope {
				validScopes = append(validScopes, scope)
				break
			}
		}
	}

	// 检查是否有有效的作用域
	if len(validScopes) == 0 {
		return "", errors.New("没有有效的作用域")
	}

	// 生成授权码
	authCode, err := auth.GenerateRandomString(32)
	if err != nil {
		return "", err
	}

	// 存储授权码，关联用户、应用和作用域，设置过期时间
	ctx := context.Background()
	expiredAt := time.Now().Add(time.Duration(configs.AppConfig.AuthCodeExpiry) * time.Second)

	authData := AuthCodeData{
		UserID:    userID,
		ClientID:  clientID,
		Scopes:    validScopes,
		ExpiredAt: expiredAt,
	}

	data, err := json.Marshal(authData)
	if err != nil {
		return "", err
	}

	key := AuthCodePrefix + authCode
	err = utils.RedisClient.Set(ctx, key, string(data), time.Duration(configs.AppConfig.AuthCodeExpiry)*time.Second).Err()
	if err != nil {
		return "", err
	}

	return authCode, nil
}

// ExchangeToken 使用授权码交换令牌
func (s *AuthService) ExchangeToken(authCode, clientID, clientSecret string) (*auth.TokenDetails, error) {
	// 验证客户端凭证
	_, err := s.ValidateClientCredentials(clientID, clientSecret)
	if err != nil {
		return nil, err
	}

	// 验证授权码，获取关联的用户ID和作用域
	ctx := context.Background()
	key := AuthCodePrefix + authCode

	data, err := utils.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, errors.New("无效的授权码或授权码已过期")
	}

	var authData AuthCodeData
	err = json.Unmarshal([]byte(data), &authData)
	if err != nil {
		return nil, err
	}

	// 验证客户端ID是否匹配
	if authData.ClientID != clientID {
		return nil, errors.New("授权码与客户端ID不匹配")
	}

	// 生成令牌
	tokenDetails, err := auth.GenerateTokens(authData.UserID)
	if err != nil {
		return nil, err
	}

	// 存储刷新令牌，关联用户和应用
	expiredAt := time.Now().Add(time.Duration(configs.AppConfig.RefreshTokenExpiry) * time.Minute)

	refreshData := RefreshTokenData{
		UserID:    authData.UserID,
		ClientID:  clientID,
		ExpiredAt: expiredAt,
	}

	refreshDataStr, err := json.Marshal(refreshData)
	if err != nil {
		return nil, err
	}

	refreshKey := RefreshTokenPrefix + tokenDetails.RefreshToken
	err = utils.RedisClient.Set(
		ctx,
		refreshKey,
		string(refreshDataStr),
		time.Duration(configs.AppConfig.RefreshTokenExpiry)*time.Minute,
	).Err()
	if err != nil {
		return nil, err
	}

	// 删除已使用的授权码
	utils.RedisClient.Del(ctx, key)

	return tokenDetails, nil
}

// RefreshToken 刷新令牌
func (s *AuthService) RefreshToken(refreshToken, clientID, clientSecret string) (*auth.TokenDetails, error) {
	// 验证客户端凭证
	_, err := s.ValidateClientCredentials(clientID, clientSecret)
	if err != nil {
		return nil, err
	}

	// 验证刷新令牌，获取关联的用户ID
	ctx := context.Background()
	key := RefreshTokenPrefix + refreshToken

	// 检查是否在黑名单中
	blacklistKey := RefreshTokenBlacklistPrefix + refreshToken
	exists, err := utils.RedisClient.Exists(ctx, blacklistKey).Result()
	if err != nil || exists > 0 {
		return nil, errors.New("刷新令牌已被撤销")
	}

	data, err := utils.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, errors.New("无效的刷新令牌或令牌已过期")
	}

	var refreshData RefreshTokenData
	err = json.Unmarshal([]byte(data), &refreshData)
	if err != nil {
		return nil, err
	}

	// 验证客户端ID是否匹配
	if refreshData.ClientID != clientID {
		return nil, errors.New("刷新令牌与客户端ID不匹配")
	}

	// 生成新的令牌
	tokenDetails, err := auth.GenerateTokens(refreshData.UserID)
	if err != nil {
		return nil, err
	}

	// 将旧的刷新令牌加入黑名单
	blacklistExpiry := time.Until(refreshData.ExpiredAt)
	if blacklistExpiry > 0 {
		utils.RedisClient.Set(
			ctx,
			RefreshTokenBlacklistPrefix+refreshToken,
			"revoked",
			blacklistExpiry,
		)
	}

	// 存储新的刷新令牌
	newExpiredAt := time.Now().Add(time.Duration(configs.AppConfig.RefreshTokenExpiry) * time.Minute)

	newRefreshData := RefreshTokenData{
		UserID:    refreshData.UserID,
		ClientID:  clientID,
		ExpiredAt: newExpiredAt,
	}

	newRefreshDataStr, err := json.Marshal(newRefreshData)
	if err != nil {
		return nil, err
	}

	newRefreshKey := RefreshTokenPrefix + tokenDetails.RefreshToken
	err = utils.RedisClient.Set(
		ctx,
		newRefreshKey,
		string(newRefreshDataStr),
		time.Duration(configs.AppConfig.RefreshTokenExpiry)*time.Minute,
	).Err()
	if err != nil {
		return nil, err
	}

	return tokenDetails, nil
}

// ValidateToken 验证令牌
func (s *AuthService) ValidateToken(tokenString string) (*auth.Claims, error) {
	return auth.ValidateToken(tokenString)
}

// CheckPermission 检查用户是否有权限
func (s *AuthService) CheckPermission(userID uint, resource, action string) (bool, error) {
	// 获取用户的角色
	roles, err := s.userRepo.GetUserRoles(userID)
	if err != nil {
		return false, err
	}

	// 检查每个角色的权限
	for _, role := range roles {
		permissions, err := s.roleRepo.GetRolePermissions(role.ID)
		if err != nil {
			continue
		}

		for _, permission := range permissions {
			if permission.Resource == resource && permission.Action == action {
				return true, nil
			}
		}
	}

	return false, nil
}

// ExchangeCodeForTokens 使用授权码交换访问令牌和刷新令牌
func (s *AuthService) ExchangeCodeForTokens(code, clientID, redirectURI string) (*auth.TokenDetails, error) {
	// 验证客户端ID
	app, err := s.appRepo.FindByClientID(clientID)
	if err != nil {
		return nil, errors.New("客户端ID无效")
	}

	// 检查应用状态
	if !app.Active {
		return nil, errors.New("应用已被禁用")
	}

	// 验证重定向URI
	allowedURIs, err := app.GetRedirectURIs()
	if err != nil {
		return nil, err
	}

	validURI := false
	for _, uri := range allowedURIs {
		if uri == redirectURI {
			validURI = true
			break
		}
	}

	if !validURI {
		return nil, errors.New("重定向URI无效")
	}

	// 从Redis中获取授权码关联的用户信息
	ctx := context.Background()
	key := AuthCodePrefix + code

	data, err := utils.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, errors.New("无效的授权码或授权码已过期")
	}

	var authData AuthCodeData
	err = json.Unmarshal([]byte(data), &authData)
	if err != nil {
		return nil, err
	}

	// 验证客户端ID是否匹配
	if authData.ClientID != clientID {
		return nil, errors.New("授权码与客户端ID不匹配")
	}

	// 生成令牌
	tokens, err := auth.GenerateTokens(authData.UserID)
	if err != nil {
		return nil, err
	}

	// 存储刷新令牌，关联用户和应用
	expiredAt := time.Now().Add(time.Duration(configs.AppConfig.RefreshTokenExpiry) * time.Minute)

	refreshData := RefreshTokenData{
		UserID:    authData.UserID,
		ClientID:  clientID,
		ExpiredAt: expiredAt,
	}

	refreshDataStr, err := json.Marshal(refreshData)
	if err != nil {
		return nil, err
	}

	refreshKey := RefreshTokenPrefix + tokens.RefreshToken
	err = utils.RedisClient.Set(
		ctx,
		refreshKey,
		string(refreshDataStr),
		time.Duration(configs.AppConfig.RefreshTokenExpiry)*time.Minute,
	).Err()
	if err != nil {
		return nil, err
	}

	// 删除已使用的授权码
	utils.RedisClient.Del(ctx, key)

	return tokens, nil
}

// RefreshTokens 刷新令牌（不需要客户端密钥）
func (s *AuthService) RefreshTokens(refreshToken, clientID string) (*auth.TokenDetails, error) {
	// 验证客户端ID
	app, err := s.appRepo.FindByClientID(clientID)
	if err != nil {
		return nil, errors.New("客户端ID无效")
	}

	// 检查应用状态
	if !app.Active {
		return nil, errors.New("应用已被禁用")
	}

	// 从Redis获取刷新令牌信息
	ctx := context.Background()

	// 检查刷新令牌是否在黑名单中
	blacklistKey := RefreshTokenBlacklistPrefix + refreshToken
	exists, err := utils.RedisClient.Exists(ctx, blacklistKey).Result()
	if err != nil || exists > 0 {
		return nil, errors.New("刷新令牌已被撤销")
	}

	// 获取刷新令牌数据
	key := RefreshTokenPrefix + refreshToken
	data, err := utils.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, errors.New("无效的刷新令牌或令牌已过期")
	}

	var refreshData RefreshTokenData
	err = json.Unmarshal([]byte(data), &refreshData)
	if err != nil {
		return nil, err
	}

	// 验证客户端ID是否匹配
	if refreshData.ClientID != clientID {
		return nil, errors.New("刷新令牌与客户端ID不匹配")
	}

	// 获取用户
	user, err := s.userRepo.FindByID(refreshData.UserID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 检查用户状态
	if !user.Active {
		return nil, errors.New("用户已被禁用")
	}

	// 生成新的令牌
	tokens, err := auth.GenerateTokens(user.ID)
	if err != nil {
		return nil, err
	}

	// 将旧的刷新令牌加入黑名单
	blacklistExpiry := time.Until(refreshData.ExpiredAt)
	if blacklistExpiry > 0 {
		utils.RedisClient.Set(
			ctx,
			RefreshTokenBlacklistPrefix+refreshToken,
			"revoked",
			blacklistExpiry,
		)
	}

	// 存储新的刷新令牌
	newExpiredAt := time.Now().Add(time.Duration(configs.AppConfig.RefreshTokenExpiry) * time.Minute)

	newRefreshData := RefreshTokenData{
		UserID:    user.ID,
		ClientID:  clientID,
		ExpiredAt: newExpiredAt,
	}

	newRefreshDataStr, err := json.Marshal(newRefreshData)
	if err != nil {
		return nil, err
	}

	newRefreshKey := RefreshTokenPrefix + tokens.RefreshToken
	err = utils.RedisClient.Set(
		ctx,
		newRefreshKey,
		string(newRefreshDataStr),
		time.Duration(configs.AppConfig.RefreshTokenExpiry)*time.Minute,
	).Err()
	if err != nil {
		return nil, err
	}

	return tokens, nil
}
