package repositories

import (
	"github.com/justseemore/sso/internal/models"
	"github.com/justseemore/sso/internal/utils"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		DB: utils.DB,
	}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) Update(user *models.User) error {
	return r.DB.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.DB.Delete(&models.User{}, id).Error
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.DB.Preload("Theme").Preload("UserRoles.Role.Permissions").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) List(page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	r.DB.Model(&models.User{}).Count(&total)

	offset := (page - 1) * limit
	err := r.DB.Preload("Theme").Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) AssignRole(userID, roleID uint) error {
	userRole := models.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	return r.DB.Create(&userRole).Error
}

func (r *UserRepository) RemoveRole(userID, roleID uint) error {
	return r.DB.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&models.UserRole{}).Error
}