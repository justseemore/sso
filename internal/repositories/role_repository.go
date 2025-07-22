package repositories

import (
	"github.com/justseemore/sso/internal/models"
	"github.com/justseemore/sso/internal/utils"
	"gorm.io/gorm"
)

type RoleRepository struct {
	DB *gorm.DB
}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{
		DB: utils.DB,
	}
}

func (r *RoleRepository) Create(role *models.Role) error {
	return r.DB.Create(role).Error
}

func (r *RoleRepository) Update(role *models.Role) error {
	return r.DB.Save(role).Error
}

func (r *RoleRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Role{}, id).Error
}

func (r *RoleRepository) FindByID(id uint) (*models.Role, error) {
	var role models.Role
	err := r.DB.Preload("Permissions").First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) FindByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.DB.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) List(page, limit int) ([]models.Role, int64, error) {
	var roles []models.Role
	var total int64

	r.DB.Model(&models.Role{}).Count(&total)

	offset := (page - 1) * limit
	err := r.DB.Limit(limit).Offset(offset).Find(&roles).Error
	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func (r *RoleRepository) AssignPermission(roleID, permissionID uint) error {
	return r.DB.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)", roleID, permissionID).Error
}

func (r *RoleRepository) RemovePermission(roleID, permissionID uint) error {
	return r.DB.Exec("DELETE FROM role_permissions WHERE role_id = ? AND permission_id = ?", roleID, permissionID).Error
}

func (r *RoleRepository) GetRolePermissions(roleID uint) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.DB.Model(&models.Role{ID: roleID}).Association("Permissions").Find(&permissions)
	return permissions, err
}