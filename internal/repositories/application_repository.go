package repositories

import (
	"github.com/justseemore/sso/internal/models"
	"github.com/justseemore/sso/internal/utils"
	"gorm.io/gorm"
)

type ApplicationRepository struct {
	DB *gorm.DB
}

func NewApplicationRepository() *ApplicationRepository {
	return &ApplicationRepository{
		DB: utils.DB,
	}
}

func (r *ApplicationRepository) Create(app *models.Application) error {
	return r.DB.Create(app).Error
}

func (r *ApplicationRepository) Update(app *models.Application) error {
	return r.DB.Save(app).Error
}

func (r *ApplicationRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Application{}, id).Error
}

func (r *ApplicationRepository) FindByID(id uint) (*models.Application, error) {
	var app models.Application
	err := r.DB.Preload("Theme").First(&app, id).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *ApplicationRepository) FindByClientID(clientID string) (*models.Application, error) {
	var app models.Application
	err := r.DB.Where("client_id = ?", clientID).First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *ApplicationRepository) FindByName(name string) (*models.Application, error) {
	var app models.Application
	err := r.DB.Where("name = ?", name).First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *ApplicationRepository) List(page, limit int) ([]models.Application, int64, error) {
	var apps []models.Application
	var total int64

	r.DB.Model(&models.Application{}).Count(&total)

	offset := (page - 1) * limit
	err := r.DB.Preload("Theme").Limit(limit).Offset(offset).Find(&apps).Error
	if err != nil {
		return nil, 0, err
	}

	return apps, total, nil
}