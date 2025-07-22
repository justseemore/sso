package services

import (
	"errors"
	"time"

	"github.com/justseemore/sso/internal/auth"
	"github.com/justseemore/sso/internal/models"
	"github.com/justseemore/sso/internal/repositories"
)

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

	validScopes := []string{}
	for _, scope := range scopes {
		for _, allowedScope := range allowedScopes {
			if scope == allowedScope {
				validScopes = append(validScopes, scope)
				break
			}
		}
	}

	// 生成授权码
	authCode, err := auth.GenerateRandomString(32)
	if err != nil {
		return "", err
	}

	// TODO: 存储授权码，关联用户、应用和作用域，设置过期时间
	// 这里可以使用 Redis 或数据库实现

	return authCode, nil
}

// ExchangeToken 使用授权码交换令牌
func (s *AuthService) ExchangeToken(authCode, clientID, clientSecret string) (*auth.TokenDetails, error) {
	// 验证客户端凭证
	app, err := s.ValidateClientCredentials(clientID, clientSecret)
	if err != nil {
		return nil, err
	}

	// TODO: 验证授权码，获取关联的用户ID和作用域
	// 这里假设已经验证了授权码，获取到了用户ID
	// 实际实现需要从存储中获取授权码关联的信息
	userID := uint(1) // 假设用户ID为1

	// 生成令牌
	tokenDetails, err := auth.GenerateTokens(userID)
	if err != nil {
		return nil, err
	}

	// TODO: 存储刷新令牌，关联用户和应用
	// 这里可以使用 Redis 或数据库实现

	return tokenDetails, nil
}

// RefreshToken 刷新令牌
func (s *AuthService) RefreshToken(refreshToken, clientID, clientSecret string) (*auth.TokenDetails, error) {
	// 验证客户端凭证
	_, err := s.ValidateClientCredentials(clientID, clientSecret)
	if err != nil {
		return nil, err
	}

	// TODO: 验证刷新令牌，获取关联的用户ID
	// 这里假设已经验证了刷新令牌，获取到了用户ID
	// 实际实现需要从存储中获取刷新令牌关联的信息
	userID := uint(1) // 假设用户ID为1

	// 生成新的令牌
	tokenDetails, err := auth.GenerateTokens(userID)
	if err != nil {
		return nil, err
	}

	// TODO: 更新存储中的刷新令牌
	// 这里可以使用 Redis 或数据库实现

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