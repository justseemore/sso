package repositories

import (
	"github.com/justseemore/sso/internal/models"
	"github.com/justseemore/sso/internal/utils"
	"gorm.io/gorm"
)

type ThemeRepository struct {
	DB *gorm.DB
}

func NewThemeRepository() *ThemeRepository {
	return &ThemeRepository{
		DB: utils.DB,
	}
}

func (r *ThemeRepository) Create(theme *models.Theme) error {
	return r.DB.Create(theme).Error
}

func (r *ThemeRepository) Update(theme *models.Theme) error {
	return r.DB.Save(theme).Error
}

func (r *ThemeRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Theme{}, id).Error
}

func (r *ThemeRepository) FindByID(id uint) (*models.Theme, error) {
	var theme models.Theme
	err := r.DB.First(&theme, id).Error
	if err != nil {
		return nil, err
	}
	return &theme, nil
}

func (r *ThemeRepository) FindByName(name string) (*models.Theme, error) {
	var theme models.Theme
	err := r.DB.Where("name = ?", name).First(&theme).Error
	if err != nil {
		return nil, err
	}
	return &theme, nil
}

func (r *ThemeRepository) List(page, limit int) ([]models.Theme, int64, error) {
	var themes []models.Theme
	var total int64

	r.DB.Model(&models.Theme{}).Count(&total)

	offset := (page - 1) * limit
	err := r.DB.Limit(limit).Offset(offset).Find(&themes).Error
	if err != nil {
		return nil, 0, err
	}

	return themes, total, nil
}