package repositories

import (
	"github.com/justseemore/sso/internal/models"
	"github.com/justseemore/sso/internal/utils"
	"gorm.io/gorm"
)

type PermissionRepository struct {
	DB *gorm.DB
}

func NewPermissionRepository() *PermissionRepository {
	return &PermissionRepository{
		DB: utils.DB,
	}
}

func (r *PermissionRepository) Create(permission *models.Permission) error {
	return r.DB.Create(permission).Error
}

func (r *PermissionRepository) Update(permission *models.Permission) error {
	return r.DB.Save(permission).Error
}

func (r *PermissionRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Permission{}, id).Error
}

func (r *PermissionRepository) FindByID(id uint) (*models.Permission, error) {
	var permission models.Permission
	err := r.DB.First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *PermissionRepository) FindByName(name string) (*models.Permission, error) {
	var permission models.Permission
	err := r.DB.Where("name = ?", name).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *PermissionRepository) FindByResourceAction(resource, action string) (*models.Permission, error) {
	var permission models.Permission
	err := r.DB.Where("resource = ? AND action = ?", resource, action).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *PermissionRepository) List(page, limit int) ([]models.Permission, int64, error) {
	var permissions []models.Permission
	var total int64

	r.DB.Model(&models.Permission{}).Count(&total)

	offset := (page - 1) * limit
	err := r.DB.Limit(limit).Offset(offset).Find(&permissions).Error
	if err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}