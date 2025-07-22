package services

import (
	"errors"
	"time"

	"github.com/justseemore/sso/internal/models"
	"github.com/justseemore/sso/internal/repositories"
)

type ThemeService struct {
	themeRepo *repositories.ThemeRepository
}

func NewThemeService() *ThemeService {
	return &ThemeService{
		themeRepo: repositories.NewThemeRepository(),
	}
}

// CreateTheme 创建主题
func (s *ThemeService) CreateTheme(theme *models.Theme) error {
	// 检查主题名是否已存在
	existTheme, _ := s.themeRepo.FindByName(theme.Name)
	if existTheme != nil {
		return errors.New("主题名已存在")
	}

	// 设置默认值
	theme.Active = true
	theme.CreatedAt = time.Now()
	theme.UpdatedAt = time.Now()

	return s.themeRepo.Create(theme)
}

// UpdateTheme 更新主题
func (s *ThemeService) UpdateTheme(theme *models.Theme) error {
	// 检查主题是否存在
	existTheme, err := s.themeRepo.FindByID(theme.ID)
	if err != nil {
		return errors.New("主题不存在")
	}

	// 如果主题名变了，检查新的主题名是否已存在
	if theme.Name != existTheme.Name {
		existTheme, _ := s.themeRepo.FindByName(theme.Name)
		if existTheme != nil {
			return errors.New("主题名已存在")
		}
	}

	// 更新时间
	theme.UpdatedAt = time.Now()
	return s.themeRepo.Update(theme)
}

// DeleteTheme 删除主题
func (s *ThemeService) DeleteTheme(id uint) error {
	return s.themeRepo.Delete(id)
}

// GetThemeByID 通过ID获取主题
func (s *ThemeService) GetThemeByID(id uint) (*models.Theme, error) {
	return s.themeRepo.FindByID(id)
}

// GetThemeByName 通过名称获取主题
func (s *ThemeService) GetThemeByName(name string) (*models.Theme, error) {
	return s.themeRepo.FindByName(name)
}

// ListThemes 列出主题
func (s *ThemeService) ListThemes(page, limit int) ([]models.Theme, int64, error) {
	return s.themeRepo.List(page, limit)
}