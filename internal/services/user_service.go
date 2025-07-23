package services

import (
	"errors"
	"time"

	"github.com/justseemore/sso/internal/auth"
	"github.com/justseemore/sso/internal/models"
	"github.com/justseemore/sso/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repositories.UserRepository
	roleRepo *repositories.RoleRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repositories.NewUserRepository(),
		roleRepo: repositories.NewRoleRepository(),
	}
}

func (s *UserService) Register(user *models.User) error {
	// 检查邮箱是否已存在
	existUser, _ := s.userRepo.FindByEmail(user.Email)
	if existUser != nil {
		return errors.New("邮箱已被注册")
	}

	// 检查用户名是否已存在
	existUser, _ = s.userRepo.FindByUsername(user.Username)
	if existUser != nil {
		return errors.New("用户名已被使用")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// 设置默认值
	user.Active = true
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// 保存用户
	return s.userRepo.Create(user)
}

func (s *UserService) Login(username, password string) (*models.User, *auth.TokenDetails, error) {
	// 查找用户
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		// 尝试通过邮箱查找
		user, err = s.userRepo.FindByEmail(username)
		if err != nil {
			return nil, nil, errors.New("用户不存在")
		}
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, nil, errors.New("密码错误")
	}

	// 检查用户状态
	if !user.Active {
		return nil, nil, errors.New("账户已被禁用")
	}

	// 生成令牌
	tokenDetails, err := auth.GenerateTokens(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, tokenDetails, nil
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.FindByEmail(email)
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	return s.userRepo.FindByUsername(username)
}

func (s *UserService) UpdateUser(user *models.User) error {
	// 更新时间
	user.UpdatedAt = time.Now()
	return s.userRepo.Update(user)
}

func (s *UserService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}

func (s *UserService) ListUsers(page, limit int) ([]models.User, int64, error) {
	return s.userRepo.List(page, limit)
}

func (s *UserService) AssignRole(userID, roleID uint) error {
	// 验证用户和角色是否存在
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	_, err = s.roleRepo.FindByID(roleID)
	if err != nil {
		return errors.New("角色不存在")
	}

	return s.userRepo.AssignRole(userID, roleID)
}

func (s *UserService) RemoveRole(userID, roleID uint) error {
	return s.userRepo.RemoveRole(userID, roleID)
}

func (s *UserService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	// 获取用户
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码
	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()
	return s.userRepo.Update(user)
}

func (s *UserService) UpdateUserProfile(userID uint, profile map[string]interface{}) error {
	// 获取用户
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 更新用户自定义属性
	user.SetUserAttributes(profile)

	// 更新用户
	user.UpdatedAt = time.Now()
	return s.userRepo.Update(user)
}

func (s *UserService) UpdateUserTheme(userID, themeID uint) error {
	// 获取用户
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 更新主题
	user.ThemeID = &themeID
	user.UpdatedAt = time.Now()
	return s.userRepo.Update(user)
}