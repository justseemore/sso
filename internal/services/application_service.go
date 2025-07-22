package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/justseemore/sso/internal/models"
	"github.com/justseemore/sso/internal/repositories"
)

type ApplicationService struct {
	appRepo   *repositories.ApplicationRepository
	themeRepo *repositories.ThemeRepository
}

func NewApplicationService() *ApplicationService {
	return &ApplicationService{
		appRepo:   repositories.NewApplicationRepository(),
		themeRepo: repositories.NewThemeRepository(),
	}
}

// 生成随机字符串
func generateRandomString(length int) (string, error) {
	b := make([]byte, length/2)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// CreateApplication 创建应用
func (s *ApplicationService) CreateApplication(app *models.Application) error {
	// 检查应用名是否已存在
	existApp, _ := s.appRepo.FindByName(app.Name)
	if existApp != nil {
		return errors.New("应用名已存在")
	}

	// 生成客户端ID和密钥
	clientID, err := generateRandomString(32)
	if err != nil {
		return err
	}

	clientSecret, err := generateRandomString(64)
	if err != nil {
		return err
	}

	app.ClientID = clientID
	app.ClientSecret = clientSecret

	// 设置默认值
	app.Active = true
	app.CreatedAt = time.Now()
	app.UpdatedAt = time.Now()

	// 设置重定向URI和作用域的默认值
	if err := app.SetRedirectURIs([]string{}); err != nil {
		return err
	}

	if err := app.SetAllowedScopes([]string{"openid", "profile", "email"}); err != nil {
		return err
	}

	return s.appRepo.Create(app)
}

// UpdateApplication 更新应用
func (s *ApplicationService) UpdateApplication(app *models.Application) error {
	// 检查应用是否存在
	existApp, err := s.appRepo.FindByID(app.ID)
	if err != nil {
		return errors.New("应用不存在")
	}

	// 如果应用名变了，检查新的应用名是否已存在
	if app.Name != existApp.Name {
		existApp, _ := s.appRepo.FindByName(app.Name)
		if existApp != nil {
			return errors.New("应用名已存在")
		}
	}

	// 保留原有的客户端ID和密钥
	app.ClientID = existApp.ClientID
	app.ClientSecret = existApp.ClientSecret

	// 更新时间
	app.UpdatedAt = time.Now()
	return s.appRepo.Update(app)
}

// DeleteApplication 删除应用
func (s *ApplicationService) DeleteApplication(id uint) error {
	return s.appRepo.Delete(id)
}

// GetApplicationByID 通过ID获取应用
func (s *ApplicationService) GetApplicationByID(id uint) (*models.Application, error) {
	return s.appRepo.FindByID(id)
}

// GetApplicationByClientID 通过客户端ID获取应用
func (s *ApplicationService) GetApplicationByClientID(clientID string) (*models.Application, error) {
	return s.appRepo.FindByClientID(clientID)
}

// ListApplications 列出应用
func (s *ApplicationService) ListApplications(page, limit int) ([]models.Application, int64, error) {
	return s.appRepo.List(page, limit)
}

// RegenerateClientSecret 重新生成客户端密钥
func (s *ApplicationService) RegenerateClientSecret(id uint) (string, error) {
	// 获取应用
	app, err := s.appRepo.FindByID(id)
	if err != nil {
		return "", errors.New("应用不存在")
	}

	// 生成新的客户端密钥
	clientSecret, err := generateRandomString(64)
	if err != nil {
		return "", err
	}

	app.ClientSecret = clientSecret
	app.UpdatedAt = time.Now()

	// 更新应用
	err = s.appRepo.Update(app)
	if err != nil {
		return "", err
	}

	return clientSecret, nil
}

// UpdateApplicationTheme 更新应用主题
func (s *ApplicationService) UpdateApplicationTheme(appID, themeID uint) error {
	// 获取应用
	app, err := s.appRepo.FindByID(appID)
	if err != nil {
		return errors.New("应用不存在")
	}

	// 检查主题是否存在
	_, err = s.themeRepo.FindByID(themeID)
	if err != nil {
		return errors.New("主题不存在")
	}

	// 更新主题
	app.ThemeID = &themeID
	app.UpdatedAt = time.Now()
	return s.appRepo.Update(app)
}

// UpdateRedirectURIs 更新重定向URI
func (s *ApplicationService) UpdateRedirectURIs(appID uint, uris []string) error {
	// 获取应用
	app, err := s.appRepo.FindByID(appID)
	if err != nil {
		return errors.New("应用不存在")
	}

	// 更新重定向URI
	err = app.SetRedirectURIs(uris)
	if err != nil {
		return err
	}

	// 更新应用
	app.UpdatedAt = time.Now()
	return s.appRepo.Update(app)
}

// UpdateAllowedScopes 更新允许的作用域
func (s *ApplicationService) UpdateAllowedScopes(appID uint, scopes []string) error {
	// 获取应用
	app, err := s.appRepo.FindByID(appID)
	if err != nil {
		return errors.New("应用不存在")
	}

	// 更新允许的作用域
	err = app.SetAllowedScopes(scopes)
	if err != nil {
		return err
	}

	// 更新应用
	app.UpdatedAt = time.Now()
	return s.appRepo.Update(app)
}

// UpdateSettings 更新应用设置
func (s *ApplicationService) UpdateSettings(appID uint, settings map[string]interface{}) error {
	// 获取应用
	app, err := s.appRepo.FindByID(appID)
	if err != nil {
		return errors.New("应用不存在")
	}

	// 更新设置
	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	app.Settings = settingsJSON
	app.UpdatedAt = time.Now()
	return s.appRepo.Update(app)
}