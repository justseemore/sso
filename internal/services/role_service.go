package services

import (
	"errors"
	"time"

	"github.com/justseemore/sso/internal/models"
	"github.com/justseemore/sso/internal/repositories"
)

type RoleService struct {
	roleRepo       *repositories.RoleRepository
	permissionRepo *repositories.PermissionRepository
}

func NewRoleService() *RoleService {
	return &RoleService{
		roleRepo:       repositories.NewRoleRepository(),
		permissionRepo: repositories.NewPermissionRepository(),
	}
}

func (s *RoleService) CreateRole(role *models.Role) error {
	// 检查角色名是否已存在
	existRole, _ := s.roleRepo.FindByName(role.Name)
	if existRole != nil {
		return errors.New("角色名已存在")
	}

	// 设置创建时间和更新时间
	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()

	return s.roleRepo.Create(role)
}

func (s *RoleService) UpdateRole(role *models.Role) error {
	// 检查角色是否存在
	existRole, err := s.roleRepo.FindByID(role.ID)
	if err != nil {
		return errors.New("角色不存在")
	}

	// 如果角色名变了，检查新的角色名是否已存在
	if role.Name != existRole.Name {
		existRole, _ := s.roleRepo.FindByName(role.Name)
		if existRole != nil {
			return errors.New("角色名已存在")
		}
	}

	// 更新时间
	role.UpdatedAt = time.Now()
	return s.roleRepo.Update(role)
}

func (s *RoleService) DeleteRole(id uint) error {
	return s.roleRepo.Delete(id)
}

func (s *RoleService) GetRoleByID(id uint) (*models.Role, error) {
	return s.roleRepo.FindByID(id)
}

func (s *RoleService) GetRoleByName(name string) (*models.Role, error) {
	return s.roleRepo.FindByName(name)
}

func (s *RoleService) ListRoles(page, limit int) ([]models.Role, int64, error) {
	return s.roleRepo.List(page, limit)
}

func (s *RoleService) AssignPermission(roleID, permissionID uint) error {
	// 验证角色和权限是否存在
	_, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return errors.New("角色不存在")
	}

	_, err = s.permissionRepo.FindByID(permissionID)
	if err != nil {
		return errors.New("权限不存在")
	}

	return s.roleRepo.AssignPermission(roleID, permissionID)
}

func (s *RoleService) RemovePermission(roleID, permissionID uint) error {
	return s.roleRepo.RemovePermission(roleID, permissionID)
}

func (s *RoleService) GetRolePermissions(roleID uint) ([]models.Permission, error) {
	return s.roleRepo.GetRolePermissions(roleID)
}

// 权限相关
func (s *RoleService) CreatePermission(permission *models.Permission) error {
	// 检查权限名是否已存在
	existPermission, _ := s.permissionRepo.FindByName(permission.Name)
	if existPermission != nil {
		return errors.New("权限名已存在")
	}

	// 检查资源和操作组合是否已存在
	existPermission, _ = s.permissionRepo.FindByResourceAction(permission.Resource, permission.Action)
	if existPermission != nil {
		return errors.New("该资源的操作权限已存在")
	}

	// 设置创建时间和更新时间
	permission.CreatedAt = time.Now()
	permission.UpdatedAt = time.Now()

	return s.permissionRepo.Create(permission)
}

func (s *RoleService) UpdatePermission(permission *models.Permission) error {
	// 更新时间
	permission.UpdatedAt = time.Now()
	return s.permissionRepo.Update(permission)
}

func (s *RoleService) DeletePermission(id uint) error {
	return s.permissionRepo.Delete(id)
}

func (s *RoleService) GetPermissionByID(id uint) (*models.Permission, error) {
	return s.permissionRepo.FindByID(id)
}

func (s *RoleService) ListPermissions(page, limit int) ([]models.Permission, int64, error) {
	return s.permissionRepo.List(page, limit)
}